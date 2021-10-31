## 1.总结几种 socket 粘包的解包方式: fix length/delimiter based/length field based frame decoder。尝试举例其应用、
​		因为 TCP 是面向连接的传输协议，TCP 传输的数据是以流的形式，而流数据是没有明确的开始结尾边界，所以 TCP 也没办法判断哪一段流属于一个消息。

* 粘包的主要原因：
  1. 发送方每次写入数据 < 套接字（Socket）缓冲区大小；
  2. 接收方读取套接字（Socket）缓冲区数据不够及时。

* 半包的主要原因：
  1. 发送方每次写入数据 > 套接字（Socket）缓冲区大小；
  2. 发送的数据大于协议的 MTU (Maximum Transmission Unit，最大传输单元)，因此必须拆包。

* 粘包和半包的解决方案有以下 3 种：
  1. 发送方和接收方规定固定大小的缓冲区，也就是发送和接收都使用固定大小的 byte[] 数组长度，当字符长度不够时使用空字符填充；
  2. 在 TCP 协议的基础上封装一层数据请求协议，即将数据包封装成数据头（存储数据正文大小）+ 数据正文的形式，这样在接收方就可以知道每个数据包的具体长度了，进而确定发送数据的具体边界解决半包和粘包的问题；
  3. 发送和接收数据的双方约定以特殊的字符结尾，比如以“\n”结尾，当收到该特殊字符就进行分帧以解决半包和粘包问题。

## 2.实现一个从 socket connection 中解码出 goim 协议的解码器。