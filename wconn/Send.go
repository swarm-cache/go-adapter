package wconn

import (
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"github.com/swarm-cache/go-adapter/lib"
)

// Sends a message to a node
func (a *Adapter) Send(meta J, data *[]byte) error {
	a.muxConn.Lock()
	defer a.muxConn.Unlock()

	err, out := lib.EncodeMessage(meta, data)
	if err != nil {
		log.Errorf("wrapped-con@Send - Error occurred: %s", err)
		return err
	}

	a.Conn.WriteMessage(websocket.BinaryMessage, *out)

	// if glob.F_LOG_IO_MSG {
	// 	log.Infof("Message sent! \nMeta: %s \nData: %s", meta, data)
	// }

	return nil
}
