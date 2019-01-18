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
