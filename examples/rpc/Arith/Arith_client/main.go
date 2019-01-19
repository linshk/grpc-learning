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
    fmt.Printf("Arith: %d * %d = %d\n", args.A, args.B, reply)
    var quotient Quotient
    err = client.Call("Arith.Divide", args, &quotient)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d / %d = %d remains %d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}