package velvet

type MouseEventType int

const (
	MousePress MouseEventType = iota
	MouseRelease
	MouseMotion
)

type MouseMsg struct {
	X, Y   int
	Button int
	Type   MouseEventType
	Alt    bool
	Ctrl   bool
	ZoneID string
}
