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
    resultCh    chan *lib.CallResult
    concurrency uint32
    tickets     lib.GoTickets
    stopSign    chan struct{}
    ctx         context.Context
    cancelFunc  context.CancelFunc
    callCount   int64
    status      uint32
}

func NewGenerator(pset ParamSet) (lib.Generator, error) {

    logger.Infoln("New a load generator...")
    if err := pset.Check(); err != nil {
        return nil, err
    }
    gen := &myGenerator{
        caller:     pset.Caller,
        timeoutNS:  pset.TimeoutNS,
        lps:        pset.LPS,
        durationNS: pset.DurationNS,
        status:     lib.STATUS_ORIGINAL,
        resultCh:   pset.ResultCh,
    }
    if err := gen.init(); err != nil {
        return nil, err
    }
    return gen, nil
}

func (gen *myGenerator) init() error {
    var buf bytes.Buffer
    buf.WriteString("Initializing the load generator...")
    // 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
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
        return &lib.RawResp{ID: -1, Err: errors.New("Invalid raw request.")}
    }
    start := time.Now().UnixNano()
    resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
    end := time.Now().UnixNano()
    elapsedTime := time.Duration(end - start)
    var rawResp lib.RawResp
    if err != nil {
        errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
        rawResp = lib.RawResp{
            ID:     rawReq.ID,
            Err:    errors.New(errMsg),
            Elapse: elapsedTime}
    } else {
        rawResp = lib.RawResp{
            ID:     rawReq.ID,
            Resp:   resp,
            Elapse: elapsedTime}
    }
    return &rawResp
}

func (gen *myGenerator) asyncCall() {
    gen.tickets.Take()
    go func() {
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
                    Code: lib.RET_CODE_FATAL_CALL,
                    Msg:  errMsg}
                gen.sendResult(result)
            }
            gen.tickets.Return()
        }()
        rawReq := gen.caller.BuildReq()
        // 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
        var callStatus uint32
        timer := time.AfterFunc(gen.timeoutNS, func() {
            if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
                return
            }
            result := &lib.CallResult{
                ID:     rawReq.ID,
                Req:    rawReq,
                Code:   lib.RET_CODE_WARNING_TIMEOUT,
                Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
                Elapse: gen.timeoutNS,
            }
            gen.sendResult(result)
        })

        // callOne 这个方法里，rawResp.ID = rawReq.ID
        rawResp := gen.callOne(&rawReq)
        // 调用完成
        if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
            return
        }
        // 调用超时
        timer.Stop()
        var result *lib.CallResult
        // 调用出错
        if rawResp.Err != nil {
            result = &lib.CallResult{
                ID:     rawResp.ID,
                Req:    rawReq,
                Code:   lib.RET_CODE_ERROR_CALL,
                Msg:    rawResp.Err.Error(),
                Elapse: rawResp.Elapse}
        } else {
            // 对正常result的检查
            result = gen.caller.CheckResp(rawReq, *rawResp)
            result.Elapse = rawResp.Elapse
        }
        gen.sendResult(result)
    }()
}

func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
    if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
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

func (gen *myGenerator) prepareToStop(ctxError error) {
    logger.Infof("Prepare to stop load generator (cause: %s)...", ctxError)
    atomic.CompareAndSwapUint32(
        &gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
    logger.Infof("Closing result channel...")
    close(gen.resultCh)
    atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
    for {
        select {
        case <-gen.ctx.Done():
            gen.prepareToStop(gen.ctx.Err())
            return
        default:
        }
        gen.asyncCall()
        if gen.lps > 0 {
            select {
            case <-throttle:
            case <-gen.ctx.Done():
                gen.prepareToStop(gen.ctx.Err())
                return
            }
        }
    }
}

func (gen *myGenerator) Start() bool {
    logger.Infoln("Starting load generator...")
    // 检查是否具备可启动的状态，顺便设置状态为正在启动
    if !atomic.CompareAndSwapUint32(
        &gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
        if !atomic.CompareAndSwapUint32(
            &gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
            return false
        }
    }

    // 设定节流阀。
    var throttle <-chan time.Time
    if gen.lps > 0 {
        interval := time.Duration(1e9 / gen.lps)
        logger.Infof("Setting throttle (%v)...", interval)
        throttle = time.Tick(interval)
    }

    // 初始化上下文和取消函数。
    gen.ctx, gen.cancelFunc = context.WithTimeout(
        context.Background(), gen.durationNS)

    // 初始化调用计数。
    gen.callCount = 0

    // 设置状态为已启动。
    atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

    go func() {
        // 生成并发送载荷。
        logger.Infoln("Generating loads...")
        gen.genLoad(throttle)
        logger.Infof("Stopped. (call count: %d)", gen.callCount)
    }()
    return true
}

func (gen *myGenerator) Stop() bool {
    if !atomic.CompareAndSwapUint32(
        &gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) {
        return false
    }
    gen.cancelFunc()
    for {
        if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED {
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
