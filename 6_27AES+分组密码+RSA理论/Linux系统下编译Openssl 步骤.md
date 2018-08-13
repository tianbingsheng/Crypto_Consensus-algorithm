#  Linux系统下编译Openssl 步骤:

##1.提前准备工作，作者提前准备的Vmware下安装的Centos7操作系统

去[https://www.openssl.org](https://www.openssl.org)网站下载openssl源码，作者下载的openssl-1.1.0h.tar版本

解压openssl开发包文件

 tar -xzf openssl-1.1.0c.tar.gz

###2.在终端下操作如下

进入刚才解压的文件夹，如图所示

![](http://om1c35wrq.bkt.clouddn.com/2018_4_4.png)



-设定Openssl 安装路径，( **--prefix** )参数为欲安装之目录，执行如下命令：

​     ./config --prefix=/usr/local/openssl

###3.执行命令./config -t

###4.执行make命令，编译Openssl，编译需要等待一定的时间。

###5、执行make install，安装 Openssl，安装也需要一定的时间。

安装完成后,openssl 会被安装到/usr/local/openssl 目录，包括头文件目录 

![](http://om1c35wrq.bkt.clouddn.com/1522819297%281%29.png)

进入lib后，libssl.a和libcrypto.a则为编译后的库文件