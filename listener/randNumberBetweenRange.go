package listener

func (el *Listener) randNumberBetweenRange(min, max int) int {
	return el.newRandGeneratorHeader().Intn(max-min) + min
}
