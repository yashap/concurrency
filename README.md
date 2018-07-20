# concurrency
Just me working through the following talks/presentations about Go concurrency patterns:
* [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
* [Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
* [Share Memory by Communicating](https://golang.org/doc/codewalk/sharemem/)

## Dev Workflow

```bash
# Run the tests
go test -v ./...

# Run the main program - just a sandbox really
go run main.go
```

## Notes
Presentation order was:

**Go Concurrency Patterns**
```go
schedule.Every // simple example of goroutines, channels, and the generator pattern
multiplex.FanInNaive // shows an example of fan in behaviour, using 2 goroutines
multiplex.FanIn // shows an example of fan in behaviour, using select statements
schedule.ReceiveWithTimeout // another select example, receives a value from a channel, with a timeout
schedule.ReceiveMultiWithTimeout // as above, but receives zero-to-many values, until a timeout
schedule.EveryStoppable // like every, but lets you stop the generator
search.* // progressively more complex examples of "real-life" concurrency
```

**Advanced Go Concurrency Patterns**
```go
// TODO
```
**Share Memory by Communicating**
```go
// TODO
```
