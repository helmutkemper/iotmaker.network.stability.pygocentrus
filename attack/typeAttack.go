package attack

import "github.com/helmutkemper/pygocentrus/changeContent"

type Attack struct {
	Enabled       bool
	Delay         RateMaxMin
	DontRespond   RateMaxMin
	ChangeContent changeContent.ChangeContent
	DeleteContent float64
}
