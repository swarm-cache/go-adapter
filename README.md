### Abstract
Golang client adaptor for nodes. Used to connect to and communicate with nodes from:  https://github.com/swarm-cache/node. Can set, get, delete data from a swarm of nodes. Operations set/get/del are asynchronous by nature but handled in synchronous manner.

### Include
```
go get github.com/swarm-cache/go-adapter
```

### Use
```go
package main

import (
	"fmt"
	"strconv"
	"time"

	cache "github.com/swarm-cache/go-adapter/export"
)

func main() {
  // Connect to a node from swarm.
  // You can pass connection parameters, see example: https://github.com/swarm-cache/go-adapter/blob/main/tests/main.go
  err, c := cache.Connect("127.0.0.1:3666", nil)
  if err != nil {
    panic(err)
  }
  
  key := "hello"
  value := []byte("world!")
  
  // Setting a key => value pair into swarm
  err := c.Set(key, &value)
  if err != nil {
    panic(err)
  }
  
  // Retrieving a value from swarm by key
  c.Get(key, value)
  
  // Deleting a value from swarm by key
  c.Del(key)
}
```
