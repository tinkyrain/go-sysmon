package pages

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func MemoryPage(data view.Memory, x1, y1, x2, y2 int) tui.Drawable {
	halfWidth := (x2 - x1) / 2
	halfHeight := (y2 - y1) / 2

	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = fieldsText(data.Info)

	memoryUsageLine := components.UsageLine(data.RAM.Label, data.RAM.Percent)
	swapUsageLine := components.UsageLine(data.Swap.Label, data.Swap.Percent)

	infoFlex := widgets.NewFlex()
	infoFlex.Direction = widgets.FlexColumn
	infoFlex.Border = false
	infoFlex.AddItem(info, 0, 1, false)

	memoryUsageFlex := widgets.NewFlex()
	memoryUsageFlex.Direction = widgets.FlexColumn
	memoryUsageFlex.Border = false
	memoryUsageFlex.AddItem(memoryUsageLine, halfHeight, 0, false)
	memoryUsageFlex.AddItem(swapUsageLine, halfHeight, 0, false)

	flex := widgets.NewFlex()
	flex.Direction = widgets.FlexRow
	flex.Border = false
	flex.AddItem(infoFlex, halfWidth, 0, false)
	flex.AddItem(memoryUsageFlex, halfWidth, 0, false)
	flex.SetRect(x1-1, y1-1, x2+1, y2+1)

	return flex
}
