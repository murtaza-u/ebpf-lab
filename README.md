# eBPF Lab

## Drop all packets on a specific port

> `eBPF` program to drop the TCP packets on a port (def: 4040). The port
> number is configurable from the userspace.

![](./ethernet-frame.png)

```sh
GOPROXY=direct go install github.com/murtaza-u/ebpf-lab/cmd/drop@latest
drop -h

# remember to run with sudo and make necessary changes to the
# interface and port.
sudo drop --interface wlp1s0 --port 80
```

Demo: [Link](https://imgur.com/PISJUlN)

## Concurrency in Go

```go
package main

import "fmt"

func main() {
    cnp := make(chan func(), 10)
    for i := 0; i < 4; i++ {
        go func() {
            for f := range cnp {
                f()
            }
        }()
    }
    cnp <- func() {
        fmt.Println("HERE1")
    }
    fmt.Println("Hello")
}
```

1. Explaining how the highlighted constructs work?

* We are first creating a *buffered* channel (capacity 10) of functions
  that take in no arguments and returns nothing.
* Then running a loop 4 times (i = 0 to 3) wherein in each iteration we
  are spawning a `goroutine`.
* Each `goroutine` iterates over all the funcs inside the channel till
  the channel is closed.
* We are then putting a function, that prints "Here1" on stdout, into
  the buffered channel.
* Lastly, we print "Hello" in the *main* goroutine.

The code above outputs "Hello" to stdout.

2. Giving use-cases of what these constructs could be used for.

The above construct can be used for concurrent processing.

3. What is the significance of the for loop with 4 iterations?

It is used to spawn 4 goroutines.

4. What is the significance of `make(chan func(), 10)`?

It creates a buffered channel (capacity 10) of functions that take in no
arguments and returns nothing.

5. Why is "HERE1" not getting printed?

The goroutines are not synchronized. As a result, the main goroutine
exits before the 4 goroutines have finished. Another problem with the
code is that the `cnp` channel is never closed. This means the for loops
never receives a signal to stop, leading to a deadlock.

A fixed version of this code would be:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	cnp := make(chan func(), 10)
	var wg sync.WaitGroup

	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			for f := range cnp {
				f()
			}
			wg.Done()
		}()
	}
	cnp <- func() {
		fmt.Println("HERE1")
	}
	fmt.Println("Hello")

	close(cnp)
	wg.Wait()
}
```
