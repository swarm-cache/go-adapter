package export

import (
	"github.com/swarm-cache/go-adapter/lib"
	"github.com/swarm-cache/go-adapter/wconn"
)

type J = lib.J

// The exported adapter for operating the wrapped connections
//
// Actual connection and data manipulation methods must be called via this structure
type operator struct {
	wconn *wconn.Adapter
}

// Returns a new operator instance.
//
// Automatically attempts to connect to the given node addr.
func Connect(addr string) (error, *operator) {
	wc := wconn.Init()

	if err := wc.Connect(addr); err != nil {
		return err, nil
	}

	return nil, &operator{
		wconn: wc,
	}
}
