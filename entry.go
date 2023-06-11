package wx

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	EntryExCaseUpper = strings.ToUpper
	EntryExCaseLower = strings.ToLower
	EntryExCaseTitle = func(s string) string {
		return cases.Title(language.French).String(strings.ToLower(s))
	}
)

type EntryEx struct {
	widget.Entry

	// widget configuration
	AcceptTab    bool
	RuneModifier func(rune) rune
	CaseModifier func(string) string // only for non multiline entries ; use RuneModifier for multiline

	ToolTipable

	// custom callbacks
	OnChanged       func(string)
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(r rune) (block bool)
	OnTypedKey      func(k *fyne.KeyEvent) (block bool)
	OnTypedShortcut func(s fyne.Shortcut) (block bool)

	readOnly bool
}

func NewEntryEx(minRows int) *EntryEx {
	e := &EntryEx{}
	e.ExtendBaseWidget(e)
	e.ToolTipable.parent = e
	e.Wrapping = fyne.TextTruncate
	if minRows > 1 {
		e.MultiLine = true
		e.Wrapping = fyne.TextWrapWord
		e.SetMinRowsVisible(minRows)
	} else {
		e.Wrapping = fyne.TextTruncate
	}
	return e
}

func (e *EntryEx) ExtendBaseWidget(wid fyne.Widget) {
	e.Entry.OnChanged = e.onChanged
	e.Entry.ExtendBaseWidget(wid)
}

func (e *EntryEx) onChanged(s string) {
	if !e.MultiLine && e.CaseModifier != nil {
		s = e.CaseModifier(s)
		e.Text = s
		e.Refresh()
	}
	if e.OnChanged != nil {
		e.OnChanged(s)
	}
}

func (e *EntryEx) AcceptsTab() bool {
	return e.AcceptTab
}

func (e *EntryEx) ReadOnly() bool { return e.readOnly }
func (e *EntryEx) SetReadOnly(b bool) {
	e.readOnly = b
	if cnv := fyne.CurrentApp().Driver().CanvasForObject(e); b && cnv != nil && cnv.Focused() == e {
		cnv.Focus(nil)
	}
	e.Entry.Refresh()
}

func (e *EntryEx) SetText(s string) {
	e.Entry.CursorColumn = 0
	e.Entry.CursorRow = 0
	e.Entry.SetText(s)
}

func (e *EntryEx) FocusGained() {
	if e.readOnly {
		return
	}
	e.Entry.FocusGained()
	if e.OnFocusGained != nil {
		e.OnFocusGained()
	}
}

func (e *EntryEx) FocusLost() {
	e.Entry.FocusLost()
	if e.OnFocusLost != nil {
		e.OnFocusLost()
	}
}

func (e *EntryEx) TypedRune(r rune) {
	if e.readOnly {
		return
	}
	if e.RuneModifier != nil {
		r = e.RuneModifier(r)
	}
	if e.OnTypedRune != nil && e.OnTypedRune(r) {
		return
	}
	e.Entry.TypedRune(r)
}

func (e *EntryEx) TypedKey(k *fyne.KeyEvent) {
	if e.readOnly {
		switch k.Name {
		case fyne.KeyEnter, fyne.KeyReturn, fyne.KeyBackspace, fyne.KeyDelete:
			return
		case fyne.KeyTab:
			if e.AcceptTab {
				return
			}
		}
	}
	if e.OnTypedKey != nil && e.OnTypedKey(k) {
		return
	}
	e.Entry.TypedKey(k)
}

func (e *EntryEx) TypedShortcut(s fyne.Shortcut) {
	if e.readOnly {
		switch s.(type) {
		case *fyne.ShortcutPaste:
			return
		case *fyne.ShortcutCut:
			s = &fyne.ShortcutCopy{Clipboard: s.(*fyne.ShortcutCut).Clipboard}
		}
	}
	if e.OnTypedShortcut != nil && e.OnTypedShortcut(s) {
		return
	}
	e.Entry.TypedShortcut(s)
}

func (e *EntryEx) MouseIn(me *desktop.MouseEvent)    { e.ToolTipable.MouseIn(me) }
func (e *EntryEx) MouseMoved(me *desktop.MouseEvent) { e.ToolTipable.MouseMoved(me) }
func (e *EntryEx) MouseOut()                         { e.ToolTipable.MouseOut() }
