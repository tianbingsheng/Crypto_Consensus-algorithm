# 006 DES、3DES加密算法

DES加密算法，为对称加密算法中的一种。70年代初由IBM研发，后1977年被美国国家标准局采纳为数据加密标准，即DES全称的由来：Data Encryption Standard。对称加密算法，是相对于非对称加密算法而言的。两者区别在于，对称加密在加密和解密时使用同一密钥，而非对称加密在加密和解密时使用不同的密钥，即公钥和私钥。常见的DES、3DES、AES均为对称加密算法，而RSA、椭圆曲线加密算法，均为非对称加密算法。
 
　　DES是以64比特的明文为一个单位来进行加密的，超过64比特的数据，要求按固定的64比特的大小分组，分组有很多模式，后续单独总结，暂时先介绍DES加密算法。DES使用的密钥长度为64比特，但由于每隔7个比特设置一个奇偶校验位，因此其密钥长度实际为56比特。奇偶校验为最简单的错误检测码，即根据一组二进制代码中1的个数是奇数或偶数来检测错误。

## Feistel网络

DES的基本结构，由IBM公司的Horst Feistel设计，因此称Feistel网络。在Feistel网络中，加密的每个步骤称为轮，经过初始置换后的64位明文，进行了16轮Feistel轮的加密过程，最后经过终结置换后形成最终的64位密文。如下为Feistel网络的示意图：

![](http://olgjbx93m.bkt.clouddn.com/5f65e67ee5cbf1514c5614e944684cc1af2a4096.jpg)  

64位明文被分为左、右两部分处理，右侧数据和子密钥经过轮函数f生成用于加密左侧数据的比特序列，与左侧数据异或运算，运算结果输出为加密后的左侧，右侧数据则直接输出为右侧。
　　其中子密钥为本轮加密使用的密钥，每次Feistel均使用不同的子密钥。子密钥的计算，以及轮函数的细节，稍后下文介绍。由于一次Feistel轮并不会加密右侧，因此需要将上一轮输出后的左右两侧对调后，重复Feistel轮的过程，DES算法共计进行16次Feistel轮，最后一轮输出后左右两侧无需对调。

DES加密和解密的过程一致，均使用Feistel网络实现，区别仅在于解密时，密文作为输入，并逆序使用子密钥。
go标准库中DES算法实现如下：

```  
func cryptBlock(subkeys []uint64, dst, src []byte, decrypt bool) {
	b := binary.BigEndian.Uint64(src)
	//初始置换
	b = permuteInitialBlock(b)
	left, right := uint32(b>>32), uint32(b)

	left = (left << 1) | (left >> 31)
	right = (right << 1) | (right >> 31)

	//共计16次feistel轮
	if decrypt {
		for i := 0; i < 8; i++ {
			left, right = feistel(left, right, subkeys[15-2*i], subkeys[15-(2*i+1)])
		}
	} else {
		for i := 0; i < 8; i++ {
			left, right = feistel(left, right, subkeys[2*i], subkeys[2*i+1])
		}
	}

	left = (left << 31) | (left >> 1)
	right = (right << 31) | (right >> 1)

	// 左右切换并执行最终置换
	preOutput := (uint64(right) << 32) | uint64(left)
	binary.BigEndian.PutUint64(dst, permuteFinalBlock(preOutput))
}
```

> DES算法共计进行16次Feistel轮，上面却只有循环了8次
> 这是由于 feistel 方法中一次操作了左右两个参数，所以循环次数减半

```
func feistel(l, r uint32, k0, k1 uint64) (lout, rout uint32) {
	var t uint32

	t = r ^ uint32(k0>>32)
	l ^= feistelBox[7][t&0x3f] ^
		feistelBox[5][(t>>8)&0x3f] ^
		feistelBox[3][(t>>16)&0x3f] ^
		feistelBox[1][(t>>24)&0x3f]

	t = ((r << 28) | (r >> 4)) ^ uint32(k0)
	l ^= feistelBox[6][(t)&0x3f] ^
		feistelBox[4][(t>>8)&0x3f] ^
		feistelBox[2][(t>>16)&0x3f] ^
		feistelBox[0][(t>>24)&0x3f]

	t = l ^ uint32(k1>>32)
	r ^= feistelBox[7][t&0x3f] ^
		feistelBox[5][(t>>8)&0x3f] ^
		feistelBox[3][(t>>16)&0x3f] ^
		feistelBox[1][(t>>24)&0x3f]

	t = ((l << 28) | (l >> 4)) ^ uint32(k1)
	r ^= feistelBox[6][(t)&0x3f] ^
		feistelBox[4][(t>>8)&0x3f] ^
		feistelBox[2][(t>>16)&0x3f] ^
		feistelBox[0][(t>>24)&0x3f]

	return l, r
}
```

### 初始置换和终结置换

进入Feistel轮之前，64位明文需做一次初始置换。Feistel轮结束后，需做一次反向操作，即终结置换。

go标准库中使用的初始置换表和终结置换表如下：

```
// 初始置换表
var initialPermutation = [64]byte{
	6, 14, 22, 30, 38, 46, 54, 62,
	4, 12, 20, 28, 36, 44, 52, 60,
	2, 10, 18, 26, 34, 42, 50, 58,
	0, 8, 16, 24, 32, 40, 48, 56,
	7, 15, 23, 31, 39, 47, 55, 63,
	5, 13, 21, 29, 37, 45, 53, 61,
	3, 11, 19, 27, 35, 43, 51, 59,
	1, 9, 17, 25, 33, 41, 49, 57,
}

// 终结置换表
var finalPermutation = [64]byte{
	24, 56, 16, 48, 8, 40, 0, 32,
	25, 57, 17, 49, 9, 41, 1, 33,
	26, 58, 18, 50, 10, 42, 2, 34,
	27, 59, 19, 51, 11, 43, 3, 35,
	28, 60, 20, 52, 12, 44, 4, 36,
	29, 61, 21, 53, 13, 45, 5, 37,
	30, 62, 22, 54, 14, 46, 6, 38,
	31, 63, 23, 55, 15, 47, 7, 39,
}
```

### 子密钥的计算

DES初始密钥为64位，其中8位用于奇偶校验，实际密钥为56位，64位初始密钥经过PC-1密钥置换后，生成56位串。经PC-1置换后56位的串，分为左右两部分，各28位，分别左移1位，形成C0和D0，C0和D0合并成56位，经PC-2置换后生成48位子密钥K0。C0和D0分别左移1位，形成C1和D1，C1和D1合并成56位，经PC-2置换后生成子密钥K1。以此类推，直至生成子密钥K15。但注意每轮循环左移的位数，有如下规定：

```
var ksRotations = [16]uint8{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}
```

![](http://olgjbx93m.bkt.clouddn.com/2131235146206893dc337d83.jpg)

DES子密钥计算的代码如下：

```
// creates 16 56-bit subkeys from the original key
func (c *desCipher) generateSubkeys(keyBytes []byte) {
	// PC-1密钥置换，生成56位串
	key := binary.BigEndian.Uint64(keyBytes)
	permutedKey := permuteBlock(key, permutedChoice1[:])

	// 56位串分左右两部分，各28位
	leftRotations := ksRotate(uint32(permutedKey >> 28))
	rightRotations := ksRotate(uint32(permutedKey<<4) >> 4)

	// 生成子密钥
	for i := 0; i < 16; i++ {
		//合并左右两部分，之后PC-2置换
		pc2Input := uint64(leftRotations[i])<<28 | uint64(rightRotations[i])
		c.subkeys[i] = unpack(permuteBlock(pc2Input, permutedChoice2[:]))
	}
}
```


## 3DES

DES是一个经典的对称加密算法，但也缺陷明显，即56位的密钥安全性不足，已被证实可以在短时间内破解。为解决此问题，出现了3DES，也称Triple DES，3DES为DES向AES过渡的加密算法，它使用3条56位的密钥对数据进行三次加密。为了兼容普通的DES，3DES并没有直接使用 `加密->加密->加密` 的方式，而是采用了`加密->解密->加密` 的方式。

![](http://oscd4dgpc.bkt.clouddn.com/00f3fb2538ec90b80eb9.jpg)

当三重密钥均相同时，前两步相互抵消，相当于仅实现了一次加密，因此可实现对普通DES加密算法的兼容。

![](http://oscd4dgpc.bkt.clouddn.com/7de5d8dfcd0bd221596384.jpg)

### 解密

3DES解密过程，与加密过程相反，即逆序使用密钥。是以密钥3、密钥2、密钥1的顺序执行 `解密->加密->解密`。

![](http://oscd4dgpc.bkt.clouddn.com/WX20180216-083713.png)

相比DES，3DES因密钥长度变长，安全性有所提高，但其处理速度不高。因此又出现了AES加密算法，AES较于3DES速度更快、安全性更高。


