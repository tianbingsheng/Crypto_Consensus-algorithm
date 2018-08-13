# 007 DES、3DES 加密算法的调用

## DES加密常用的概念

* 加密模式

    * **ECB模式** 全称Electronic Codebook模式，译为电子密码本模式
    * **CBC模式** 全称Cipher Block Chaining模式，译为密文分组链接模式
    * **CFB模式** 全称Cipher FeedBack模式，译为密文反馈模式
    * **OFB模式** 全称Output Feedback模式，译为输出反馈模式。
    * **CTR模式** 全称Counter模式，译为计数器模式。

* 初始向量
    
    当加密第一个明文分组时，由于不存在 “前一个密文分组”，因此需要事先准备一个长度为一个分组的比特序列来代替 “前一个密文分组”，这个比特序列称为初始化向量(InitializationVector)，通常缩写为 IV。
    
* 填充方式
    
    当明文长度不为分组长度的整数倍时，需要在最后一个分组中填充一些数据使其凑满一个分组长度。
    
    * NoPadding

        API或算法本身不对数据进行处理，加密数据由加密双方约定填补算法。例如若对字符串数据进行加解密，可以补充\0或者空格，然后trim 

    * PKCS5Padding

        加密前：数据字节长度对8取余，余数为m，若m>0,则补足8-m个字节，字节数值为8-m，即差几个字节就补几个字节，字节数值即为补充的字节数，若为0则补充8个字节的8 
    
        解密后：取最后一个字节，值为m，则从数据尾部删除m个字节，剩余数据即为加密前的原文。
        
        加密字符串为为AAA，则补位为AAA55555;加密字符串为BBBBBB，则补位为BBBBBB22；加密字符串为CCCCCCCC，则补位为CCCCCCCC88888888。
        
    * PKCS7Padding
    
        PKCS7Padding 的填充方式和PKCS5Padding 填充方式一样。只是加密块的字节数不同。PKCS5Padding明确定义了加密块是8字节，PKCS7Padding加密快可以是1-255之间。

## 实现 DES 加密和解密

采用3DES、CBC模式、pkcs5padding，初始向量用key充当；另外，对于zero padding，还得约定好，对于数据长度刚好是block size的整数倍时，是否需要额外填充。

### crypto/des包

Go中crypto/des包实现了 Data Encryption Standard (DES) and the Triple Data Encryption Algorithm (TDEA)。[查看该包文档](https://golang.org/pkg/crypto/des/)，发现相当简单：定义了DES块大小（8bytes），定义了一个KeySizeError。另外定义了两个我们需要特别关注的函数，即

```
    func NewCipher(key []byte) (cipher.Block, error)

    func NewTripleDESCipher(key []byte) (cipher.Block, error)
```

他们都是用来获得一个cipher.Block。从名字可以很容易知道，DES使用NewCipher，3DES使用NewTripleDESCipher。参数都是密钥（key）

### 加密

使用DES加密 `（des.NewCipher）` ，加密模式为CBC `（cipher.NewCBCEncrypter(block, key)） `，填充方式 `PKCS5Padding`

```
func DesEncrypt(origData, key []byte) ([]byte, error) {

    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    origData = PKCS5Padding(origData, block.BlockSize())
    blockMode := cipher.NewCBCEncrypter(block, key)
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}
```

### 解密

```
func DesDecrypt(crypted, key []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockMode := cipher.NewCBCDecrypter(block, key)
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS5UnPadding(origData)
    return origData, nil
}
```

### des加密的详细代码

```
package main

import (
    "bytes"
    "crypto/cipher"
    "crypto/des"
    "encoding/base64"
    "fmt"
)

func main() {
	key := []byte("12345678")
    result, err := DesEncrypt([]byte("hello world"), key)
    if err != nil {
        panic(err)
    }
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	
	
    origData, err := DesDecrypt(result, key)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(origData))
}

func DesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    origData = PKCS5Padding(origData, block.BlockSize())
    blockMode := cipher.NewCBCEncrypter(block, key)
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockMode := cipher.NewCBCDecrypter(block, key)
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS5UnPadding(origData)
    return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
```

## 3DES 加密及解密

### 加密

对比DES，发现只是换了NewTripleDESCipher。不过，需要注意的是，密钥长度必须24byte，否则直接返回错误。

```
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    origData = PKCS5Padding(origData, block.BlockSize())
    blockMode := cipher.NewCBCEncrypter(block, key[:8])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}
```

### 主函数

```
func main() {
    // 3DES加解密
	key := []byte("123456789012345678901234")
    result, err := TripleDesEncrypt([]byte("hello world"), key)
    if err != nil {
        panic(err)
    }
    fmt.Println(base64.StdEncoding.EncodeToString(result))
    origData, err := TripleDesDecrypt(result, key)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(origData))
}
```

如果我们把主函数中 `key` 改为25位的 `1234567890123456789012345` 运行 `go run 3des.go`，提示出现如下错误： 

```
go run 3des.go
panic: crypto/des: invalid key size 25

goroutine 1 [running]:
```

### 详细代码

```
package main

import (
    "bytes"
    "crypto/cipher"
    "crypto/des"
    "encoding/base64"
    "fmt"
)

func main() {
    // 3DES加解密
	key := []byte("123456789012345678901234")
    result, err := TripleDesEncrypt([]byte("hello world"), key)
    if err != nil {
        panic(err)
    }
    fmt.Println(base64.StdEncoding.EncodeToString(result))
    origData, err := TripleDesDecrypt(result, key)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(origData))
}


// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    origData = PKCS5Padding(origData, block.BlockSize())
    blockMode := cipher.NewCBCEncrypter(block, key[:8])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    blockMode := cipher.NewCBCDecrypter(block, key[:8])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = PKCS5UnPadding(origData)
    return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
```


