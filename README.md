



## gRPC简介

gRPC是一个高性能、通用的开源RPC框架，其由Google主要面向移动应用开发并基于HTTP/2协议标准而设计，基于ProtoBuf(Protocol Buffers)序列化协议开发，且支持众多开发语言。 gRPC提供了一种简单的方法来精确地定义服务和为iOS、Android和后台支持服务自动生成可靠性很强的客户端功能库。 客户端充分利用高级流和链接功能，从而有助于节省带宽、降低TCP连接次数、节省CPU使用和电池寿命。 



> RPC(Remote Procedure Call，远程过程调用)是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络细节的应用程序通信协议。RPC协议构建于TCP或UDP,或者是HTTP上。允许开发者直接调用另一台服务器上的程序，而开发者无需另外的为这个调用过程编写网络通信相关代码，使得开发网络分布式程序在内的应用程序更加容易
>
> RPC采用客户端-服务器端的工作模式，请求程序就是一个客户端，而服务提供程序就是一个服务器端。当执行一个远程过程调用时，客户端程序首先先发送一个带有参数的调用信息到服务端，然后等待服务端响应。在服务端，服务进程保持睡眠状态直到客户端的调用信息到达。当一个调用信息到达时，服务端获得进程参数，计算出结果，并向客户端发送应答信息。然后等待下一个调用。



## RPC实践

在认识gRPC之前，我们先来看看一个简单的RPC服务端和客户端具体是怎样的。 下面以Go语言的net/rpc包为例实现一个提供乘法和除法运算服务的RPC服务端和RPC客户端。



服务端上的RPC服务可简单理解为一组可被客户端调用的方法，这些方法的形式必须如下：

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

- 该方法所属的类型T必须是导出的（即类型名首字母大写）
- 方法是导出的（即方法名首字母大写）
- 方法带有两个参数，T1和T2（T2必须是指针类型，用于写入该方法的返回结果），且T1和T2都是导出类型或内建类型
- 方法返回类型为error



**服务端实现：**

```go
package main

import (
    "net"
    "net/http"
    "net/rpc"
    "log"
    "errors"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}


func main() {
    arith := new(Arith)
    rpc.Register(arith)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":2345")
    if e != nil {
        log.Fatal("listen error:", e)
    }
    log.Println("rpcserver listening at 127.0.0.1:2345")
    http.Serve(l, nil)
}

```



上面的代码主要有两部分：

**服务的定义和实现：**

Arith将作为服务名，该服务包含两个方法可被远程调用，Multiply和Divide。这两个函数都遵循rpc包要求的形式，通过结构体来传送1个或多个函数的参数，然后将函数调用的结果写入用于存放响应内容的结构体中。

**服务的监听和处理：**

- rpc.Register方法注册Arith的一个实例，使得Arith服务对外可调用
- HandleHTTP方法用于注册一个HTTP handler来处理RPC消息
- net.Listen设置使用的协议及端口，并返回一个Listener对象l用于监听新建立的连接
- http.Serve启动一个http server处理Listener对象l接收到的连接



**客户端实现：**

```go
package main

import (
    "fmt"
    "net/rpc"
    "log"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}


var serverAddress string = "127.0.0.1"

func main() {
    client, err := rpc.DialHTTP("tcp", serverAddress + ":2345")
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // Synchronous call
    args := &Args{9,4}
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d * %d=%d\n", args.A, args.B, reply)
    var quotient Quotient
    err = client.Call("Arith.Divide", args, &quotient)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d / %d = %d remains %d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}
```



客户端的代码主要有两部分：

**定义和服务端Arith服务相同或相似的数据结构：**

客户端定义了一样的结构体Args和Quotient，用于调用服务端Arith服务的Multiply方法和Divide方法。事实上，客户端Args和Quotient的定义也不一定要完全一致，但不能有冲突。

例如，以下定义也是可以的

```go
type Args struct {
    A, B, C int
}

type Args struct {
    A, B *int
}
```

但以下定义都是不可以的

```go
type Args struct {
    A int; B float
}

type Args struct {
    C, D int
}
```



**连接服务端及服务调用：**

- 通过rpc.DialHTTP与服务端建立连接并获取Client对象
- 通过client.Call调用Arith服务的方法，Call方法第一个参数的形式为“服务名.方法名”，服务名和方法名都必须和服务端的定义一致



RPC客户端只要遵循服务端的服务定义的格式就可以像调用本地方法一样调用远程的服务。



**代码运行:**

上述代码已放在github仓库，可直接下载运行：

```bash
$ mkdir -p $GOPATH/src/github.com/linshk
$ cd $GOPATH/src/github.com/linshk
$ git clone https://github.com/linshk/grpc-learning
$ cd grpc-learning
# 以后台方式运行服务端
$ go run rpcserver/main.go &
$ go run rpcclient/main.go
# 客户端输出：
# Arith: 9 * 4=36
# Arith: 9 / 4 = 2 remains 1
```





## gRPC helloworld

进入examples目录查看示例项目

```sh
$ cd $GOPATH/src/google.golang.org/grpc/examples/helloworld
```

helloword目录结构：

- helloworld
  - helloworld
    - helloworld.pb.go（protoc-gen-go根据helloworld.proto的服务定义自动生成的服务端客户端代码）
    - helloworld.proto （Greeter服务（gRPC服务）的定义）
  - greeter_client
    - main.go（借助helloworld.pb.go构建的gRPC客户端）
  - greeter_server
    - main.go（借助helloworld.pb.go构建的gRPC服务端）
  - mock_helloworld
    - hw_mock.go（MockGen自动生成的mock服务器代码）
    - hw_mock_test.go



首先查看helloworld.proto

```protobuf
syntax = "proto3";

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";

package helloworld;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```



前面我们手动构建RPC服务端代码时需要先定义服务的类型（Arith），函数参数及计算结果的类型（Args为参数的类型，Int为Multiply的结果类型，Quotient为Divide方法的结果类型），服务可供调用的方法（Multiply和Quotient），而Greeter服务的定义也是如此：



- 服务的定义：

```protobuf
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}
```

服务名为Greeter，服务可供调用的方法为SayHello，且SayHello方法的形式与前面的Multiply和Divide有所不同，SayHello返回的返回值即是方法的调用结果。



- 参数类型及结果类型的定义：

```protobuf
// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```



接着再看由该服务定义自动生成的代码helloworld.pb.go

```go
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: helloworld.proto

package helloworld

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type HelloRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_helloworld_71e208cbdc16936b, []int{0}
}
func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
}
func (dst *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(dst, src)
}
func (m *HelloRequest) XXX_Size() int {
	return xxx_messageInfo_HelloRequest.Size(m)
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type HelloReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReply) Reset()         { *m = HelloReply{} }
func (m *HelloReply) String() string { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()    {}
func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_helloworld_71e208cbdc16936b, []int{1}
}
func (m *HelloReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloReply.Unmarshal(m, b)
}
func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
}
func (dst *HelloReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReply.Merge(dst, src)
}
func (m *HelloReply) XXX_Size() int {
	return xxx_messageInfo_HelloReply.Size(m)
}
func (m *HelloReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReply.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReply proto.InternalMessageInfo

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServer is the server API for Greeter service.
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld.proto",
}

func init() { proto.RegisterFile("helloworld.proto", fileDescriptor_helloworld_71e208cbdc16936b) }

var fileDescriptor_helloworld_71e208cbdc16936b = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x48, 0xcd, 0xc9,
	0xc9, 0x2f, 0xcf, 0x2f, 0xca, 0x49, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x28, 0x29, 0x71, 0xf1, 0x78, 0x80, 0x78, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42,
	0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92,
	0x1a, 0x17, 0x17, 0x54, 0x4d, 0x41, 0x4e, 0xa5, 0x90, 0x04, 0x17, 0x7b, 0x6e, 0x6a, 0x71, 0x71,
	0x62, 0x3a, 0x4c, 0x11, 0x8c, 0x6b, 0xe4, 0xc9, 0xc5, 0xee, 0x5e, 0x94, 0x9a, 0x5a, 0x92, 0x5a,
	0x24, 0x64, 0xc7, 0xc5, 0x11, 0x9c, 0x58, 0x09, 0xd6, 0x25, 0x24, 0xa1, 0x87, 0xe4, 0x02, 0x64,
	0xcb, 0xa4, 0xc4, 0xb0, 0xc8, 0x14, 0xe4, 0x54, 0x2a, 0x31, 0x38, 0x19, 0x70, 0x49, 0x67, 0xe6,
	0xeb, 0xa5, 0x17, 0x15, 0x24, 0xeb, 0xa5, 0x56, 0x24, 0xe6, 0x16, 0xe4, 0xa4, 0x16, 0x23, 0xa9,
	0x75, 0xe2, 0x07, 0x2b, 0x0e, 0x07, 0xb1, 0x03, 0x40, 0x5e, 0x0a, 0x60, 0x4c, 0x62, 0x03, 0xfb,
	0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x0f, 0xb7, 0xcd, 0xf2, 0xef, 0x00, 0x00, 0x00,
}
```



helloworld.pb.go已实现的内容主要有：

- 参数类型和结果类型的定义及相关方法的生成（HelloRequest和HelloReply）

如HelloRequest类型的定义及方法

```go
// The request message containing the user's name.
type HelloRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
func (m *HelloRequest) Reset() 
func (m *HelloRequest) String() string 
func (*HelloRequest) ProtoMessage()    
func (*HelloRequest) Descriptor() ([]byte, []int) 
func (m *HelloRequest) XXX_Unmarshal(b []byte) error 
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) 
func (dst *HelloRequest) XXX_Merge(src proto.Message) 
func (m *HelloRequest) XXX_Size() int 
func (m *HelloRequest) XXX_DiscardUnknown() 
func (m *HelloRequest) GetName() string
```

HelloRequest自动生成的方法



- 类型的注册

```go
func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}
```



- Greeter服务的服务端和客户端接口及接口实现

客户端：

```go
// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}



```

GreeterClient接口包含了Greeter服务中的可调用方法（即SayHello），greeterClient通过grpc.ClientConn.Invoke来实现对SayHello的调用。



服务端：

```go
// GreeterServer is the server API for Greeter service.
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}
```

GreeterServer接口也包含了Greeter服务中的可调用方法（即SayHello）。RegisterGreeterServer方法将一个GreeterServer实例注册到grpc.Server中，完成服务的注册。  _Greeter_SayHello_Handler用于处理客户端的请求，而GreeterServer接口并没有实现，因为我们需要自己实现GreeterServer接口来实现我们的Greeter服务。

从代码中还可看出Greeter服务的路径（方法全称）表示为："/helloworld.Greeter/SayHello"。



- 服务定义文件的注册（将服务定义文件与其压缩后的描述对象映射起来）

```go
func init() { proto.RegisterFile("helloworld.proto", fileDescriptor_helloworld_71e208cbdc16936b) }
```



借助上面自动生成的helloworld.pb.go，我们将可以很方便的实现我们自己的服务端和客户端，即greeter_server/main.go和greeter_client/main.go中的实现。



## gRPC实践

上面我们已经手动创建了一个简单的RPC服务端和客户端，接下来我们借助工具更方便地自动生成Arith服务的gRPC服务端和客户端的代码。



### 准备工作

- Go >= 1.9

  可参考[Go环境搭建](https://blog.csdn.net/linshk_ver18/article/details/82872634)

- 安装cpp版本的protobuf（Protocol Buffers的简称），[下载地址](https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-cpp-3.6.1.tar.gz)

  以下安装方法为默认安装，若需自定义安装可以参考[linux 下安装protobuf ](https://blog.csdn.net/xiexievv/article/details/47396725)

```sh
# 下载完成并解压后进入解压生成的文件夹
$ cd yourpath/to/protobuf
$ ./configure
$ sudo make
$ sudo make install
# 验证是否安装成功
$ protoc --help
# 若有如下报错
# protoc: error while loading shared libraries: libprotoc.so.17: cannot open shared object file: No such file or directory
# 则只需再设置一下环境变量
$ sudo echo "export LD_LIBRARY_PATH=/usr/local/lib" >> ~/.bashrc
$ source ~/.bashrc
```



- 安装grpc

```bash
$ go get -u google.golang.org/grpc
```

该安装方式显然需要翻*，也可以通过以下方式来安装

```bash
$ go get -d github.com/grpc/grpc-go
# 创建所需的文件夹
$ mkdir -p $GOPATH/src/google.golang.org/
$ mv $GOPATH/src/github.com/grpc/grpc-go 
$ GOPATH/src/google.golang.org/grpc
```

更详细的内容可参考[Get Golang Packages on Golang.org in China](https://github.com/northbright/Notes/blob/master/Golang/china/get-golang-packages-on-golang-org-in-china.md)



- 安装Go语言的代码生成器（用于自动生成gRPC服务端和客户端）

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```



