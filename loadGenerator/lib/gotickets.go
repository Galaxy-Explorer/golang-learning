package lib

import (
    "errors"
    "fmt"
)

type GoTickets interface {
    Take()
    Return()
    Total() uint32
    Active() bool
    Remainder() uint32
}

type myGoTickets struct {
    total    uint32
    ticketCh chan struct{}
    active   bool
}

func NewGoTickets(total uint32) (GoTickets, error) {
    gt := myGoTickets{}
    if !gt.init(total) {
        errMsg := fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
        return nil, errors.New(errMsg)
    }

    return &gt, nil
}

func (gt *myGoTickets) init(total uint32) bool {
    // 判断票池是否已经初始化
    if gt.active {
        return false
    }
    // 票池的容量是否大于0
    if total == 0 {
        return false
    }

    // 初始化票池的缓冲管道
    ch := make(chan struct{}, total)
    n := int(total)
    for i := 0; i < n; i++ {
        ch <- struct{}{}
    }

    gt.ticketCh = ch
    gt.total = total
    gt.active = true
    return true
}

func (gt *myGoTickets) Take() {
    <-gt.ticketCh
}

func (gt *myGoTickets) Return() {
    gt.ticketCh <- struct{}{}
}

func (gt *myGoTickets) Active() bool {
    return gt.active
}

func (gt *myGoTickets) Total() uint32 {
    return gt.total
}

func (gt *myGoTickets) Remainder() uint32 {
    return uint32(len(gt.ticketCh))
}
