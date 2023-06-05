package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	widget.Button

	Menu *fyne.Menu

	OnRequestMenu func() *fyne.Menu

	ToolTipable
}

func NewButton(text string, icon fyne.Resource, importance widget.ButtonImportance, action func(), menu *fyne.Menu) *Button {
	b := &Button{
		Button: widget.Button{
			Text: text, Icon: icon, Importance: importance, OnTapped: action,
		},
		Menu: menu,
	}
	b.ToolTipable.parent = b
	b.ExtendBaseWidget(b)
	return b
}

func NewTBButton(text string, icon fyne.Resource, action func(), menu *fyne.Menu) *Button {
	b := &Button{
		Button: widget.Button{
			Text: text, Icon: icon, OnTapped: action,
			Importance: widget.LowImportance,
		},
		Menu: menu,
	}
	b.ToolTipable.parent = b
	b.ExtendBaseWidget(b)
	return b
}

func (b *Button) ToolbarObject() fyne.CanvasObject {
	return b
}

func (b *Button) Tap() {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	pos.Y += b.Size().Height + theme.Padding()
	b.Tapped(&fyne.PointEvent{
		Position:         fyne.NewPos(0, b.Size().Height+theme.Padding()),
		AbsolutePosition: pos,
	})
}

func (b *Button) TapSecondary() {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	pos.Y += b.Size().Height + theme.Padding()
	b.TappedSecondary(&fyne.PointEvent{
		Position:         fyne.NewPos(0, b.Size().Height+theme.Padding()),
		AbsolutePosition: pos,
	})
}

func (b *Button) Tapped(e *fyne.PointEvent) {
	if b.OnTapped == nil {
		b.TappedSecondary(e)
	} else {
		b.Button.Tapped(e)
	}
}

func (b *Button) TappedSecondary(e *fyne.PointEvent) {
	if b.Disabled() {
		return
	}
	if b.OnRequestMenu != nil {
		b.Menu = b.OnRequestMenu()
	}
	if b.Menu != nil && len(b.Menu.Items) > 0 {
		widget.ShowPopUpMenuAtPosition(b.Menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
	}
}

func (b *Button) MouseIn(e *desktop.MouseEvent) {
	b.ToolTipable.MouseIn(e)
	b.Button.MouseIn(e)
}

func (b *Button) MouseMoved(e *desktop.MouseEvent) {
	b.ToolTipable.MouseMoved(e)
	b.Button.MouseMoved(e)
}

func (b *Button) MouseOut() {
	b.ToolTipable.MouseOut()
	b.Button.MouseOut()
}
