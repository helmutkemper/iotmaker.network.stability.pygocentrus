package pygocentrus

type Protocol int

func (el Protocol) String() string {
	return protocols[el]
}

var protocols = [...]string{
	"",
	"tcp",
}

const (
	KProtocolTCP Protocol = iota + 1
)
