package testhelper

import (
    "encoding/json"
    . "golang_learning/loadGenerator/lib"
    "math/rand"
    "net"
    "time"
)

const (
    DELIM = '\n'
)

var operators = [4]string{"+", "-", "*", "/"}

type TcpComm struct {
    addr string
}

func NewTcpComm(addr string) *TcpComm {
    return &TcpComm{addr: addr}
}

func (comm *TcpComm) BuildReq() RawReq {
    ID := time.Now().UnixNano()

    sreq := ServerReq{
        ID:       ID,
        Operands: []int{int(rand.Int31n(1000) + 1), int(rand.Int31n(1000) + 1)},
        Operator: operators[rand.Int31n(100)%4],
    }

    bytes, err := json.Marshal(sreq)
    if err != nil {
        panic(err)
    }
    rawReq := RawReq{ID: ID, Req: bytes}

    return rawReq
}

func (comm *TcpComm) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
    conn, err := net.DialTimeout("TCP", comm.addr, timeoutNS)
    if err != nil {
        return nil, err
    }
    _, err = write(conn, req, DELIM)
    if err != nil {
        return nil, err
    }

    return read(conn, DELIM)

}

func (comm *TcpComm) CheckResp(req RawReq, resp RawResp) *CallResult {
    var commResult CallResult

    return &commResult
}
