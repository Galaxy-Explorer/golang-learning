package lib

import (
    "time"
)

type RawReq struct {
    ID  int64
    Req []byte
}

type RawResp struct {
    ID     int64
    Resp   []byte
    Err    error
    Elapse time.Duration
}

type RetCode int

const (
    RET_CODE_SUCCESS         RetCode = 0
    RET_CODE_WARNING_TIMEOUT         = 1001
    RET_CODE_ERROR_CALL              = 2001
    RET_CODE_ERROR_RESPONSE          = 2002
    RET_CODE_ERROR_CALEE             = 2003
    RET_CODE_FATAL_CALL              = 3001
)

func GetRetCodePlain(code RetCode) string {
    var codePlain string
    switch code {
    case RET_CODE_SUCCESS:
        codePlain = ""
    case RET_CODE_WARNING_TIMEOUT:
        codePlain = ""
    case RET_CODE_ERROR_CALL:
        codePlain = ""
    case RET_CODE_ERROR_RESPONSE:
        codePlain = ""
    case RET_CODE_ERROR_CALEE:
        codePlain = ""
    case RET_CODE_FATAL_CALL:
        codePlain = ""
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
    STATUS_ORIGINAL uint32 = 0
    STATUS_STARTING uint32 = 1
    STATUS_STARTED  uint32 = 2
    STATUS_STOPPING uint32 = 3
    STATUS_STOPPED  uint32 = 4
)

type Generator interface {
    Start() bool
    Stop() bool
    Status() uint32
    CallCount() int64
}
