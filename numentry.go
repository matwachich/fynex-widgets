package wx

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type NumEntry struct {
	widget.Entry

	Float  bool
	Signed bool

	OnChanged func(string)

	lastValidInt int
	OnChangedInt func(int)

	lastValidFloat float64
	OnChangedFloat func(float64)

	ToolTipable

	// custom callbacks
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(r rune) (block bool)
	OnTypedKey      func(k *fyne.KeyEvent) (block bool)
	OnTypedShortcut func(s fyne.Shortcut) (block bool)

	readOnly bool
	minCols  int
}

func NewNumEntry() *NumEntry {
	n := &NumEntry{}
	n.ExtendBaseWidget(n)
	n.Entry.OnChanged = func(s string) {
		if n.OnChanged != nil {
			n.OnChanged(s)
		}
		if !n.Float && n.OnChangedInt != nil {
			i := n.GetInt()
			if i != n.lastValidInt {
				n.OnChangedInt(i)
				n.lastValidInt = i
			}
		}
		if n.Float && n.OnChangedFloat != nil {
			f := n.GetFloat()
			if f != n.lastValidFloat {
				n.OnChangedFloat(f)
				n.lastValidFloat = f
			}
		}
	}
	return n
}

func (n *NumEntry) ReadOnly() bool { return n.readOnly }
func (n *NumEntry) SetReadOnly(b bool) {
	n.readOnly = b
	n.Refresh()
}

func (e *NumEntry) SetMinColsVisible(c int) {
	e.minCols = c
	e.Refresh()
}

func (e *NumEntry) MinSize() fyne.Size {
	sz := e.Entry.MinSize()
	if e.minCols > 0 {
		sz.Width = fyne.MeasureText(strings.Repeat("0", e.minCols), theme.TextSize(), e.TextStyle).Width + 2*theme.InnerPadding() + 2*theme.InputBorderSize()
	}
	return sz
}

func (n *NumEntry) SetText(s string) {
	n.Entry.SetText("")
	if s == "" {
		return
	}

	old := n.Entry.OnChanged
	n.Entry.OnChanged = nil

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		n.TypedRune(runes[i])
	}

	n.Entry.OnChanged = old
	n.Entry.OnChanged(n.Entry.Text)
}

func (n *NumEntry) SetInt(i int) {
	if i == 0 {
		n.Entry.SetText("")
	}
	n.Entry.Text = strconv.Itoa(i)
	n.Entry.CursorColumn = len(n.Entry.Text)
	n.Entry.Refresh()
	n.Entry.OnChanged(n.Entry.Text)
}

func (n *NumEntry) SetFloat(f float64) {
	if f == 0 {
		n.Entry.SetText("")
	}
	n.Entry.Text = strings.Replace(fmt.Sprint(f), ".", ",", 1)
	n.Entry.CursorColumn = len(n.Entry.Text)
	n.Entry.Refresh()
	n.Entry.OnChanged(n.Entry.Text)
}

func (n *NumEntry) GetInt() int {
	return int(math.Abs(n.GetFloat()))
}

func (n *NumEntry) GetFloat() float64 {
	f, _ := strconv.ParseFloat(strings.Replace(n.Entry.Text, ",", ".", 1), 64)
	return f
}

func (n *NumEntry) TypedRune(r rune) {
	if n.readOnly {
		return
	}

	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if n.fnHasSign() {
			n.Entry.CursorColumn = 1
		}
		if n.OnTypedRune != nil && n.OnTypedRune(r) {
			return
		}
		n.Entry.TypedRune(r)
	// ---
	case '.', ',':
		if n.Float && !strings.Contains(n.Entry.Text, ",") {
			if n.fnHasSign() {
				n.Entry.CursorColumn = 1
			}
			if n.OnTypedRune != nil && n.OnTypedRune(r) {
				return
			}
			n.Entry.TypedRune(',')
		}
	// ---
	case '+', '-':
		if n.Signed {
			n.updateSign(r)
			if n.OnChanged != nil {
				n.OnChanged(n.Entry.Text)
			}
		}
	}
}

func (n *NumEntry) TypedKey(ke *fyne.KeyEvent) {
	if n.OnTypedKey != nil && n.OnTypedKey(ke) {
		return
	}
	n.Entry.TypedKey(ke)
}

func (n *NumEntry) FocusGained() {
	if n.readOnly {
		return
	}
	n.Entry.FocusGained()
	if n.OnFocusGained != nil {
		n.OnFocusGained()
	}
}

func (n *NumEntry) FocusLost() {
	n.Entry.FocusLost()
	if n.OnFocusLost != nil {
		n.OnFocusLost()
	}
}

func (n *NumEntry) TypedShortcut(s fyne.Shortcut) {
	if n.OnTypedShortcut != nil && n.OnTypedShortcut(s) {
		return
	}
	switch s := s.(type) {
	case *fyne.ShortcutPaste:
		if n.readOnly {
			return
		}
		for _, r := range s.Clipboard.Content() {
			n.TypedRune(r)
		}
	default:
		n.Entry.TypedShortcut(s)
	}
}

func (n *NumEntry) MouseIn(me *desktop.MouseEvent)    { n.ToolTipable.MouseIn(me) }
func (n *NumEntry) MouseMoved(me *desktop.MouseEvent) { n.ToolTipable.MouseMoved(me) }
func (n *NumEntry) MouseOut()                         { n.ToolTipable.MouseOut() }

func (n *NumEntry) updateSign(r rune) {
	if len(n.Text) > 0 {
		if n.Text[0] == '+' || n.Text[0] == '-' {
			n.Entry.SetText(string(r) + n.Text[1:])
		} else {
			n.Entry.SetText(string(r) + n.Text)
			n.Entry.CursorColumn += 1
		}
	} else {
		n.Entry.Text = string(r)
		n.Entry.CursorColumn = 1
		n.Entry.Refresh()
	}
}

func (n *NumEntry) fnHasSign() bool {
	return len(n.Entry.Text) > 0 && n.Entry.CursorColumn < 1 && (n.Entry.Text[0] == '+' || n.Entry.Text[0] == '-')
}
