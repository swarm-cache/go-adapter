package export

import (
	"fmt"
	"sync"
	"time"

	"github.com/swarm-cache/go-adapter/glob"
	"github.com/swarm-cache/go-adapter/lib"
)

// Retrieves a key from nodes swarm
//
// Method waits until a response is received
func (o *operator) Get(key string) (error, *[]byte) {
	msgID, _ := lib.GenerateRandomString(glob.V_MSG_ID_LENGTH)
	mux := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	var rMeta J
	var rData *[]byte

	// Callback will wait for incoming messages from node
	// then will be invoked automatically by the traversal method.
	//
	// See wconn/wrapped-connections.go @ handleMessage
	o.wconn.PushCB(msgID, func(name string, meta J, data *[]byte) {
		if meta["msgID"].(string) == name {
			mux.Lock()
			rMeta = meta
			rData = data
			mux.Unlock()

			o.wconn.DelCB(name)
			wg.Done()
		}
	})

	// Send request
	o.wconn.Send(J{
		"msgID": msgID,
		"type":  "comm",
		"comm":  "get",
		"key":   key,
	}, nil)

	// handles timeout
	go func() {
		time.Sleep(glob.V_NODE_RES_TIMEOUT * time.Millisecond)

		mux.Lock()
		defer mux.Unlock()

		// it must be a timeout if rMeta is nil!
		if rMeta == nil {
			o.wconn.DelCB(msgID)
			wg.Done()
		}
	}()

	wg.Wait()

	// This occurs only when timeout is hit
	if rMeta == nil {
		return fmt.Errorf("Timeout!"), nil
	}

	// General switch executing action depending on meta response code
	switch rMeta["code"].(float64) {

	case glob.RES_SUCCESS:
		return nil, rData

	case glob.RES_NOT_FOUND:
		return fmt.Errorf("Item not found!"), nil

	default:
		return fmt.Errorf("Error occurred! %s", rMeta["message"]), nil

	}

}
