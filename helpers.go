package velvet

import (
	"github.com/ProductionPanic/velvet/style"
	"github.com/mattn/go-runewidth"
)
import "regexp"

func NewCellStyle() *style.CellStyle {
	return style.DefaultCellStyle()
}

func getTemplateStringWidth(template string) int {
	re := regexp.MustCompile(`\[[^\]]*\]`)
	cleaned := re.ReplaceAllString(template, "")
	return runewidth.StringWidth(cleaned)
}
