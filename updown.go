package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type UpDownButton struct {
	widget.DisableableWidget
}

// geometry

// MinSize returns the minimum size this object needs to be drawn.
func (w *UpDownButton) MinSize() fyne.Size {
	panic("not implemented") // TODO: Implement
}

// Move moves this object to the given position relative to its parent.
// This should only be called if your object is not in a container with a layout manager.
func (w *UpDownButton) Move(_ fyne.Position) {
	panic("not implemented") // TODO: Implement
}

// Position returns the current position of the object relative to its parent.
func (w *UpDownButton) Position() fyne.Position {
	panic("not implemented") // TODO: Implement
}

// Resize resizes this object to the given size.
// This should only be called if your object is not in a container with a layout manager.
func (w *UpDownButton) Resize(_ fyne.Size) {
	panic("not implemented") // TODO: Implement
}

// Size returns the current size of this object.
func (w *UpDownButton) Size() fyne.Size {
	panic("not implemented") // TODO: Implement
}

// visibility

// Hide hides this object.
func (w *UpDownButton) Hide() {
	panic("not implemented") // TODO: Implement
}

// Visible returns whether this object is visible or not.
func (w *UpDownButton) Visible() bool {
	panic("not implemented") // TODO: Implement
}

// Show shows this object.
func (w *UpDownButton) Show() {
	panic("not implemented") // TODO: Implement
}

// Refresh must be called if this object should be redrawn because its inner state changed.
func (w *UpDownButton) Refresh() {
	panic("not implemented") // TODO: Implement
}

// CreateRenderer returns a new WidgetRenderer for this widget.
// This should not be called by regular code, it is used internally to render a widget.
func (w *UpDownButton) CreateRenderer() fyne.WidgetRenderer {
	panic("not implemented") // TODO: Implement
}

type upDownRenderer struct {
	widget *UpDownButton
}

// Destroy is for internal use.
func (r *upDownRenderer) Destroy() {
	panic("not implemented") // TODO: Implement
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
func (r *upDownRenderer) Layout(_ fyne.Size) {
	panic("not implemented") // TODO: Implement
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (r *upDownRenderer) MinSize() fyne.Size {
	panic("not implemented") // TODO: Implement
}

// Objects returns all objects that should be drawn.
func (r *upDownRenderer) Objects() []fyne.CanvasObject {
	panic("not implemented") // TODO: Implement
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (r *upDownRenderer) Refresh() {
	panic("not implemented") // TODO: Implement
}
