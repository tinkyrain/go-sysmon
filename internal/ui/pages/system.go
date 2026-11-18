package pages

import (
	"fmt"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func SystemPage(data view.System, x1, y1, x2, y2 int) tui.Drawable {
	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = fieldsText(data.Info)

	loadAvg := widgets.NewParagraph()
	loadAvg.Title = "LoadAvg"
	loadAvg.Text = fieldsText(data.LoadAvg)

	uptime := widgets.NewParagraph()
	uptime.Title = "Uptime"
	uptime.TitleAlignment = tui.AlignCenter
	uptime.TextAlignment = tui.AlignCenter
	uptime.VerticalAlignment = tui.AlignMiddle
	uptime.Text = fmt.Sprintf("Total: %s    Idle: %s", data.Uptime.Total, data.Uptime.Idle)

	infoFlex := widgets.NewFlex()
	infoFlex.Border = false
	infoFlex.Direction = widgets.FlexRow
	infoFlex.AddItem(info, 0, 1, false)
	infoFlex.AddItem(loadAvg, 0, 1, false)

	infoFlex.PaddingTop = -1
	infoFlex.PaddingBottom = -1
	infoFlex.PaddingLeft = -1
	infoFlex.PaddingRight = -1

	uptimeFlex := widgets.NewFlex()
	uptimeFlex.Border = false
	uptimeFlex.Direction = widgets.FlexColumn
	uptimeFlex.AddItem(uptime, 0, 1, false)

	uptimeFlex.PaddingTop = -1
	uptimeFlex.PaddingBottom = -1
	uptimeFlex.PaddingLeft = -1
	uptimeFlex.PaddingRight = -1

	flex := widgets.NewFlex()
	flex.Border = false
	flex.Direction = widgets.FlexColumn
	flex.AddItem(infoFlex, 0, 1, false)
	flex.AddItem(uptimeFlex, 5, 0, false)
	flex.SetRect(x1-1, y1-1, x2+1, y2+1)

	return flex
}
