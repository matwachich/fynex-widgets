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

// InputBox is a complexe data input widget, intended to be used in popups (but not only...).
type InputBox struct {
	widget.BaseWidget

	// configuration
	Header           string // markdown supported
	Validate, Cancel string

	OnReady   func()              // called just after the widget is showed for the first time
	OnChanged func(name string)   // called every time an entry is modified
	OnAction  func(name string)   // called when an entry button is pressed
	OnSubmit  func(validate bool) // called when validate or cancel button is pressed

	//
	vbox *fyne.Container
	tabs *container.AppTabs

	btnOk, btnCancel *widget.Button

	inputs        map[string]fyne.Widget
	inputsOrder   []string
	onReadyCalled bool
}

func NewInputBox() (ib *InputBox) {
	ib = &InputBox{Validate: "Valider", Cancel: "Annuler", inputs: make(map[string]fyne.Widget)}
	ib.ExtendBaseWidget(ib)
	return ib
}

func (ib *InputBox) CreateRenderer() fyne.WidgetRenderer {
	// Header
	var header fyne.CanvasObject
	if ib.Header != "" {
		header = container.NewVBox(
			widget.NewRichTextFromMarkdown(ib.Header),
			widget.NewSeparator(),
		)
	}

	// Content
	var content fyne.CanvasObject
	if ib.tabs == nil {
		content = ib.vbox
	} else {
		content = ib.tabs
	}
	if content == nil {
		content = widget.NewLabel("Aucun champ d'entrée défini")
	}

	// Buttons
	if ib.Cancel == "" {
		ib.Cancel = "Annuler"
	}
	ib.btnCancel = &widget.Button{Text: ib.Cancel, OnTapped: func() { ib.callOnSubmit(false) }}

	var btnBar *fyne.Container
	if ib.Validate != "" {
		ib.btnOk = &widget.Button{Text: ib.Validate, Importance: widget.HighImportance, OnTapped: func() { ib.callOnSubmit(true) }}
		btnBar = container.NewGridWithColumns(4, layout.NewSpacer(), ib.btnCancel, ib.btnOk, layout.NewSpacer())
	} else {
		btnBar = container.NewGridWithColumns(5, layout.NewSpacer(), layout.NewSpacer(), ib.btnCancel, layout.NewSpacer(), layout.NewSpacer())
	}

	// Build layout
	return widget.NewSimpleRenderer(container.NewBorder(
		header,
		container.NewVBox(
			widget.NewSeparator(),
			btnBar,
		),
		nil, nil,
		content,
	))
}

func (ib *InputBox) Show() {
	ib.BaseWidget.Show()

	// focus the first focusable input
	if cnv := fyne.CurrentApp().Driver().CanvasForObject(ib); cnv != nil && len(ib.inputsOrder) > 0 {
		for _, name := range ib.inputsOrder {
			if input, ok := ib.inputs[name]; ok && input.Visible() {
				if w, ok := input.(fyne.Focusable); ok {
					cnv.Focus(w)
					break
				}
			}
		}
	}

	// call on ready
	if !ib.onReadyCalled && ib.OnReady != nil {
		ib.OnReady()
		ib.onReadyCalled = true
	}
}

// ------------------------------------------------------------------------------------------------

func (ib *InputBox) AddTab(label string) {
	if label == "" || ib.vbox != nil {
		return
	}
	if ib.tabs == nil {
		ib.tabs = container.NewAppTabs()
	}
	ib.tabs.Append(container.NewTabItem(label, container.NewVBox()))
}

func (ib *InputBox) AddSeparator() {
	ib.currentVBox().Add(widget.NewSeparator())
}

func (ib *InputBox) AddTitle(text string) {
	if text != "" {
		ib.currentVBox().Add(widget.NewRichTextFromMarkdown(text))
	}
}

// ---

func (ib *InputBox) AddLabel(name, label string, value string, bold, italic bool, alignement fyne.TextAlign) {
	ib.addEntryWidget(name, label, &widget.Label{
		Text:      value,
		Alignment: alignement,
		TextStyle: fyne.TextStyle{
			Italic: italic,
			Bold:   bold,
		},
	})
}

func (ib *InputBox) AddText(name, label, value string, lines int) {
	w := NewEntryEx(lines)
	w.Text = value
	w.OnChanged = func(_ string) { ib.callOnChanged(name) }
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddPassword(name, label string) {
	w := widget.NewEntry()
	w.Password = true
	w.OnChanged = func(_ string) { ib.callOnChanged(name) }
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddDate(name, label string, value string) {
	w := NewDateEntry()
	w.SetText(value)
	w.OnChanged = func(_ time.Time) { ib.callOnChanged(name) }
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddNumber(name, label string, value string, float, signed bool) {
	w := NewNumEntry()
	w.SetText(value)
	w.OnChanged = func(_ string) { ib.callOnChanged(name) }
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddSelect(name, label string, options []string, value string, editable bool) {
	var w fyne.Widget
	if !editable {
		w = widget.NewSelect(options, func(_ string) { ib.callOnChanged(name) })
		w.(*widget.Select).SetSelected(value)
	} else {
		w = widget.NewSelectEntry(options)
		w.(*widget.SelectEntry).SetText(value)
		w.(*widget.SelectEntry).OnChanged = func(s string) { ib.callOnChanged(name) }
	}
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddCheck(name, label, option string, value bool) {
	w := widget.NewCheck(option, func(_ bool) { ib.callOnChanged(name) })
	w.Checked = value
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddCheckGroup(name, label string, options []string, values []string, horizontal bool) {
	w := widget.NewCheckGroup(options, func(_ []string) { ib.callOnChanged(name) })
	w.SetSelected(values)
	w.Horizontal = horizontal
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddRadioGroup(name, label string, options []string, value string, horizontal bool) {
	w := widget.NewRadioGroup(options, func(_ string) { ib.callOnChanged(name) })
	w.SetSelected(value)
	w.Horizontal = horizontal
	ib.addEntryWidget(name, label, w)
}

func (ib *InputBox) AddButton(name, label, text string, importance widget.ButtonImportance) {
	ib.addEntryWidget(name, label, &widget.Button{
		Text:       text,
		Importance: importance,
		OnTapped:   func() { ib.callOnAction(name) },
	})
}

// ------------------------------------------------------------------------------------------------

func (ib *InputBox) Inputs() []string {
	return ib.inputsOrder
}

func (ib *InputBox) Widget(name string) fyne.Widget {
	if w, ok := ib.inputs[name]; ok {
		return w
	}
	return nil
}

func (ib *InputBox) ReadString(name string) (ret string) {
	w := ib.Widget(name)
	if w == nil {
		return
	}
	switch w := w.(type) {
	case *widget.Label:
		ret = w.Text
	case *EntryEx:
		ret = w.Text
	case *widget.Entry: // password
		ret = w.Text
	case *DateEntry:
		ret = w.GetText()
	case *NumEntry:
		ret = w.Text
	case *widget.Select:
		ret = w.Selected
	case *widget.SelectEntry:
		ret = w.Text
	case *widget.Check:
		ret = strconv.FormatBool(w.Checked)
	case *widget.CheckGroup:
		ret = strings.Join(w.Selected, "|")
	case *widget.RadioGroup:
		ret = w.Selected
	case *widget.Button:
		ret = w.Text
	}
	return
}

func (ib *InputBox) WriteString(name string, value string) {
	w := ib.Widget(name)
	if w == nil {
		return
	}
	switch w := w.(type) {
	case *widget.Label:
		w.SetText(value)
	case *EntryEx:
		w.SetText(value)
	case *widget.Entry: // password
		w.SetText(value)
	case *DateEntry:
		w.SetText(value)
	case *NumEntry:
		w.SetText(value)
	case *widget.Select:
		w.SetSelected(value)
	case *widget.SelectEntry:
		w.SetText(value)
	case *widget.Check:
		w.SetChecked(value != "")
	case *widget.CheckGroup:
		w.SetSelected(strings.Split(value, "|"))
	case *widget.RadioGroup:
		w.SetSelected(value)
	case *widget.Button:
		w.SetText(value)
	}
}

func (ib *InputBox) SetOptions(name string, value string) {
	switch w := ib.Widget(name).(type) {
	case *widget.Select:
		w.Options = strings.Split(value, "|")
		w.Refresh()
	case *widget.SelectEntry:
		w.SetOptions(strings.Split(value, "|"))
	case *widget.Check:
		w.Text = value
		w.Refresh()
	case *widget.CheckGroup:
		w.Options = strings.Split(value, "|")
		w.Refresh()
	case *widget.RadioGroup:
		w.Options = strings.Split(value, "|")
		w.Refresh()
	}
}

func (ib *InputBox) Focus(name string) {
	if cnv := fyne.CurrentApp().Driver().CanvasForObject(ib); cnv != nil {
		if input, ok := ib.inputs[name]; ok && input.Visible() {
			if focusable, ok := input.(fyne.Focusable); ok {
				cnv.Focus(focusable)
			}
		}
	}
}

func (ib *InputBox) SetStatus(name string, enable bool) {
	if input, ok := ib.inputs[name]; ok && input.Visible() {
		if w, ok := input.(fyne.Disableable); ok {
			if enable {
				w.Enable()
			} else {
				w.Disable()
			}
		}
	}
}

func (ib *InputBox) GetStatus(name string) bool {
	if input, ok := ib.inputs[name]; ok && input.Visible() {
		if w, ok := input.(fyne.Disableable); ok {
			return !w.Disabled()
		}
	}
	return false
}

func (ib *InputBox) SetSubmitable(b bool) {
	if ib.btnOk != nil {
		if b {
			ib.btnOk.Enable()
		} else {
			ib.btnOk.Disable()
		}
	}
}

func (ib *InputBox) GetSubmitable() bool {
	if ib.btnOk != nil {
		return !ib.btnOk.Disabled()
	}
	return false
}

// ------------------------------------------------------------------------------------------------

func (ib *InputBox) currentVBox() *fyne.Container {
	if ib.tabs == nil {
		if ib.vbox == nil {
			ib.vbox = container.NewVBox()
		}
		return ib.vbox
	} else {
		return ib.tabs.Items[len(ib.tabs.Items)-1].Content.(*fyne.Container)
	}
}

func (ib *InputBox) currentForm() *fyne.Container {
	vbox := ib.currentVBox()
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

func (ib *InputBox) addEntryWidget(name, label string, w fyne.Widget) {
	if name == "" {
		return
	}
	if _, ok := ib.inputs[name]; !ok {
		switch w.(type) {
		case *widget.Button:
			if label == "" {
				vbox := ib.currentVBox()
				vbox.Add(w)
			} else {
				form := ib.currentForm()
				form.Add(&widget.Label{Text: label, TextStyle: fyne.TextStyle{Bold: true}})
				form.Add(w)
			}
		default:
			form := ib.currentForm()
			form.Add(&widget.Label{Text: label, TextStyle: fyne.TextStyle{Bold: true}})
			form.Add(w)
		}

		ib.inputs[name] = w
		ib.inputsOrder = append(ib.inputsOrder, name)
	}
}

// ---

func (ib *InputBox) callOnChanged(name string) {
	if ib.OnChanged != nil {
		ib.OnChanged(name)
	}
}

func (ib *InputBox) callOnAction(name string) {
	if ib.OnAction != nil {
		ib.OnAction(name)
	}
}

func (ib *InputBox) callOnSubmit(validate bool) {
	if ib.OnSubmit != nil {
		ib.OnSubmit(validate)
	}
}
