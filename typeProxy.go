package iotmakernetworkstabilitypygocentrus

type Proxy struct {
	min        int
	max        int
	bufferSize int

	parser ParserInterface
}
