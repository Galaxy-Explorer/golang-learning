package channel

import "fmt"

// Send the sequence 2, 3, 4, ... to returned channel
func generate2() chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 2; ; i++ {
            ch <- i
        }
    }()
    return ch
}

// Filter out input values divisible by 'prime', send rest to returned channel
func filter2(in chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            if i := <-in; i%prime != 0 {
                out <- i
            }
        }
    }()
    return out
}
func sieve2(n int) chan int {
    out := make(chan int)
    stop := false
    go func(stop bool) {
        defer close(out)
        ch := generate2()
        for !stop {
            prime := <-ch
            if prime > n {
                stop = true
            }
            ch = filter2(ch, prime)
            if prime < n {
                out <- prime
            }

        }
    }(stop)
    return out
}
func Prime2() {
    primes := sieve2(100)
    for x := range primes {
        fmt.Println(x)
    }
}
