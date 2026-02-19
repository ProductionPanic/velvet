package velvet

import "github.com/ProductionPanic/velvet/buffer"

// Drawable represents any element that can be rendered to a buffer.Grid
type Drawable interface {
	Render() *buffer.Grid
}
