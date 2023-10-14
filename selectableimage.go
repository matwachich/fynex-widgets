package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SelectableImage struct {
	widget.BaseWidget

	win fyne.Window

	Image *canvas.Image
	bk    *canvas.Circle
	check *canvas.Image

	Selected   bool
	OnSelected func(b bool)
}

func NewSelectableImage(win fyne.Window) *SelectableImage {
	w := &SelectableImage{win: win}
	w.ExtendBaseWidget(w)

	w.Image = &canvas.Image{}

	w.bk = &canvas.Circle{
		Hidden:    true,
		FillColor: theme.PrimaryColor(),
	}

	w.check = &canvas.Image{
		Resource: theme.ConfirmIcon(),
	}
	w.check.SetMinSize(fyne.NewSquareSize(theme.IconInlineSize() + 2*theme.Padding()))

	return w
}

func (w *SelectableImage) Refresh() {
	if w.Selected {
		w.check.Show()
		w.bk.Show()
	} else {
		w.check.Hide()
		w.bk.Hide()
	}
}

func (w *SelectableImage) Tapped(_ *fyne.PointEvent) {
	w.Selected = !w.Selected
	w.Refresh()

	if w.OnSelected != nil {
		w.OnSelected(w.Selected)
	}
}

func (w *SelectableImage) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (w *SelectableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(
		w.Image,
		container.NewCenter(container.NewStack(w.bk, w.check)),
	))
}
