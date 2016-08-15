package helpers

// DrainChannel takes all messages out of channel
// No generics!! wtf can't write chan interface{}
func DrainChannel(c chan bool) {
	for {
		select {
		case <-c:
		default:
			return
		}
	}
}
