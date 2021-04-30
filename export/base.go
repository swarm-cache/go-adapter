package export

import (
	"time"

	"github.com/swarm-cache/go-adapter/glob"
	"github.com/swarm-cache/go-adapter/lib"
	"github.com/swarm-cache/go-adapter/wconn"
)

// The general purpose json map type.
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
func Connect(addr string, cfg *J) (error, *operator) {

	// first register configs
	if cfg != nil {
		for n, v := range *cfg {
			if n == "V_NODE_RES_TIMEOUT" {
				glob.V_NODE_RES_TIMEOUT = time.Duration(v.(int))
			}

			if n == "V_MSG_ID_LENGTH" {
				glob.V_MSG_ID_LENGTH = v.(int)
			}

			if n == "V_LOG_IO_MSG" {
				glob.V_LOG_IO_MSG = v.(bool)
			}
		}
	}

	// then init
	wc := wconn.Init()

	if err := wc.Connect(addr); err != nil {
		return err, nil
	}

	return nil, &operator{
		wconn: wc,
	}
}
