### Redo - Re-run failing a function continuously
Redo is a small package that helps you rerun a function that continuously returns an 
error. It gives you the ability to determine the number of times you want the function
to be re-executed and the delay between each re-execution.

#### Installation
```
go get github.com/ncodes/redo
```

#### Full Example

```go
// Create a Redo instance, specifying the max retries and delay
redo := NewRedo(3, 100*time.Millisecond)

// Run a function that receives a method to stop the further re-execution.
// The function will be re-executed as long as it continues to return error.
err := redo.Do(func(stop func()) error {
    // stop() - Call to further stop re-execution
    return fmt.Errorf("something bad. redo")
})

// returns the error from the last execution
redo.LastErr    

// stop further re-execution
redo.Stop()
```

#### Documentation
[GoDoc](https://godoc.org/github.com/ncodes/redo)