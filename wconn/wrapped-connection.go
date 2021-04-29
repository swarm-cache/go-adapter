package wconn

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/swarm-cache/go-adapter/glob"
	"github.com/swarm-cache/go-adapter/lib"
)

// Alias
type J = glob.J

// Callback types
type callback func(string, J, *[]byte)
type callbacks map[string]callback

// The base structure. Holds all required data to initialize the adapter
type Adapter struct {
	// The websocket connection
	Conn *websocket.Conn

	// The callbacks traversed
	callbacks map[string]callback

	// The callbacks mutex
	muxCB *sync.Mutex

	// The connection mutex
	muxConn *sync.Mutex
}

// Initializes the adapter
func Init() *Adapter {
	return &Adapter{
		callbacks: make(callbacks, 0),
		muxCB:     &sync.Mutex{},
		muxConn:   &sync.Mutex{},
	}
}

// Dial connection to a node addr.
func (a *Adapter) Connect(addr string) error {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+addr, nil)
	if err != nil {
		return fmt.Errorf("Could not connect to node! %s", err)
	}

	a.Conn = conn

	go a.handleMessage()

	return nil
}

// Messages handler.
//
// An infinite loop waiting for incoming messages upon which executes callbacks.
func (a *Adapter) handleMessage() {
	for {
		_, input, err := a.Conn.ReadMessage()
		if err != nil {
			log.Errorf("Could not read message! %s", err)
			break
		}

		// Decodes input message
		err, meta, data := lib.DecodeMessage(&input)
		if err != nil {
			log.Errorf("Could not decode incoming message! %s", err)
		}

		go a.TraverseCBs(meta, data)
	}
}
