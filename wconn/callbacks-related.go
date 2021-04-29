package wconn

// Pushes a new callback onto map
//
// Name is used to uniquely identify the callback
func (a *Adapter) PushCB(name string, f callback) {
	a.muxCB.Lock()
	defer a.muxCB.Unlock()

	a.callbacks[name] = f
}

// Deletes a callback from map by name
func (a *Adapter) DelCB(name string) {
	a.muxCB.Lock()
	defer a.muxCB.Unlock()

	delete(a.callbacks, name)
}

// Traverses callbacks
func (a *Adapter) TraverseCBs(meta J, data *[]byte) {
	a.muxCB.Lock()
	defer a.muxCB.Unlock()

	for n, f := range a.callbacks {
		go f(n, meta, data)
	}
}
