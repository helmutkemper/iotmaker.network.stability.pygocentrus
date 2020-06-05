package pygocentrus

import "github.com/helmutkemper/pygocentrus/changeContent"

type Pygocentrus struct {
	Enabled       bool
	Delay         RateMaxMin
	DontRespond   RateMaxMin
	ChangeContent changeContent.ChangeContent
	DeleteContent float64
}
