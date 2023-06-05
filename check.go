package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type CheckEx struct {
	widget.Check

	ToolTipable

	// custom callbacks
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(rune) (block bool)
	OnTypedKey      func(*fyne.KeyEvent) (block bool)
	OnTypedShortcut func(fyne.Shortcut) (block bool)
}

func NewCheckEx(text string, changed func(bool)) *CheckEx {
	c := &CheckEx{}
	c.Check = widget.Check{
		DisableableWidget: widget.DisableableWidget{},
		Text:              text,
		OnChanged:         changed,
	}
	c.ToolTipable.parent = c
	c.ExtendBaseWidget(c)
	return c
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (c *CheckEx) FocusGained() {
	c.Check.FocusGained()
	if c.OnFocusGained != nil {
		c.OnFocusGained()
	}
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (c *CheckEx) FocusLost() {
	c.Check.FocusLost()
	if c.OnFocusLost != nil {
		c.OnFocusLost()
	}
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (c *CheckEx) TypedRune(r rune) {
	if c.OnTypedRune != nil && c.OnTypedRune(r) {
		return
	}
	c.Check.TypedRune(r)
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (c *CheckEx) TypedKey(k *fyne.KeyEvent) {
	/*if k.Name == fyne.KeyEnter {
		c.TypedKey(&fyne.KeyEvent{
			Name: fyne.KeySpace,
		})
	}*/
	if c.OnTypedKey != nil && c.OnTypedKey(k) {
		return
	}
	c.Check.TypedKey(k)
}

func (c *CheckEx) TypedShortcut(s fyne.Shortcut) {
	if c.OnTypedShortcut != nil && c.OnTypedShortcut(s) {
		return
	}
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (c *CheckEx) MouseIn(me *desktop.MouseEvent) {
	c.ToolTipable.MouseIn(me)
	c.Check.MouseIn(me)
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (c *CheckEx) MouseMoved(me *desktop.MouseEvent) {
	c.ToolTipable.MouseMoved(me)
	c.Check.MouseMoved(me)
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (c *CheckEx) MouseOut() {
	c.ToolTipable.MouseOut()
	c.Check.MouseOut()
}
