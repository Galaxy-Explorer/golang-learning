package loadGenerator

import (
    loadgenlib "golang_learning/loadGenerator/lib"
    helper "golang_learning/loadGenerator/testhelper"
    "testing"
    "time"
)

func TestParameterCheck(t *testing.T) {
    t.Logf("Check Parameter Set...\n")
    serverAddr := "127.0.0.1:8080"
    pSet := ParamSet{
        Caller:     helper.NewTcpComm(serverAddr),
        TimeoutNS:  50 * time.Millisecond,
        LPS:        uint32(1000),
        DurationNS: 10 * time.Second,
        ResultCh:   make(chan *loadgenlib.CallResult, 50),
    }
    err := pSet.Check()
    if err != nil {
        t.Logf("Initialize load generator error")
    } else {
        t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
            pSet.TimeoutNS, pSet.LPS, pSet.DurationNS)
    }

}
