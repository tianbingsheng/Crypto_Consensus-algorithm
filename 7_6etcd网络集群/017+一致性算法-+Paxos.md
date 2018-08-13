# 017 一致性算法- Paxos

Paxos是一种基于消息传递的分布式一致性算法，由Leslie Lamport（莱斯利·兰伯特）于1990提出。是目前公认的解决分布式一致性问题的最有效算法之一。

## Paxos 是什么

Paxos协议是一个解决分布式系统中，多个节点之间就某个值（提案）达成一致（决议）的通信协议。它能够处理在少数节点离线的情况下，剩余的多数节点仍然能达成一致。

## 概念定义

**Proposal**：为了就某一个值达成一致而发起的提案，包括提案编号和提案的值。

涉及角色如下：
　　
　　**Proposer**：提案发起者，为了就某一个值达成一致，Proposer可以以任意速度、发起任意数量的提案，可以停止或重启。
　　
　　**Acceptor**：提案批准者，负责处理接收到的提案，响应、作出承诺、或批准提案。
　　
　　**Learner**：提案学习者，可以从Acceptor处获取已被批准的提案。
 
## Paxos 协议

Paxos需要遵循如下约定：

1. 一个Acceptor必须批准它收到的第一个提案。

2. 如果编号为n的提案被批准了，那么所有编号大于n的提案，其值必须与编号为n的提案的值相同。

Paxos 协议是一个两阶段协议，分为Prepare 阶段 和 Accept阶段。

### Prepare 阶段

#### Proposer 发送 Prepare

`Proposer` 生成一个全局唯一且递增的提案ID，向 `Acceptor` 发送请求，只携带提案ID即可。

#### Acceptor 应答 Prepare

`Acceptor` 接收到提案请求后，如下情况会收到应答

* 当前提交的编号大于之前的其他机器 Prepare 的编号，

* 当前是第一个提交 `Prepare` 的机器

### Accept 阶段

#### Proposer 发送 Accept

如果Proposer收到半数以上Acceptor对其发出的编号为N的Prepare请求的响应，那么它就会发送一个针对[N,V]提案的Accept请求给半数以上的Acceptor。注意：V就是收到的响应中编号最大的提案的value，如果响应中不包含任何提案，那么V就由Proposer自己决定。

#### Acceptor 应答 Accept

如果Acceptor收到一个针对编号为N的提案的Accept请求，只要该Acceptor没有对编号大于N的Prepare请求做出过响应，它就接受该提案。


### Learner
 
　　一旦Acceptor批准了某个提案，即将该提案发给所有的Learner。为了避免大量通信，Acceptor也可以将批准的提案，发给主Learner，由主Learner分发给其他Learner。考虑到主Learner单点问题，也可以考虑Acceptor将批准的提案，发给主Learner组，由主Learner组分发给其他Learner。

## Paxos 算法演示

![](http://olgjbx93m.bkt.clouddn.com/20180125-20310.jpg)


## 死循环问题
 
　　如果Proposer1提出编号为n1的提案，并完成了阶段一。与此同时Proposer2提出了编号为n2的提案，n2>n1，同样也完成了阶段一。于是Acceptor承诺不再批准编号小于n2的提案，当Proposer1进入阶段二时，将会被忽略。同理，此时，Proposer1可以提出编号为n3的提案，n3>n2，又会导致Proposer2的编号为n2的提案进入阶段二时被忽略。以此类推，将进入死循环。
 
　　**解决办法**：
 
　　可以选择一个Proposer作为主Proposer，并约定只有主Proposer才可以提出提案。因此，只要主Proposer可以与过半的Acceptor保持通信，那么但凡主Proposer提出的编号更高的提案，均会被批准。
 

 

