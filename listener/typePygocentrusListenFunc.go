package listener

type pygocentrusListenFunc func(inData []byte) (int, []byte)
