package changeContent

import (
	"math/rand"
	"time"
)

func (el *ChangeContent) GetRandomByLength(length int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(length)
}
