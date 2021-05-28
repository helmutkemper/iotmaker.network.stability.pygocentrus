package iotmakernetworkstabilitypygocentrus

func (e *Proxy) intRange(min, max int) int {
	return e.rand().Intn(max-min) + min
}
