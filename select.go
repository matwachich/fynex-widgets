package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Select struct {
	widget.Select

	ToolTipable

	// custom callbacks
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(rune) (block bool)
	OnTypedKey      func(*fyne.KeyEvent) (block bool)
	OnTypedShortcut func(fyne.Shortcut) (block bool)
}

func NewSelect(options []string, onChanged func(s string)) *Select {
	sel := &Select{}
	sel.ToolTipable.parent = sel
	sel.ExtendBaseWidget(sel)
	sel.Options = options
	sel.OnChanged = onChanged
	return sel
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (sel *Select) FocusGained() {
	sel.Select.FocusGained()
	if sel.OnFocusGained != nil {
		sel.OnFocusGained()
	}
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (sel *Select) FocusLost() {
	sel.Select.FocusLost()
	if sel.OnFocusLost != nil {
		sel.OnFocusLost()
	}
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (sel *Select) TypedRune(r rune) {
	if sel.OnTypedRune != nil && sel.OnTypedRune(r) {
		return
	}
	sel.Select.TypedRune(r)
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (sel *Select) TypedKey(e *fyne.KeyEvent) {
	if sel.OnTypedKey != nil && sel.OnTypedKey(e) {
		return
	}
	sel.Select.TypedKey(e)
}

func (sel *Select) TypedShortcut(s fyne.Shortcut) {
	if sel.OnTypedShortcut != nil && sel.OnTypedShortcut(s) {
		return
	}
	// sel.Select.TypedShortcut(s) // doesn't exists
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (sel *Select) MouseIn(me *desktop.MouseEvent) {
	sel.ToolTipable.MouseIn(me)
	sel.Select.MouseIn(me)
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (sel *Select) MouseMoved(me *desktop.MouseEvent) {
	sel.ToolTipable.MouseMoved(me)
	sel.Select.MouseMoved(me)
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (sel *Select) MouseOut() {
	sel.ToolTipable.MouseOut()
	sel.Select.MouseOut()
}
