package testhelper

import (
    "encoding/json"
    "fmt"
    lib "golang_learning/loadGenerator/lib"
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

func (comm *TcpComm) BuildReq() lib.RawReq {
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
    rawReq := lib.RawReq{ID: ID, Req: bytes}

    return rawReq
}

func (comm *TcpComm) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
    conn, err := net.DialTimeout("tcp", comm.addr, timeoutNS)
    if err != nil {
        return nil, err
    }
    _, err = write(conn, req, DELIM)
    if err != nil {
        return nil, err
    }

    return read(conn, DELIM)

}

func (comm *TcpComm) CheckResp(rawReq lib.RawReq, rawResp lib.RawResp) *lib.CallResult {
    var commResult lib.CallResult
    commResult.ID = rawResp.ID
    commResult.Req = rawReq
    commResult.Resp = rawResp
    var sreq ServerReq
    err := json.Unmarshal(rawReq.Req, &sreq)
    if err != nil {
        commResult.Code = lib.RetCodeFatalCall
        commResult.Msg =
            fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
        return &commResult
    }
    var sresp ServerResp
    err = json.Unmarshal(rawResp.Resp, &sresp)
    if err != nil {
        commResult.Code = lib.RetCodeErrorResponse
        commResult.Msg =
            fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
        return &commResult
    }
    if sresp.ID != sreq.ID {
        commResult.Code = lib.RetCodeErrorResponse
        commResult.Msg =
            fmt.Sprintf("Inconsistent raw id! (%d != %d)\n", rawReq.ID, rawResp.ID)
        return &commResult
    }
    if sresp.Err != nil {
        commResult.Code = lib.RetCodeErrorCalee
        commResult.Msg =
            fmt.Sprintf("Abnormal server: %s!\n", sresp.Err)
        return &commResult
    }
    if sresp.Result != op(sreq.Operands, sreq.Operator) {
        commResult.Code = lib.RetCodeErrorResponse
        commResult.Msg =
            fmt.Sprintf(
                "Incorrect result: %s!\n",
                genFormula(sreq.Operands, sreq.Operator, sresp.Result, false))
        return &commResult
    }
    commResult.Code = lib.RetCodeSuccess
    commResult.Msg = fmt.Sprintf("Success. (%s)", sresp.Formula)
    return &commResult
}
