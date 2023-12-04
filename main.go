package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    helper "golang_learning/loadGenerator/testhelper"
    "math/rand"
    "net"
    "time"
)

func read(conn net.Conn, delim byte) ([]byte, error) {
    reader := bufio.NewReader(conn)
    content, err := reader.ReadBytes(delim)
    if err != nil {
        return nil, err
    }

    return content, nil
}

func write(conn net.Conn, content []byte, delim byte) (int, error) {
    writer := bufio.NewWriter(conn)
    n, err := writer.Write(content)
    if err == nil {
        writer.WriteByte(delim)
    }
    if err == nil {
        err = writer.Flush()
    }
    return n, nil
}

var operators = []string{"+", "-", "*", "/"}

func main() {
    server := helper.NewTCPServer()
    defer server.Close()
    serverAddr := "127.0.0.1:8080"
    fmt.Printf("Startup TCP server(%s)...\n", serverAddr)
    err := server.Listen(serverAddr)
    if err != nil {
        fmt.Printf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
    }
    go func() {
        id := int64(1)
        sreq := helper.ServerReq{
            ID:       id,
            Operands: []int{int(rand.Int31n(1000) + 1), int(rand.Int31n(1000) + 1)},
            Operator: func() string { return operators[rand.Int31n(100)%4] }(),
        }
        bytes, _ := json.Marshal(sreq)
        conn, _ := net.DialTimeout("tcp", serverAddr, time.Second)
        write(conn, bytes, '\n')
        time.Sleep(time.Second * 10)
        fmt.Println("test1")
    }()

    go func() {
        conn, _ := net.DialTimeout("tcp", serverAddr, time.Second)
        fmt.Println("test2")
        result1, _ := read(conn, '\n')
        fmt.Println(string(result1))
    }()

    time.Sleep(time.Second * 100)
}
