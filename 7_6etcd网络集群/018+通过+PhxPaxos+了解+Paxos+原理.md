# 018 通过 PhxPaxos 了解 Paxos 原理

##  Prepare阶段

### Prepare


```
// src/algorithm/proposer.cpp

void Proposer :: Prepare(const bool bNeedNewBallot)
{
    PLGHead("START Now.InstanceID %lu MyNodeID %lu State.ProposalID %lu State.ValueLen %zu",
            GetInstanceID(), m_poConfig->GetMyNodeID(), m_oProposerState.GetProposalID(),
            m_oProposerState.GetValue().size());

    BP->GetProposerBP()->Prepare();
    m_oTimeStat.Point();
    
    ExitAccept();

    //表明Proposer正处于Prepare阶段
    m_bIsPreparing = true;

    //不能跳过Prepare阶段
    m_bCanSkipPrepare = false;

    //目前还未被任意一个Acceptor拒绝
    m_bWasRejectBySomeone = false;

    m_oProposerState.ResetHighestOtherPreAcceptBallot();

    //如果需要产生新的投票，就调用NewPrepare产生新的ProposalID，新的ProposalID为当前已知的最大ProposalID+1
    if (bNeedNewBallot)
    {
        m_oProposerState.NewPrepare();
    }

    PaxosMsg oPaxosMsg;

    //设置Prepare消息的各个字段
    oPaxosMsg.set_msgtype(MsgType_PaxosPrepare);
    oPaxosMsg.set_instanceid(GetInstanceID());
    oPaxosMsg.set_nodeid(m_poConfig->GetMyNodeID());
    oPaxosMsg.set_proposalid(m_oProposerState.GetProposalID());

    //MsgCount是专门用来统计票数的，根据计算的结果确定是否通过Prepare阶段或者Accept阶段
    m_oMsgCounter.StartNewRound();

    //Prepare超时定时器
    AddPrepareTimer();

    PLGHead("END OK");

    //将Prepare消息发送到各个节点
    BroadcastMessage(oPaxosMsg);
}
```

Proposer在Prepare阶段主要做了这么几件事：

* 重置各个状态位，表明当前正处于Prepare阶段。

* 获取提案编号ProposalID。当bNeedNewBallot为true时需要将ProposalID+1。否则沿用之前的ProposalID。bNeedNewBallot是在NewValue中调用Prepare方法时传入的m_bWasRejectBySomeone参数。也就是如果之前没有被任何Acceptor拒绝（说明还没有明确出现更大的ProposalID），则不需要获取新的ProposalID。对应的场景是Prepare阶段超时了，在超时时间内没有收到过半Acceptor同意的消息，因此需要重新执行Prepare阶段，此时只需要沿用原来的ProposalID即可。

* 发送Prepare请求。该请求PaxosMsg是Protocol Buffer定义的一个message，包含MsgType、InstanceID、NodeID、ProposalID等字段。在BroadcastMessage(oPaxosMsg)中还会将oPaxosMsg序列化后才发送出去。

PaxosMsg的定义如下，Prepare和Accept阶段Proposer和Acceptor的所有消息都用PaxosMsg来表示：

```
// src/comm/paxos_msg.proto

message PaxosMsg
{
	required int32 MsgType = 1;
	optional uint64 InstanceID = 2;
	optional uint64 NodeID = 3;
	optional uint64 ProposalID = 4;
	optional uint64 ProposalNodeID = 5;
	optional bytes Value = 6;
	optional uint64 PreAcceptID = 7;
	optional uint64 PreAcceptNodeID = 8;
	optional uint64 RejectByPromiseID = 9;
	optional uint64 NowInstanceID = 10;
	optional uint64 MinChosenInstanceID = 11;
	optional uint32 LastChecksum = 12;
	optional uint32 Flag = 13;
	optional bytes SystemVariables = 14;
	optional bytes MasterVariables = 15;
};
```

### OnPrepareReply

Proposer发出Prepare请求后就开始等待Acceptor的回复。当Proposer所在节点收到PaxosPrepareReply消息后，就会调用Proposer的OnPrepareReply(oPaxosMsg)，其中oPaxosMsg是Acceptor回复的消息。

```
// src/algorithm/proposer.cpp

void Proposer :: OnPrepareReply(const PaxosMsg & oPaxosMsg)
{
    PLGHead("START Msg.ProposalID %lu State.ProposalID %lu Msg.from_nodeid %lu RejectByPromiseID %lu",
            oPaxosMsg.proposalid(), m_oProposerState.GetProposalID(), 
            oPaxosMsg.nodeid(), oPaxosMsg.rejectbypromiseid());

    BP->GetProposerBP()->OnPrepareReply();
    
    //如果Proposer不是在Prepare阶段，则忽略该消息
    if (!m_bIsPreparing)
    {
        BP->GetProposerBP()->OnPrepareReplyButNotPreparing();
        //PLGErr("Not preparing, skip this msg");
        return;
    }

    //如果ProposalID不同，也忽略
    if (oPaxosMsg.proposalid() != m_oProposerState.GetProposalID())
    {
        BP->GetProposerBP()->OnPrepareReplyNotSameProposalIDMsg();
        //PLGErr("ProposalID not same, skip this msg");
        return;
    }

    //加入一个收到的消息，用于MsgCounter统计
    m_oMsgCounter.AddReceive(oPaxosMsg.nodeid());

    //如果该消息不是拒绝，即Acceptor同意本次Prepare请求
    if (oPaxosMsg.rejectbypromiseid() == 0)
    {
        BallotNumber oBallot(oPaxosMsg.preacceptid(), oPaxosMsg.preacceptnodeid());
        PLGDebug("[Promise] PreAcceptedID %lu PreAcceptedNodeID %lu ValueSize %zu", 
                oPaxosMsg.preacceptid(), oPaxosMsg.preacceptnodeid(), oPaxosMsg.value().size());
        //加入MsgCounter用于统计投票
        m_oMsgCounter.AddPromiseOrAccept(oPaxosMsg.nodeid());
        //将Acceptor返回的它接受过的编号最大的提案记录下来（如果有的话），用于确定Accept阶段的Value
        m_oProposerState.AddPreAcceptValue(oBallot, oPaxosMsg.value());
    }

    //Acceptor拒绝了Prepare请求
    else
    {
        PLGDebug("[Reject] RejectByPromiseID %lu", oPaxosMsg.rejectbypromiseid());

        //同样也要记录到MsgCounter用于统计投票
        m_oMsgCounter.AddReject(oPaxosMsg.nodeid());

        //记录被Acceptor拒绝过，待会儿如果重新进入Prepare阶段的话就需要获取更大的ProposalID
        m_bWasRejectBySomeone = true;

        //记录下别的Proposer提出的更大的ProposalID。这样重新发起Prepare请求时才知道需要用多大的ProposalID
        m_oProposerState.SetOtherProposalID(oPaxosMsg.rejectbypromiseid());
    }


    //本次Prepare请求通过了。也就是得到了半数以上Acceptor的同意
    if (m_oMsgCounter.IsPassedOnThisRound())
    {
        int iUseTimeMs = m_oTimeStat.Point();
        BP->GetProposerBP()->PreparePass(iUseTimeMs);
        PLGImp("[Pass] start accept, usetime %dms", iUseTimeMs);
        m_bCanSkipPrepare = true;

        //进入Accept阶段
        Accept();
    }

    //本次Prepare请求没有通过
    else if (m_oMsgCounter.IsRejectedOnThisRound()
            || m_oMsgCounter.IsAllReceiveOnThisRound())
    {
        BP->GetProposerBP()->PrepareNotPass();
        PLGImp("[Not Pass] wait 30ms and restart prepare");

         //随机等待一段时间后重新发起Prepare请求
        AddPrepareTimer(OtherUtils::FastRand() % 30 + 10);
    }

    PLGHead("END");
}
```

该阶段Proposer主要做了以下事情：

* 判断消息是否有效。包括ProposalID是否相同，自身是否处于Prepare阶段等。因为网络是不可靠的，有些消息可能延迟很久，等收到的时候已经不需要了，所以需要做这些判断。

* 将收到的消息加入MsgCounter用于统计。

* 根据收到的消息更新自身状态。包括Acceptor承诺过的ProposalID，以及Acceptor接受过的编号最大的提案等。

* 根据MsgCounter统计的Acceptor投票结果决定是进入Acceptor阶段还是重新发起Prepare请求。这里如果判断需要重新发起Prepare请求的话，也不是立即进行，而是等待一段随机的时间，这样做的好处是减少不同Proposer之间的冲突，采取的策略跟raft中leader选举冲突时在一段随机的选举超时时间后重新发起选举的做法类似。

注：这里跟Paxos算法中提案编号对应的并不是ProposalID，而是BallotNumber。BallotNumber由ProposalID和NodeID组成。还实现了运算符重载。如果ProposalID大，则BallotNumber（即提案编号）大。在ProposalID相同的情况下，NodeID大的BallotNumber大。

## Accept 阶段

接下来Proposer就进入Accept阶段：

### Accept

```
// src/algorithm/proposer.cpp

void Proposer :: Accept()
{
    PLGHead("START ProposalID %lu ValueSize %zu ValueLen %zu", 
            m_oProposerState.GetProposalID(), m_oProposerState.GetValue().size(), m_oProposerState.GetValue().size());

    BP->GetProposerBP()->Accept();
    m_oTimeStat.Point();
    
    ExitPrepare();
    m_bIsAccepting = true;
    
    //设置Accept请求的消息内容
    PaxosMsg oPaxosMsg;
    oPaxosMsg.set_msgtype(MsgType_PaxosAccept);
    oPaxosMsg.set_instanceid(GetInstanceID());
    oPaxosMsg.set_nodeid(m_poConfig->GetMyNodeID());
    oPaxosMsg.set_proposalid(m_oProposerState.GetProposalID());
    oPaxosMsg.set_value(m_oProposerState.GetValue());
    oPaxosMsg.set_lastchecksum(GetLastChecksum());

    m_oMsgCounter.StartNewRound();

    AddAcceptTimer();

    PLGHead("END");

    //发给各个节点
    BroadcastMessage(oPaxosMsg, BroadcastMessage_Type_RunSelf_Final);
}
```

Accept请求中 `PaxosMsg`里的Value是这样确定的：如果Prepare阶段有Acceptor的回复中带有提案值，则该Value为所有的Acceptor的回复中，编号最大的提案的值。否则就是Proposer在最初调用NewValue时传入的值。

### OnAcceptReply

```
// src/algorithm/proposer.cpp

void Proposer :: OnAcceptReply(const PaxosMsg & oPaxosMsg)
{
    PLGHead("START Msg.ProposalID %lu State.ProposalID %lu Msg.from_nodeid %lu RejectByPromiseID %lu",
            oPaxosMsg.proposalid(), m_oProposerState.GetProposalID(), 
            oPaxosMsg.nodeid(), oPaxosMsg.rejectbypromiseid());

    BP->GetProposerBP()->OnAcceptReply();

    if (!m_bIsAccepting)
    {
        //PLGErr("Not proposing, skip this msg");
        BP->GetProposerBP()->OnAcceptReplyButNotAccepting();
        return;
    }

    if (oPaxosMsg.proposalid() != m_oProposerState.GetProposalID())
    {
        //PLGErr("ProposalID not same, skip this msg");
        BP->GetProposerBP()->OnAcceptReplyNotSameProposalIDMsg();
        return;
    }

    m_oMsgCounter.AddReceive(oPaxosMsg.nodeid());

    if (oPaxosMsg.rejectbypromiseid() == 0)
    {
        PLGDebug("[Accept]");
        m_oMsgCounter.AddPromiseOrAccept(oPaxosMsg.nodeid());
    }
    else
    {
        PLGDebug("[Reject]");
        m_oMsgCounter.AddReject(oPaxosMsg.nodeid());

        m_bWasRejectBySomeone = true;

        m_oProposerState.SetOtherProposalID(oPaxosMsg.rejectbypromiseid());
    }

    if (m_oMsgCounter.IsPassedOnThisRound())
    {
        int iUseTimeMs = m_oTimeStat.Point();
        BP->GetProposerBP()->AcceptPass(iUseTimeMs);
        PLGImp("[Pass] Start send learn, usetime %dms", iUseTimeMs);
        ExitAccept();

        //让Leaner学习被选定（Chosen）的值
        m_poLearner->ProposerSendSuccess(GetInstanceID(), m_oProposerState.GetProposalID());
    }
    else if (m_oMsgCounter.IsRejectedOnThisRound()
            || m_oMsgCounter.IsAllReceiveOnThisRound())
    {
        BP->GetProposerBP()->AcceptNotPass();
        PLGImp("[Not pass] wait 30ms and Restart prepare");
        AddAcceptTimer(OtherUtils::FastRand() % 30 + 10);
    }

    PLGHead("END");
}
```
这里跟OnPrepareReply的过程基本一致。比较大的区别在于最后如果过半的Acceptor接受了该Accept请求，则说明该Value被选定（Chosen）了，就发送消息，让每个节点上的Learner学习该Value。

## Acceptor

### OnPrepare

OnPrepare用于处理收到的Prepare请求，逻辑如下：

```
int Acceptor :: OnPrepare(const PaxosMsg & oPaxosMsg)
{
    PLGHead("START Msg.InstanceID %lu Msg.from_nodeid %lu Msg.ProposalID %lu",
            oPaxosMsg.instanceid(), oPaxosMsg.nodeid(), oPaxosMsg.proposalid());

    BP->GetAcceptorBP()->OnPrepare();
    
    PaxosMsg oReplyPaxosMsg;
    oReplyPaxosMsg.set_instanceid(GetInstanceID());
    oReplyPaxosMsg.set_nodeid(m_poConfig->GetMyNodeID());
    oReplyPaxosMsg.set_proposalid(oPaxosMsg.proposalid());
    oReplyPaxosMsg.set_msgtype(MsgType_PaxosPrepareReply);

    //构造接收到的Prepare请求里的提案编号
    BallotNumber oBallot(oPaxosMsg.proposalid(), oPaxosMsg.nodeid());
    
    //提案编号大于承诺过的提案编号
    if (oBallot >= m_oAcceptorState.GetPromiseBallot())
    {
        PLGDebug("[Promise] State.PromiseID %lu State.PromiseNodeID %lu "
                "State.PreAcceptedID %lu State.PreAcceptedNodeID %lu",
                m_oAcceptorState.GetPromiseBallot().m_llProposalID, 
                m_oAcceptorState.GetPromiseBallot().m_llNodeID,
                m_oAcceptorState.GetAcceptedBallot().m_llProposalID,
                m_oAcceptorState.GetAcceptedBallot().m_llNodeID);
        
        //返回之前接受过的提案的编号
        oReplyPaxosMsg.set_preacceptid(m_oAcceptorState.GetAcceptedBallot().m_llProposalID);
        oReplyPaxosMsg.set_preacceptnodeid(m_oAcceptorState.GetAcceptedBallot().m_llNodeID);
        //如果接受过的提案编号大于0（<=0说明没有接受过提案），则设置接受过的提案的Value
        if (m_oAcceptorState.GetAcceptedBallot().m_llProposalID > 0)
        {
            oReplyPaxosMsg.set_value(m_oAcceptorState.GetAcceptedValue());
        }

        //更新承诺的提案编号为新的提案编号（因为新的提案编号更大）
        m_oAcceptorState.SetPromiseBallot(oBallot);

        //信息持久化
        int ret = m_oAcceptorState.Persist(GetInstanceID(), GetLastChecksum());
        if (ret != 0)
        {
            BP->GetAcceptorBP()->OnPreparePersistFail();
            PLGErr("Persist fail, Now.InstanceID %lu ret %d",
                    GetInstanceID(), ret);
            
            return -1;
        }

        BP->GetAcceptorBP()->OnPreparePass();
    }

    //提案编号小于承诺过的提案编号，需要拒绝
    else
    {
        BP->GetAcceptorBP()->OnPrepareReject();

        PLGDebug("[Reject] State.PromiseID %lu State.PromiseNodeID %lu", 
                m_oAcceptorState.GetPromiseBallot().m_llProposalID, 
                m_oAcceptorState.GetPromiseBallot().m_llNodeID);
        
        //拒绝该Prepare请求，并返回承诺过的ProposalID      
        oReplyPaxosMsg.set_rejectbypromiseid(m_oAcceptorState.GetPromiseBallot().m_llProposalID);
    }

    nodeid_t iReplyNodeID = oPaxosMsg.nodeid();

    PLGHead("END Now.InstanceID %lu ReplyNodeID %lu",
            GetInstanceID(), oPaxosMsg.nodeid());;

    //向发出Prepare请求的Proposer回复消息
    SendMessage(iReplyNodeID, oReplyPaxosMsg);

    return 0;
}
```

### OnAccept

再来看看OnAccept：

```
void Acceptor :: OnAccept(const PaxosMsg & oPaxosMsg)
{
    PLGHead("START Msg.InstanceID %lu Msg.from_nodeid %lu Msg.ProposalID %lu Msg.ValueLen %zu",
            oPaxosMsg.instanceid(), oPaxosMsg.nodeid(), oPaxosMsg.proposalid(), oPaxosMsg.value().size());

    BP->GetAcceptorBP()->OnAccept();

    PaxosMsg oReplyPaxosMsg;
    oReplyPaxosMsg.set_instanceid(GetInstanceID());
    oReplyPaxosMsg.set_nodeid(m_poConfig->GetMyNodeID());
    oReplyPaxosMsg.set_proposalid(oPaxosMsg.proposalid());
    oReplyPaxosMsg.set_msgtype(MsgType_PaxosAcceptReply);

    BallotNumber oBallot(oPaxosMsg.proposalid(), oPaxosMsg.nodeid());

    //提案编号不小于承诺过的提案编号（注意：这里是“>=”，而再OnPrepare中是“>”，可以先思考下为什么），需要接受该提案
    if (oBallot >= m_oAcceptorState.GetPromiseBallot())
    {
        PLGDebug("[Promise] State.PromiseID %lu State.PromiseNodeID %lu "
                "State.PreAcceptedID %lu State.PreAcceptedNodeID %lu",
                m_oAcceptorState.GetPromiseBallot().m_llProposalID, 
                m_oAcceptorState.GetPromiseBallot().m_llNodeID,
                m_oAcceptorState.GetAcceptedBallot().m_llProposalID,
                m_oAcceptorState.GetAcceptedBallot().m_llNodeID);

        //更新承诺的提案编号；接受的提案编号、提案值
        m_oAcceptorState.SetPromiseBallot(oBallot);
        m_oAcceptorState.SetAcceptedBallot(oBallot);
        m_oAcceptorState.SetAcceptedValue(oPaxosMsg.value());
        
        //信息持久化
        int ret = m_oAcceptorState.Persist(GetInstanceID(), GetLastChecksum());
        if (ret != 0)
        {
            BP->GetAcceptorBP()->OnAcceptPersistFail();

            PLGErr("Persist fail, Now.InstanceID %lu ret %d",
                    GetInstanceID(), ret);
            
            return;
        }

        BP->GetAcceptorBP()->OnAcceptPass();
    }

    //需要拒绝该提案
    else
    {
        BP->GetAcceptorBP()->OnAcceptReject();

        PLGDebug("[Reject] State.PromiseID %lu State.PromiseNodeID %lu", 
                m_oAcceptorState.GetPromiseBallot().m_llProposalID, 
                m_oAcceptorState.GetPromiseBallot().m_llNodeID);
        
        //拒绝的消息中附上承诺过的ProposalID
        oReplyPaxosMsg.set_rejectbypromiseid(m_oAcceptorState.GetPromiseBallot().m_llProposalID);
    }

    nodeid_t iReplyNodeID = oPaxosMsg.nodeid();

    PLGHead("END Now.InstanceID %lu ReplyNodeID %lu",
            GetInstanceID(), oPaxosMsg.nodeid());

    //将响应发送给Proposer
    SendMessage(iReplyNodeID, oReplyPaxosMsg);
}
```

