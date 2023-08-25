# Go Uitls for linolab

## Async

ExecOnce

```go
import "github.com/linolab/go-utils/async"

func main() {
  hw := async.ExecOnceWrap(func() {
    time.Sleep(1 * time.Second)
    fmt.Println("Hello World")
  })

  go hw()
  go hw()
}

// Output:
// Hello World
```
