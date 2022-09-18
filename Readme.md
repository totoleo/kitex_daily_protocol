# 介绍
为了实现对 kitex 的通信内容的解析。

# 前置工作
1. 安装 thriftgo 和 kitex
```
go install github.com/cloudwego/thriftgo@latest
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```
2. 生成必要的 idl 结构体
```
kitex --module exp -thrift frugal_tag item.thrift
//或者完全使用 frugral 的 jit 编码也可以，可以参考 frugal 的文档进行
```

# 运行
`go run .`