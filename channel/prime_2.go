package channel

import "fmt"

// Send the sequence 2, 3, 4, ... to returned channel
func generate() chan int {
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
func filter(in chan int, prime int) chan int {
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
func sieve(n int) chan int {
    out := make(chan int)
    stop := false
    go func(stop bool) {
        defer close(out)
        ch := generate()
        for !stop {
            prime := <-ch
            if prime > n {
                stop = true
            }
            ch = filter(ch, prime)
            if prime < n {
                out <- prime
            }

        }
    }(stop)
    return out
}
func Prime2() {
    primes := sieve(100)
    for x := range primes {
        fmt.Println(x)
    }
}
