package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SelectEntryEx struct {
	widget.SelectEntry

	readOnly bool

	// custom callbacks
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(rune) (block bool)
	OnTypedKey      func(*fyne.KeyEvent) (block bool)
	OnTypedShortcut func(fyne.Shortcut) (block bool)
}

func NewSelectEntryEx() *SelectEntryEx {
	sel := &SelectEntryEx{}
	sel.Wrapping = fyne.TextTruncate
	sel.ExtendBaseWidget(sel)
	return sel
}

func (sel *SelectEntryEx) ReadOnly() bool { return sel.readOnly }
func (sel *SelectEntryEx) SetReadOnly(b bool) {
	sel.readOnly = b
	if b {
		sel.SelectEntry.Entry.ActionItem.(fyne.Disableable).Disable()
	} else {
		sel.SelectEntry.Entry.ActionItem.(fyne.Disableable).Enable()
	}
}

func (sel *SelectEntryEx) SetText(s string) {
	sel.SelectEntry.CursorColumn = 0
	sel.SelectEntry.CursorRow = 0
	sel.SelectEntry.SetText(s)
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (sel *SelectEntryEx) FocusGained() {
	if sel.readOnly {
		return
	}
	sel.SelectEntry.FocusGained()
	if sel.OnFocusGained != nil {
		sel.OnFocusGained()
	}
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (sel *SelectEntryEx) FocusLost() {
	sel.SelectEntry.FocusLost()
	if sel.OnFocusLost != nil {
		sel.OnFocusLost()
	}
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (sel *SelectEntryEx) TypedRune(r rune) {
	if sel.readOnly {
		return
	}
	if sel.OnTypedRune != nil && sel.OnTypedRune(r) {
		return
	}
	sel.SelectEntry.TypedRune(r)
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (sel *SelectEntryEx) TypedKey(e *fyne.KeyEvent) {
	if sel.readOnly {
		switch e.Name {
		case fyne.KeyEnter, fyne.KeyReturn, fyne.KeyBackspace, fyne.KeyDelete:
			return
		}
	}
	if sel.OnTypedKey != nil && sel.OnTypedKey(e) {
		return
	}
	sel.SelectEntry.TypedKey(e)
}

func (sel *SelectEntryEx) TypedShortcut(s fyne.Shortcut) {
	if sel.readOnly {
		switch s.(type) {
		case *fyne.ShortcutPaste:
			return
		case *fyne.ShortcutCut:
			s = &fyne.ShortcutCopy{Clipboard: s.(*fyne.ShortcutCut).Clipboard}
		}
	}
	if sel.OnTypedShortcut != nil && sel.OnTypedShortcut(s) {
		return
	}
	sel.SelectEntry.TypedShortcut(s)
}
