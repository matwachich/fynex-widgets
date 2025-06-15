package wx

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TODO intégrer le fyne-x Calendar ; bofbof... il est bugué

type DateEntry struct {
	widget.Entry

	//valid bool

	lastValidTime time.Time

	ToolTipable

	// custom callbacks
	OnChanged       func(time.Time)
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(r rune) (block bool)
	OnTypedKey      func(k *fyne.KeyEvent) (block bool)
	OnTypedShortcut func(s fyne.Shortcut) (block bool)

	readOnly bool

	popup *widget.PopUp
	cal   *Calendar
	today *widget.Button
}

func NewDateEntry() *DateEntry {
	d := &DateEntry{}
	d.ToolTipable.parent = d
	d.ExtendBaseWidget(d)

	d.Text = "__/__/____"
	d.Entry.OnChanged = func(s string) {
		tm := d.GetTime()
		if !tm.Equal(d.lastValidTime) {
			if d.OnChanged != nil {
				d.OnChanged(tm)
			}
			d.lastValidTime = tm

			d.calendarSetWithoutCallback(tm)
		}
	}
	/*d.Entry.Validator = func(s string) (err error) { // la fonction n'est pas appelée! j'sais pas pourquoi
		if s == "__/__/____" {
			return
		}
		_, err = time.Parse("02/01/2006", s)
		return
	}*/

	d.Entry.ActionItem = &widget.Button{Icon: theme.CalendarIcon(), Importance: widget.LowImportance, OnTapped: func() {
		if d.popup == nil {
			c := fyne.CurrentApp().Driver().CanvasForObject(d)
			d.popup = widget.NewPopUp(container.NewVBox(d.cal, d.today), c)
		}

		pos := d.Position()
		pos.Y += d.MinSize().Height
		d.popup.ShowAtRelativePosition(pos, d)

		// set date after show, because cal internal widgets are create in CreateRenderer
		d.calendarSetWithoutCallback(d.GetTime())
	}}
	d.cal = NewCalendar(time.Now(), time.Time{}, func(t time.Time) {
		d.SetTime(t)
		d.popup.Hide()
	})
	d.cal.Selectable = true

	d.today = &widget.Button{Text: "Aujourd'hui", Alignment: widget.ButtonAlignCenter, Importance: widget.LowImportance, OnTapped: func() {
		d.SetTime(time.Now())
		d.popup.Hide()
	}}
	return d
}

func (d *DateEntry) calendarSetWithoutCallback(date time.Time) {
	oldCB := d.cal.OnChanged
	d.cal.OnChanged = nil
	d.cal.SetSelectedDate(date)
	d.cal.SetDisplayedDate(date)
	d.cal.OnChanged = oldCB
}

func (d *DateEntry) SetWeekStart(wd time.Weekday) {
	d.cal.SetWeekStart(wd)
}

func (d *DateEntry) callOnChanged() {
	d.Entry.OnChanged(d.Text)
}

func (d *DateEntry) SetText(s string) {
	d.Text = "__/__/____"
	d.CursorColumn = 0
	for _, r := range s {
		d.TypedRune(r)
	}
	d.Refresh()
	d.callOnChanged()
}

func (d *DateEntry) GetText() string {
	tm, err := time.ParseInLocation("02/01/2006", d.Text, time.Local)
	if err != nil || tm.IsZero() {
		return ""
	}
	return d.Text
}

func (d *DateEntry) SetTime(tm time.Time) {
	if tm.IsZero() {
		d.Text = "__/__/____"
		d.CursorColumn = 0
	} else {
		d.Text = tm.Format("02/01/2006")
		d.CursorColumn = 10
	}
	d.Refresh()
	d.callOnChanged()
}

func (d *DateEntry) GetTime() time.Time {
	tm, _ := time.ParseInLocation("02/01/2006", d.Text, time.Local)
	return tm
}

func (d *DateEntry) ReadOnly() bool {
	return d.readOnly
}
func (d *DateEntry) SetReadOnly(b bool) {
	d.readOnly = b
	if cnv := fyne.CurrentApp().Driver().CanvasForObject(d); b && cnv != nil && cnv.Focused() == d {
		cnv.Focus(nil)
	}
	if b && d.popup != nil {
		d.popup.Hide()
	}
	if b {
		d.Entry.ActionItem.(fyne.Disableable).Disable()
	} else {
		d.Entry.ActionItem.(fyne.Disableable).Enable()
	}
	d.Refresh()
}

func (d *DateEntry) MinSize() fyne.Size {
	s := d.Entry.MinSize()
	s.Width = fyne.MeasureText("00/00/0000", theme.TextSize(), d.TextStyle).Width + 2*theme.InnerPadding() + 2*theme.InputBorderSize() + s.Height // trick! pour ajouter la largeur du bouton d'action
	return s
}

func (d *DateEntry) FocusGained() {
	if d.readOnly {
		return
	}
	d.Entry.FocusGained()
	if d.OnFocusGained != nil {
		d.OnFocusGained()
	}
}

func (d *DateEntry) FocusLost() {
	d.Entry.FocusLost()
	if d.OnFocusLost != nil {
		d.OnFocusLost()
	}
}

func (d *DateEntry) TypedRune(r rune) {
	if d.OnTypedRune != nil && d.OnTypedRune(r) {
		return
	}

	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		switch d.CursorColumn {
		// __/__/____ 0, 1, /2, 3, 4, /5, 6, 7, 8, 9 [, 10]
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
			t := []rune(d.Text)
			if d.CursorColumn == 2 || d.CursorColumn == 5 {
				d.CursorColumn += 1
			}
			t[d.CursorColumn] = r
			d.CursorColumn += 1
			if d.CursorColumn == 2 || d.CursorColumn == 5 {
				d.CursorColumn += 1
			}
			d.Text = string(t)
			d.Refresh()
			d.callOnChanged()
		}
	}
}

func (d *DateEntry) TypedKey(k *fyne.KeyEvent) {
	if d.OnTypedKey != nil && d.OnTypedKey(k) {
		return
	}

	switch k.Name {
	case fyne.KeyRight:
		d.CursorColumn += 1
		if d.CursorColumn >= 10 {
			d.CursorColumn = 10
		}
		if d.CursorColumn == 2 || d.CursorColumn == 5 {
			d.CursorColumn += 1
		}
	case fyne.KeyLeft:
		d.CursorColumn -= 1
		if d.CursorColumn <= 0 {
			d.CursorColumn = 0
		}
		if d.CursorColumn == 2 || d.CursorColumn == 5 {
			d.CursorColumn -= 1
		}
	case fyne.KeyUp:
		switch d.CursorColumn {
		case 0, 1, 2:
			d.setDay(d.getDay()+1, true)
		case 3, 4, 5:
			d.setMonth(d.getMonth()+1, true)
		case 6, 7, 8, 9, 10:
			d.setYear(d.getYear() + 1)
		}
		d.callOnChanged()
	case fyne.KeyDown:
		switch d.CursorColumn {
		case 0, 1, 2:
			d.setDay(d.getDay()-1, true)
		case 3, 4, 5:
			d.setMonth(d.getMonth()-1, true)
		case 6, 7, 8, 9, 10:
			d.setYear(d.getYear() - 1)
		}
		d.callOnChanged()
	case fyne.KeyBackspace:
		// __/__/____ 0, 1, /2, 3, 4, /5, 6, 7, 8, 9 [, 10]
		t := []rune(d.Text)
		switch d.CursorColumn {
		case 10, 9, 8, 7, 5, 4, 2, 1:
			t[d.CursorColumn-1] = '_'
			d.CursorColumn -= 1
			if d.CursorColumn == 3 || d.CursorColumn == 6 {
				d.CursorColumn -= 1
			}
		case 3, 6:
			d.CursorColumn -= 1
		}
		d.Text = string(t)
		d.callOnChanged()
	case fyne.KeyDelete, fyne.KeyEscape:
		d.Text = "__/__/____"
		d.CursorColumn = 0
		d.callOnChanged()
	default:
		return
	}
	d.Refresh()
}

func (d *DateEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if d.OnTypedShortcut != nil && d.OnTypedShortcut(shortcut) {
		return
	}

	if s, ok := shortcut.(*fyne.ShortcutPaste); ok {
		for _, r := range s.Clipboard.Content() {
			d.TypedRune(r)
		}
		d.callOnChanged()
	} else {
		d.Entry.TypedShortcut(shortcut)
	}
}

func (d *DateEntry) MouseIn(me *desktop.MouseEvent)    { d.ToolTipable.MouseIn(me) }
func (d *DateEntry) MouseMoved(me *desktop.MouseEvent) { d.ToolTipable.MouseMoved(me) }
func (d *DateEntry) MouseOut()                         { d.ToolTipable.MouseOut() }

// ------------------------------------------------------------------------------------------------

func (d *DateEntry) setDay(day int, loop bool) {
	maxDay := 30
	switch d.getMonth() {
	case 4, 6, 9, 11:
		maxDay = 30
	case 1, 3, 5, 7, 8, 10, 12:
		maxDay = 31
	case 2:
		year := d.getYear()
		if year == 0 {
			maxDay = 28
		} else {
			if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
				maxDay = 29
			} else {
				maxDay = 28
			}
		}
	}
	if day > maxDay {
		if loop {
			day = 1
		} else {
			day = maxDay
		}
	}
	if day < 1 {
		if loop {
			day = maxDay
		} else {
			day = 1
		}
	}
	d.Text = fmt.Sprintf("%02d", day) + d.Text[2:]
}
func (d *DateEntry) getDay() int {
	ret, err := strconv.Atoi(d.Text[:2]) // __/__/____
	if err != nil {
		return 0
	}
	return ret
}

func (d *DateEntry) setMonth(month int, loop bool) {
	if month > 12 {
		if loop {
			month = 1
		} else {
			month = 12
		}
	}
	if month < 1 {
		if loop {
			month = 12
		} else {
			month = 1
		}
	}
	sMonth := fmt.Sprintf("%02d", month)
	d.Text = d.Text[:3] + sMonth + d.Text[5:]
}
func (d *DateEntry) getMonth() int {
	ret, err := strconv.Atoi(d.Text[3:5]) // __/__/____
	if err != nil {
		return 0
	}
	return ret
}

func (d *DateEntry) setYear(year int) {
	if year > 9999 {
		year = 9999
	}
	if year < 1 {
		year = 1
	}
	sYear := fmt.Sprintf("%04d", year)
	d.Text = d.Text[:6] + sYear
}
func (d *DateEntry) getYear() int {
	ret, err := strconv.Atoi(d.Text[6:]) // __/__/____
	if err != nil {
		return 0
	}
	return ret
}
