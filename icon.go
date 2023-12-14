package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Icon struct {
	widget.Icon

	OnTapped          func(e *fyne.PointEvent)
	OnTappedSecondary func(e *fyne.PointEvent)

	PointerCursor bool
}

func NewIcon(res fyne.Resource) *Icon {
	w := &Icon{}
	w.ExtendBaseWidget(w)
	w.Resource = res
	return w
}

func (w *Icon) SetColor(clr fyne.ThemeColorName) {
	w.SetResource(theme.NewColoredResource(w.Resource, clr))
}

func (w *Icon) Tapped(e *fyne.PointEvent) {
	if w.OnTapped != nil {
		w.OnTapped(e)
	}
}

func (w *Icon) TappedSecondary(e *fyne.PointEvent) {
	if w.OnTappedSecondary != nil {
		w.OnTappedSecondary(e)
	}
}

func (w *Icon) Cursor() desktop.Cursor {
	if w.PointerCursor && (w.OnTapped != nil || w.OnTappedSecondary != nil) {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}
