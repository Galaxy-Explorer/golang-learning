package lib

import "time"

type RawReq struct {
    ID  int64
    Req []byte
}

type RawResp struct {
    ID      int64
    Resp    []byte
    ErrMsg  error
    Elapsed time.Duration
}

type CallResult struct {
    ID      int64
    RawReq  RawReq
    RawResp RawResp
}
