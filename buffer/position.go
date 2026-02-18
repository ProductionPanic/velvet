package buffer

type PositionX int
type PositionY int

const (
	AlignLeft   PositionX = -1000
	AlignCenter PositionX = -1001
	AlignRight  PositionX = -1002
)

const (
	AlignTop    PositionY = -1000
	AlignMiddle PositionY = -1001
	AlignBottom PositionY = -1002
)

func AbsoluteX(x int) PositionX {
	return PositionX(x)
}

func AbsoluteY(y int) PositionY {
	return PositionY(y)
}
