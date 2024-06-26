package wx

import (
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FieldID = string

type inputField struct {
	fyne.Widget
	check *widget.Check
}

type InputFields struct {
	widget.BaseWidget

	OnChanged func(id FieldID)
	OnAction  func(id FieldID)

	inputs map[FieldID]*inputField
	order  []FieldID

	vbox *fyne.Container
	tabs *container.AppTabs
}

// ----------------------------------------------------------------------------
// Creation

func NewInputFields() (w *InputFields) {
	w = &InputFields{inputs: make(map[FieldID]*inputField)}
	w.ExtendBaseWidget(w)
	return w
}

func (w *InputFields) CreateRenderer() fyne.WidgetRenderer {
	if w.vbox != nil {
		return widget.NewSimpleRenderer(w.vbox)
	} else if w.tabs != nil {
		return widget.NewSimpleRenderer(w.tabs)
	} else {
		return widget.NewSimpleRenderer(widget.NewLabel("(!) ERREUR: Aucun champs d'entrée défini."))
	}
}

// ----------------------------------------------------------------------------
// Structure

func (w *InputFields) AddTab(title string, icon fyne.Resource) {
	if title == "" || w.vbox != nil {
		return
	}
	if w.tabs == nil {
		w.tabs = container.NewAppTabs()
	}
	w.tabs.Append(container.NewTabItemWithIcon(title, icon, container.NewVBox()))
}

func (w *InputFields) AddSeparator() {
	w.currentVBox().Add(widget.NewSeparator())
}

func (w *InputFields) AddTitle(id FieldID, text string, style fyne.TextStyle, alignement fyne.TextAlign) {
	w.currentVBox().Add(&widget.Label{Text: text, TextStyle: style, Alignment: alignement})
}

func (w *InputFields) AddTitleMkd(id FieldID, text string) {
	w.currentVBox().Add(widget.NewRichTextFromMarkdown(text))
}

func (w *InputFields) AddLabel(id FieldID, label string, text string, style fyne.TextStyle, alignement fyne.TextAlign) {
	w.addWidget(id, false, label, &widget.Label{
		Text:      text,
		Alignment: alignement,
		TextStyle: style,
	})
}

// ----------------------------------------------------------------------------
// Inputs

func (w *InputFields) AddText(id FieldID, nullable bool, label string, value string, lines int) {
	w.dummyId(&id)
	wid := NewEntryEx(lines)
	wid.Text = value
	wid.OnChanged = func(_ string) { w.onChanged(id) }
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddPassword(id FieldID, nullable bool, label string, value string) {
	w.dummyId(&id)
	wid := widget.NewEntry()
	wid.Text = value
	wid.Password = true
	wid.OnChanged = func(_ string) { w.onChanged(id) }
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddNumber(id FieldID, nullable bool, label string, value string, float, signed bool) {
	w.dummyId(&id)
	wid := NewNumEntry()
	wid.Float = float
	wid.Signed = signed
	wid.SetText(value)
	if !float {
		wid.OnChangedInt = func(_ int) { w.onChanged(id) }
	} else {
		wid.OnChangedFloat = func(_ float64) { w.onChanged(id) }
	}
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddDate(id FieldID, nullable bool, label string, value string) {
	w.dummyId(&id)
	wid := NewDateEntry()
	wid.SetText(value)
	wid.OnChanged = func(_ time.Time) { w.onChanged(id) }
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddSelect(id FieldID, nullable bool, label string, options []string, value string, editable bool) {
	w.dummyId(&id)
	var wid fyne.Widget
	if editable {
		wid = widget.NewSelectEntry(options)
		wid.(*widget.SelectEntry).Text = value
		wid.(*widget.SelectEntry).OnChanged = func(_ string) { w.onChanged(id) }
	} else {
		wid = widget.NewSelect(options, nil)
		wid.(*widget.Select).Selected = value
		wid.(*widget.Select).OnChanged = func(_ string) { w.onChanged(id) }
	}
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddCheck(id FieldID, nullable bool, label string, text string, value bool) {
	w.dummyId(&id)
	wid := widget.NewCheck(text, func(_ bool) { w.onChanged(id) })
	wid.Checked = value
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddCheckGroup(id FieldID, nullable bool, label string, options []string, values []string, horizontal bool) {
	w.dummyId(&id)
	wid := widget.NewCheckGroup(options, func(_ []string) { w.onChanged(id) })
	wid.Selected = values
	wid.Horizontal = horizontal
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddRadioGroup(id FieldID, nullable bool, label string, options []string, value string, horizontal bool) {
	w.dummyId(&id)
	wid := widget.NewRadioGroup(options, func(_ string) { w.onChanged(id) })
	wid.Selected = value
	wid.Horizontal = horizontal
	w.addWidget(id, nullable, label, wid)
}

func (w *InputFields) AddActionButton(id FieldID, label, btnText string, importance widget.Importance) {
	w.dummyId(&id)
	w.addWidget(id, false, label, &widget.Button{
		Text:       btnText,
		Importance: importance,
		OnTapped:   func() { w.onAction(id) },
	})
}

// ----------------------------------------------------------------------------
// Lire/Ecrire les inputs

func (w *InputFields) ReadString(id FieldID) (ret string) {
	f := w.inputs[id]
	if f == nil {
		return
	}

	if f.check != nil && !f.check.Checked {
		return ""
	}

	switch wid := f.Widget.(type) {
	case *widget.Label:
		ret = wid.Text
	case *EntryEx:
		ret = wid.Text
	case *widget.Entry: // password
		ret = wid.Text
	case *DateEntry:
		ret = wid.GetText()
	case *NumEntry:
		ret = wid.Text
		/*if ret == "" {
			ret = "0"
			if wid.Float {
				ret += ".0"
			}
		}*/
	case *widget.Select:
		ret = wid.Selected
	case *widget.SelectEntry:
		ret = wid.Text
	case *widget.Check:
		ret = strconv.FormatBool(wid.Checked)
	case *widget.CheckGroup:
		ret = strings.Join(wid.Selected, "|")
	case *widget.RadioGroup:
		ret = wid.Selected
	case *widget.Button:
		ret = wid.Text
	}
	return
}

func (w *InputFields) WriteString(id FieldID, value string) {
	f := w.inputs[id]
	if f == nil {
		return
	}

	old := w.OnChanged
	w.OnChanged = nil
	defer func() { w.OnChanged = old }()

	switch wid := f.Widget.(type) {
	case *widget.Label:
		wid.SetText(value)
	case *EntryEx:
		wid.SetText(value)
	case *widget.Entry: // password
		wid.SetText(value)
	case *DateEntry:
		wid.SetText(value)
	case *NumEntry:
		wid.SetText(value)
	case *widget.Select:
		wid.SetSelected(value)
	case *widget.SelectEntry:
		wid.SetText(value)
	case *widget.Check:
		if b, err := strconv.ParseBool(value); err == nil {
			wid.SetChecked(b)
		}
	case *widget.CheckGroup:
		wid.SetSelected(strings.Split(value, "|"))
	case *widget.RadioGroup:
		wid.SetSelected(value)
	case *widget.Button:
		wid.SetText(value)
	}
}

func (w *InputFields) WriteOptions(id FieldID, options []string) {
	f := w.inputs[id]
	if f == nil {
		return
	}

	switch wid := f.Widget.(type) {
	case *widget.Select:
		wid.SetOptions(options)
	case *widget.SelectEntry:
		wid.SetOptions(options)
	case *widget.CheckGroup:
		wid.Options = options
		wid.Refresh()
	case *widget.RadioGroup:
		wid.Options = options
		wid.Refresh()
	}
}

func (w *InputFields) ReadAllString() (ret map[FieldID]string) {
	ret = make(map[FieldID]string)
	for _, id := range w.order {
		ret[id] = w.ReadString(id)
	}
	return
}

func (w *InputFields) WriteAllString(data map[FieldID]string) {
	for k, v := range data {
		w.WriteString(k, v)
	}
}

// ----------------------------------------------------------------------------
// Manipuler les inputs

func (w *InputFields) Inputs() []FieldID {
	return w.order
}

func (w *InputFields) Widget(id FieldID) fyne.Widget {
	if f, ok := w.inputs[id]; ok {
		return f.Widget
	}
	return nil
}

func (w *InputFields) SetNull(id FieldID, b bool) {
	if f, ok := w.inputs[id]; ok {
		if f.check != nil {
			f.check.SetChecked(!b)
		}
	}
}

func (w *InputFields) GetNull(id FieldID) (b bool) {
	if f, ok := w.inputs[id]; ok {
		if f.check != nil {
			b = !f.check.Checked
		}
	}
	return
}

func (w *InputFields) SetStatus(id FieldID, b bool) {
	if f, ok := w.inputs[id]; ok {
		if dis, ok := f.Widget.(fyne.Disableable); ok {
			if b {
				if f.check == nil {
					dis.Enable()
				} else {
					f.check.Enable()
					f.check.SetChecked(f.check.Checked) // will update widget
				}
			} else {
				dis.Disable()
				if f.check != nil {
					f.check.Disable()
				}
			}
		}
	}
}

func (w *InputFields) GetStatus(id FieldID) (b bool) {
	if f, ok := w.inputs[id]; ok {
		if dis, ok := f.Widget.(fyne.Disableable); ok {
			b = !dis.Disabled()
		}
	}
	return
}

func (w *InputFields) SetFocus(id FieldID) {
	if f, ok := w.inputs[id]; ok {
		if foc, ok := f.Widget.(fyne.Focusable); ok {
			fyne.CurrentApp().Driver().CanvasForObject(f.Widget).Focus(foc)
		}
	}
}

// ----------------------------------------------------------------------------
// Disableable

func (w *InputFields) Enable() {
	for _, f := range w.inputs {
		if dis, ok := f.Widget.(fyne.Disableable); ok {
			dis.Enable()
			if f.check != nil {
				f.check.Enable()
			}
		}
	}
}

func (w *InputFields) Disable() {
	for _, f := range w.inputs {
		if dis, ok := f.Widget.(fyne.Disableable); ok {
			if f.check == nil {
				dis.Disable()
			} else {
				f.check.Disable()
				f.check.SetChecked(f.check.Checked) // will update widget
			}
		}
	}
}

func (w *InputFields) Disabled() (ret bool) {
	for _, f := range w.inputs {
		if dis, ok := f.Widget.(fyne.Disableable); ok && dis.Disabled() {
			ret = true
			break
		}
	}
	return
}

// TODO wx.Readonlyable ? (pour afficher un document avec DocType.Inputs en readonly (autre user))

// ----------------------------------------------------------------------------
// internals

func (w *InputFields) dummyId(id *FieldID) {
	if *id == "" {
		*id = "_noname_input_" + strconv.Itoa(len(w.order)+1)
	}
}

func (w *InputFields) addWidget(id FieldID, nullable bool, label string, wid fyne.Widget) {
	if _, ok := w.inputs[id]; ok {
		return
	}

	w.dummyId(&id)

	var cnt fyne.CanvasObject

	f := &inputField{Widget: wid}
	if dis, ok := wid.(fyne.Disableable); ok && nullable {
		f.check = widget.NewCheck("", func(b bool) {
			if b {
				dis.Enable()
			} else {
				dis.Disable()
			}
			w.onChanged(id)
		})
		f.check.Checked = true

		cnt = container.NewBorder(nil, nil, f.check, nil, wid)
	} else {
		cnt = wid
	}

	switch wid.(type) {
	case *widget.Button:
		if label == "" {
			w.currentVBox().Add(cnt)
		} else {
			form := w.currentForm()
			form.Objects = append(form.Objects,
				&widget.Label{Text: label, TextStyle: fyne.TextStyle{Bold: true}},
				cnt,
			)
		}
	default:
		form := w.currentForm()
		form.Objects = append(form.Objects,
			&widget.Label{Text: label, TextStyle: fyne.TextStyle{Bold: true}},
			cnt,
		)
	}

	w.inputs[id] = f
	w.order = append(w.order, id)
}

func (w *InputFields) currentVBox() *fyne.Container {
	if w.tabs == nil {
		if w.vbox == nil {
			w.vbox = container.NewVBox()
		}
		return w.vbox
	} else {
		return w.tabs.Items[len(w.tabs.Items)-1].Content.(*fyne.Container)
	}
}

func (w *InputFields) currentForm() *fyne.Container {
	vbox := w.currentVBox()
	if len(vbox.Objects) <= 0 {
		form := container.New(layout.NewFormLayout())
		vbox.Add(form)
		return form
	}
	if form, ok := vbox.Objects[len(vbox.Objects)-1].(*fyne.Container); ok {
		return form
	}
	form := container.New(layout.NewFormLayout())
	vbox.Add(form)
	return form
}

func (w *InputFields) onChanged(id FieldID) {
	if w.OnChanged != nil {
		w.OnChanged(id)
	}
}

func (w *InputFields) onAction(id FieldID) {
	if w.OnAction != nil {
		w.OnAction(id)
	}
}
