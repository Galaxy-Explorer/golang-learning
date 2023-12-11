package lib

import "time"

type RawReq struct {
    ID  int64
    Req []byte
}

type RawResp struct {
    ID      int64
    Resp    []byte
    Err     error
    Elapsed time.Duration
}

type RetCode int

const (
    RetCodeSuccess            RetCode = 0    // 成功。
    RetCodeWarningCallTimeout         = 1001 // 调用超时警告。
    RetCodeErrorCall                  = 2001 // 调用错误。
    RetCodeErrorResponse              = 2002 // 响应内容错误。
    RetCodeErrorCalee                 = 2003 // 被调用方（被测软件）的内部错误。
    RetCodeFatalCall                  = 3001 // 调用过程中发生了致命错误！
)

func GetRetCodePlain(code RetCode) string {
    var codePlain string
    switch code {
    case RetCodeSuccess:
        codePlain = "Success"
    case RetCodeWarningCallTimeout:
        codePlain = "Call Timeout Warning"
    case RetCodeErrorCall:
        codePlain = "Call Error"
    case RetCodeErrorResponse:
        codePlain = "Response Error"
    case RetCodeErrorCalee:
        codePlain = "Callee Error"
    case RetCodeFatalCall:
        codePlain = "Call Fatal Error"
    default:
        codePlain = "Unknown result code"
    }
    return codePlain
}

type CallResult struct {
    ID     int64
    Req    RawReq
    Resp   RawResp
    Code   RetCode
    Msg    string
    Elapse time.Duration
}

const (
    // StatusOriginal 代表原始。
    StatusOriginal uint32 = 0
    // StatusStarting 代表正在启动。
    StatusStarting uint32 = 1
    // StatusStarted 代表已启动。
    StatusStarted uint32 = 2
    // StatusStopping 代表正在停止。
    StatusStopping uint32 = 3
    // StatusStopped 代表已停止。
    StatusStopped uint32 = 4
)

type Generator interface {
    Start() bool
    Stop() bool
    Status() uint32
    CallCount() int64
}
