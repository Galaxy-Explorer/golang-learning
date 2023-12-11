package testhelper

import (
    "errors"
    "fmt"
    "golang_learning/loadGenerator/lib"
    "testing"
    "time"
)

// 总是会有请求失败的问题
func TestHelper(t *testing.T) {
    serverAddr := "127.0.0.1:8080"
    server := NewTCPServer()
    err := server.Listen(serverAddr)
    if err != nil {
        logger.Errorf("Startup TCP server(%s)...\n", err)
    }

    conn := NewTcpComm(serverAddr)
    rawReq := conn.BuildReq()
    start := time.Now().UnixNano()
    bytes, err := conn.Call(rawReq.Req, 1000*time.Microsecond)
    if err != nil {
        logger.Errorf("request call is error:(%s)...\n", err)
    }
    end := time.Now().UnixNano()
    elapsedTime := time.Duration(end - start)

    var rawResp lib.RawResp
    if err != nil {
        errMsg := fmt.Sprintf("Call Error: %s.", err)
        rawResp = lib.RawResp{
            ID:      rawReq.ID,
            Err:     errors.New(errMsg),
            Elapsed: elapsedTime,
        }
    } else {
        rawResp = lib.RawResp{
            ID:      rawReq.ID,
            Resp:    bytes,
            Elapsed: elapsedTime}
    }

    resp := conn.CheckResp(rawReq, rawResp)

    t.Logf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n",
        resp.ID, resp.Code, resp.Msg, resp.Elapse)

}
