package pages

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

const mountHeight = 4

func DiskPage(data view.Disk, x1, y1, x2, y2 int) tui.Drawable {
	flex := widgets.NewFlex()
	flex.Direction = widgets.FlexColumn
	flex.Border = false
	flex.SetRect(x1-1, y1-1, x2+1, y2+1)

	mountCount := (y2 - y1) / mountHeight
	for _, mount := range visibleMounts(data.Mounts, mountCount) {
		widget := widgets.NewParagraph()
		widget.Title = mount.Path
		widget.Text = fieldsText(mount.Info)
		flex.AddItem(widget, mountHeight, 0, false)
	}

	return flex
}

func visibleMounts(mounts []view.Mount, count int) []view.Mount {
	if count < 1 {
		return nil
	}

	if len(mounts) > count {
		return mounts[:count]
	}

	return mounts
}
