package listener

import (
	"math/rand"
	"time"
)

func (el *Listener) newRandGeneratorHeader() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
