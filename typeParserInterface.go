package iotmakernetworkstabilitypygocentrus

type ParserInterface interface {
	Parser(data []byte, direction string) (dataSize int, err error)
}
