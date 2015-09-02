# abutil [![GoDoc](https://godoc.org/github.com/bahlo/abutil?status.svg)](https://godoc.org/github.com/bahlo/abutil) [![Build Status](https://travis-ci.org/bahlo/abutil.svg?branch=master)](https://travis-ci.org/bahlo/abutil) [![Coverage Status](https://coveralls.io/repos/bahlo/abutil/badge.svg?branch=master&service=github)](https://coveralls.io/github/bahlo/abutil?branch=master)

abutil is a collection of often-used Golang helpers.

## Contents
- [Functions](#functions)
  - [OnSignal](#onsignal)
  - [Parallel](#parallel)
  - [RollbackErr](#rollbackerr)
  - [RemoteIP](#remoteip)
  - [GracefulServer](#gracefulserver)
- [License](#license)

## Functions

#### [OnSignal](https://godoc.org/github.com/bahlo/abutil#OnSignal)
Listens to various signals and executes the given function with the received
signal.

```go
go abutil.OnSignal(func(s os.Signal) {
  fmt.Printf("Got signal %s\n", s)
})
```

#### [Parallel](https://godoc.org/github.com/bahlo/abutil#Parallel)
Executes the given function n times concurrently.

```go
var m sync.Mutex
c := 0
abutil.Parallel(4, func() {
    m.Lock()
    defer m.Unlock()

    fmt.Print(c)
    c++
})
```

#### [RollbackErr](https://godoc.org/github.com/bahlo/abutil#RollbackErr)
Does a rollback on the given transaction and returns either the rollback-error,
if occured, or the given one.

```go
insertSomething := func(db *sql.DB) error {
    tx, _ := db.Begin()

    _, err := tx.Exec("INSERT INTO some_table (some_column) VALUES (?)",
        "foobar")
    if err != nil {
        // The old way, imagine doing this 10 times in a method
        if err := tx.Rollback(); err != nil {
            return err
        }

        return err
    }

    _, err = tx.Exec("DROP DATABASE foobar")
    if err != nil {
        // With RollbackErr
        return abutil.RollbackErr(tx, err)
    }

    return nil
}
```

#### [RemoteIP](https://godoc.org/github.com/bahlo/abutil#RemoteIP)
Tries everything to get the remote ip.

```go
someHandler := func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("New request from %s\n", abutil.RemoteIP(r))
}
```

#### [GracefulServer](https://godoc.org/github.com/bahlo/abutil#GracefulServer)
A wrapper around `graceful.Server` from <http://github.com/tylerb/graceful>
with state variable and easier handling.

```go
s := abutil.NewGracefulServer(1337, someHandlerFunc)

// This channel blocks until all connections are closed or the time is up
sc := s.StopChan()

// Some go func that stops the server after 2 seconds for no reason
time.After(2 * time.Second, func() {
    s.Stop(10 * time.Second)
})

if err := s.ListenAndServe(); err != nil && !s.Stopped() {
    // We didn't stop the server, so something must be wrong
    panic(err)
}

// Wait for the server to finish
<-sc
```

## License

This project is licensed under the WTFPL, for more information see the LICENSE
file.
