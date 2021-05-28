package pygocentrus

type Pygocentrus struct {
	Enabled          bool
	Delay            RateMaxMin
	DontRespond      RateMaxMin
	ChangeLength     float64
	ChangeContent    ChangeContent
	DeleteContent    float64
	successfulAttack bool
}

func (el *Pygocentrus) SetAttack() {
	el.successfulAttack = true
}

func (el *Pygocentrus) GetAttack() bool {
	return el.successfulAttack
}

func (el *Pygocentrus) ClearAttack() {
	el.successfulAttack = false
}
