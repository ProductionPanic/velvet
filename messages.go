package velvet

type Msg interface{}

type Cmd func() Msg

type KeyMsg struct {
	Key rune
}
