package pages

import (
	"fmt"
	"strings"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
)

func DiskPage(stat disk.Stats, x1, y1, x2, y2 int) tui.Drawable {
	flex := widgets.NewFlex()
	flex.Direction = widgets.FlexColumn
	flex.Border = false
	flex.SetRect(x1, y1, x2, y2)

	for _, mount := range stat.Mounts {
		widget := widgets.NewParagraph()
		widget.Title = mount.Path
		widget.Text = strings.Join([]string{
			fmt.Sprintf("Total: %s", HumanizeSize(float64(mount.Total))),
			fmt.Sprintf("Available: %s", HumanizeSize(float64(mount.Available))),
		}, "\n")
		flex.AddItem(widget, 4, 0, false)
	}

	return flex
}
