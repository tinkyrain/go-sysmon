package components

import (
	"time"

	"github.com/metaspartan/gotui/v5/widgets"
)

func Time(time time.Time, x1, y1, x2, y2 int) *widgets.Paragraph {
	paragraph := widgets.NewParagraph()
	paragraph.Border = false
	paragraph.Text = time.Format("02.01.2006 15:04:05")
	paragraph.SetRect(x1, y1, x2, y2)
	return paragraph
}
