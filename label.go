package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Label struct {
	widget.Label

	minWidth float32
}

func NewLabel(text string) *Label {
	w := &Label{}
	w.ExtendBaseWidget(w)
	return w
}

func (w *Label) SetMinWidth(width float32) {
	w.minWidth = width
}

func (w *Label) MinSize() fyne.Size {
	sz := w.Label.MinSize()
	if sz.Width < w.minWidth {
		sz.Width = w.minWidth
	}
	return sz
}
