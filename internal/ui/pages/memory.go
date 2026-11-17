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

	memoryUsageBlock := components.UsageBlock(data.RAM.Label, data.RAM.Percent)
	swapUsageBlock := components.UsageBlock(data.Swap.Label, data.Swap.Percent)

	infoFlex := widgets.NewFlex()
	infoFlex.Direction = widgets.FlexColumn
	infoFlex.Border = false

	// Negative padding aligns flex content by removing empty lines that appear when border = false
	infoFlex.PaddingTop = -1
	infoFlex.PaddingBottom = -1
	infoFlex.PaddingLeft = -1
	infoFlex.PaddingRight = -1

	infoFlex.AddItem(info, 0, 1, false)

	memoryUsageFlex := widgets.NewFlex()
	memoryUsageFlex.Direction = widgets.FlexColumn
	memoryUsageFlex.Border = false

	// Negative padding aligns flex content by removing empty lines that appear when border = false
	memoryUsageFlex.PaddingTop = -1
	memoryUsageFlex.PaddingBottom = -1
	memoryUsageFlex.PaddingLeft = -1
	memoryUsageFlex.PaddingRight = -1

	memoryUsageFlex.AddItem(memoryUsageBlock, halfHeight, 0, false)
	memoryUsageFlex.AddItem(swapUsageBlock, (y2-y1)-halfHeight, 0, false)

	flex := widgets.NewFlex()
	flex.Direction = widgets.FlexRow
	flex.Border = false
	flex.AddItem(infoFlex, halfWidth, 0, false)
	flex.AddItem(memoryUsageFlex, (x2-x1)-halfWidth, 0, false)
	flex.SetRect(x1-1, y1-1, x2+1, y2+1)

	return flex
}
