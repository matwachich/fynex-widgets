package wx

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ToolTipable struct {
	ToolTip fyne.CanvasObject

	parent  fyne.Widget
	timer   *time.Timer
	tooltip *widget.PopUp
	pos     fyne.Position
}

func (tt *ToolTipable) SetToolTip(title, text string, icon fyne.Resource) {
	if title == "" && text == "" {
		tt.ToolTip = nil
		return
	}
	if title == "" || text == "" {
		tt.ToolTip = widget.NewLabel(title + text)
	} else {
		if icon != nil {
			tt.ToolTip = container.NewVBox(
				container.NewHBox(
					&widget.Label{Text: title, TextStyle: fyne.TextStyle{Bold: true}},
					layout.NewSpacer(), widget.NewIcon(icon),
				),
				widget.NewLabel(text),
			)
		} else {
			tt.ToolTip = container.NewVBox(
				&widget.Label{Text: title, TextStyle: fyne.TextStyle{Bold: true}},
				widget.NewLabel(text),
			)
		}
	}
}

func (tt *ToolTipable) MouseIn(me *desktop.MouseEvent) {
	if p, ok := tt.parent.(fyne.Disableable); ok && p.Disabled() {
		return
	}
	if tt.ToolTip != nil && tt.timer == nil {
		tt.pos = me.AbsolutePosition
		tt.timer = time.AfterFunc(1500*time.Millisecond, func() {
			tt.tooltip = widget.NewPopUp(tt.ToolTip, fyne.CurrentApp().Driver().CanvasForObject(tt.parent))
			tt.tooltip.ShowAtPosition(tt.pos)
		})
	}
}

func (tt *ToolTipable) MouseMoved(me *desktop.MouseEvent) {
	tt.pos = me.AbsolutePosition
}

func (tt *ToolTipable) MouseOut() {
	if tt.timer != nil {
		tt.timer.Stop()
		tt.timer = nil
	}
}
