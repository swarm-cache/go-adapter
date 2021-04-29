package export

import (
	"fmt"
	"sync"
	"time"

	"github.com/swarm-cache/go-adapter/glob"
	"github.com/swarm-cache/go-adapter/lib"
)

// Deletes value by key
func (o *operator) Del(key string) error {
	msgID, _ := lib.GenerateRandomString(glob.V_MSG_ID_LENGTH)
	mux := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	var rMeta J

	// push callback
	o.wconn.PushCB(msgID, func(name string, meta J, fData *[]byte) {
		if meta["msgID"].(string) == name {
			mux.Lock()
			rMeta = meta
			mux.Unlock()

			o.wconn.DelCB(name)
			wg.Done()
		}
	})

	o.wconn.Send(J{
		"msgID": msgID,
		"type":  "comm",
		"comm":  "del",
		"key":   key,
	}, nil)

	// timeout handler
	go func() {
		time.Sleep(glob.V_NODE_RES_TIMEOUT * time.Millisecond)

		mux.Lock()
		defer mux.Unlock()

		if rMeta == nil {
			o.wconn.DelCB(msgID)
			wg.Done()
		}
	}()

	//
	wg.Wait()

	if rMeta == nil {
		return fmt.Errorf("Timeout!")
	}

	switch rMeta["code"].(float64) {

	case glob.RES_SUCCESS:
		return nil

	default:
		return fmt.Errorf("Error occurred! Data may have been deleted already. %s", rMeta["message"])

	}
}
