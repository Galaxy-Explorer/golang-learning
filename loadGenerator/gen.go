package loadGenerator

import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "golang_learning/helper/log"
    "golang_learning/loadGenerator/lib"
    "math"
    "sync/atomic"
    "time"
)

var logger = log.DLogger()

type myGenerator struct {
    caller      lib.Caller
    timeoutNS   time.Duration
    lps         uint32
    durationNS  time.Duration
    concurrency uint32
    tickets     lib.GoTickets
    ctx         context.Context
    cancelFunc  context.CancelFunc
    callCount   int64
    status      uint32
    resultCh    chan *lib.CallResult
}

func NewGenerator(pSet ParamSet) (lib.Generator, error) {
    logger.Infoln("New a load generator...")
    if err := pSet.Check(); err != nil {
        return nil, err
    }
    // 新建一个指针类型，该指针类型实现了Generator，所以可以返回
    gen := &myGenerator{
        caller:     pSet.Caller,
        timeoutNS:  pSet.TimeoutNS,
        lps:        pSet.LPS,
        durationNS: pSet.DurationNS,
        status:     lib.StatusOriginal,
        resultCh:   pSet.ResultCh,
    }
    if err := gen.init(); err != nil {
        return nil, err
    }
    return gen, nil
}

func (gen *myGenerator) init() error {
    var buf bytes.Buffer
    buf.WriteString("Initializing the load generator...")
    var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
    if total64 > math.MaxInt32 {
        total64 = math.MaxInt32
    }
    gen.concurrency = uint32(total64)
    tickets, err := lib.NewGoTickets(gen.concurrency)
    if err != nil {
        return err
    }
    gen.tickets = tickets
    buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
    logger.Infoln(buf.String())
    return nil
}

func (gen *myGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
    atomic.AddInt64(&gen.callCount, 1)
    if rawReq == nil {
        return &lib.RawResp{ID: -1, Err: errors.New("invalid raw request")}
    }
    start := time.Now().UnixNano()
    resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
    end := time.Now().UnixNano()
    elapsedTime := time.Duration(end - start)
    var rawResp lib.RawResp
    if err != nil {
        errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
        rawResp = lib.RawResp{
            ID:      rawReq.ID,
            Err:     errors.New(errMsg),
            Elapsed: elapsedTime,
        }
    } else {
        rawResp = lib.RawResp{
            ID:      rawReq.ID,
            Resp:    resp,
            Elapsed: elapsedTime}
    }
    return &rawResp
}

func (gen *myGenerator) asyncCall() {
    gen.tickets.Take()
    defer func() {
        if p := recover(); p != nil {
            err, ok := interface{}(p).(error)
            var errMsg string
            if ok {
                errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
            } else {
                errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
            }
            logger.Errorln(errMsg)
            result := &lib.CallResult{
                ID:   -1,
                Code: lib.RetCodeFatalCall,
                Msg:  errMsg,
            }
            gen.sendResult(result)
        }
        gen.tickets.Return()
    }()
    rawReq := gen.caller.BuildReq()
    var callStatus uint32
    timer := time.AfterFunc(gen.timeoutNS, func() {
        if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
            return
        }
        result := &lib.CallResult{
            ID:     rawReq.ID,
            Req:    rawReq,
            Code:   lib.RetCodeWarningCallTimeout,
            Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
            Elapse: gen.timeoutNS,
        }

        gen.sendResult(result)
    })
    rawResp := gen.callOne(&rawReq)
    if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
        return
    }
    timer.Stop()
    var result *lib.CallResult
    if rawResp.Err != nil {
        result = &lib.CallResult{
            ID:     rawResp.ID,
            Req:    rawReq,
            Code:   lib.RetCodeErrorCall,
            Msg:    rawResp.Err.Error(),
            Elapse: rawResp.Elapsed}
    } else {
        result = gen.caller.CheckResp(rawReq, *rawResp)
        result.Elapse = rawResp.Elapsed
    }
    gen.sendResult(result)
}

func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
    if atomic.LoadUint32(&gen.status) != lib.StatusStarted {
        gen.printIgnoredResult(result, "stopped load generator")
        return false
    }
    select {
    case gen.resultCh <- result:
        return true
    default:
        gen.printIgnoredResult(result, "full result channel")
        return false
    }

}

func (gen *myGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
    resultMsg := fmt.Sprintf(
        "ID=%d, Code=%d, Msg=%s, Elapse=%v",
        result.ID, result.Code, result.Msg, result.Elapse)
    logger.Warnf("Ignored result: %s. (cause: %s)\n", resultMsg, cause)
}

func (gen *myGenerator) prepareStop(ctxError error) {
    logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
    atomic.CompareAndSwapUint32(&gen.status, lib.StatusStarted, lib.StatusStopping)
    logger.Infof("Closing result channel...")
    close(gen.resultCh)
    atomic.StoreUint32(&gen.status, lib.StatusStopped)
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
    for {
        select {
        case <-gen.ctx.Done():
            gen.prepareStop(gen.ctx.Err())
            return
        default:
        }
        gen.asyncCall()
        if gen.lps > 0 {
            select {
            case <-throttle:
            case <-gen.ctx.Done():
                return
            }
        }
    }
}

func (gen *myGenerator) Start() bool {
    logger.Infoln("Starting load generator...")
    if !atomic.CompareAndSwapUint32(&gen.status, lib.StatusOriginal, lib.StatusStarting) {
        if !atomic.CompareAndSwapUint32(&gen.status, lib.StatusStopped, lib.StatusStarting) {
            return false
        }
    }

    var throttle <-chan time.Time
    if gen.lps > 0 {
        interval := time.Duration(1e9 / gen.lps)
        logger.Infof("Setting throttle (%v)...", interval)
        throttle = time.Tick(interval)
    }
    gen.ctx, gen.cancelFunc = context.WithTimeout(context.Background(), gen.durationNS)

    gen.callCount = 0

    atomic.StoreUint32(&gen.status, lib.StatusStarted)

    go func() {
        logger.Infoln("Generating loads...")
        gen.genLoad(throttle)
        logger.Infof("Stopped. (call count: %d)", gen.callCount)
    }()
    return true
}

func (gen *myGenerator) Stop() bool {
    if atomic.CompareAndSwapUint32(&gen.status, lib.StatusStarted, lib.StatusStopping) {
        return false
    }
    gen.cancelFunc()
    for {
        if atomic.LoadUint32(&gen.status) == lib.StatusStopped {
            break
        }
        time.Sleep(time.Microsecond)
    }
    return true

}

func (gen *myGenerator) Status() uint32 {
    return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
    return atomic.LoadInt64(&gen.callCount)
}
