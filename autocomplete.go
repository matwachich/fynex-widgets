package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AutoCompleteDataProvider interface {
	Length() int
	Create() fyne.CanvasObject
	Update(id int, co fyne.CanvasObject)
	Complete(id int) string
}

func (ac *AutoComplete) data_length() int {
	if ac.Data == nil {
		return len(ac.Options)
	} else {
		return ac.Data.Length()
	}
}

func (ac *AutoComplete) data_create() fyne.CanvasObject {
	if ac.Data == nil {
		return &widget.Label{}
	} else {
		return ac.Data.Create()
	}
}

func (ac *AutoComplete) data_update(id int, co fyne.CanvasObject) {
	if ac.Data == nil {
		co.(*widget.Label).Text = ac.Options[id]
	} else {
		ac.Data.Update(id, co)
	}
}

func (ac *AutoComplete) data_complete(id int) string {
	if ac.Data == nil {
		if ac.OnCompleted == nil {
			return ac.Options[id]
		} else {
			return ac.OnCompleted(ac.Options[id])
		}
	} else {
		return ac.Data.Complete(id)
	}
}

// ----------------------------------------------

type AutoComplete struct {
	widget.Entry // AutoComplete extends widget.Entry

	// autocomplete
	Options     []string
	OnCompleted func(string) string

	Data AutoCompleteDataProvider

	SubmitOnCompleted bool // if true, completing from list triggers OnSubmited

	/*CustomCreate func() fyne.CanvasObject
	CustomUpdate func(id widget.ListItemID, co fyne.CanvasObject)*/

	popup    *widget.PopUp
	list     *autoCompleteList
	selected widget.ListItemID
	pause    bool

	// custom callbacks
	OnFocusGained   func()
	OnFocusLost     func()
	OnTypedRune     func(r rune) (block bool)
	OnTypedKey      func(k *fyne.KeyEvent) (block bool)
	OnTypedShortcut func(s fyne.Shortcut) (block bool)

	// tooltips
	ToolTipable

	readonly bool
}

func NewAutoComplete(minLines int) *AutoComplete {
	ac := &AutoComplete{}
	ac.ExtendBaseWidget(ac)
	if minLines > 1 {
		ac.Entry.MultiLine = true
		ac.Entry.Wrapping = fyne.TextWrapWord
		ac.Entry.SetMinRowsVisible(minLines)
	} else {
		ac.Entry.Wrapping = fyne.TextTruncate
	}
	ac.ToolTipable.parent = ac
	return ac
}

func (ac *AutoComplete) ReadOnly() bool { return ac.readonly }
func (ac *AutoComplete) SetReadOnly(b bool) {
	ac.readonly = b
	if cnv := fyne.CurrentApp().Driver().CanvasForObject(ac); b && cnv != nil && cnv.Focused() == ac {
		cnv.Focus(nil)
	}
	ac.ListHide()
	ac.Refresh()
}

func (ac *AutoComplete) AcceptsTab() bool { return false }

func (ac *AutoComplete) FocusGained() {
	if ac.readonly {
		return
	}
	ac.Entry.FocusGained()
	if ac.OnFocusGained != nil {
		ac.OnFocusGained()
	}
}

func (ac *AutoComplete) FocusLost() {
	ac.Entry.FocusLost()
	ac.ListHide()
	if ac.OnFocusLost != nil {
		ac.OnFocusLost()
	}
}

func (ac *AutoComplete) TypedRune(r rune) {
	if ac.readonly {
		return
	}
	if ac.OnTypedRune != nil && ac.OnTypedRune(r) {
		return
	}
	ac.Entry.TypedRune(r)
}

func (ac *AutoComplete) TypedKey(k *fyne.KeyEvent) {
	if ac.readonly {
		return
	}
	if ac.OnTypedKey != nil && ac.OnTypedKey(k) {
		return
	}
	ac.Entry.TypedKey(k)
}

func (ac *AutoComplete) TypedShortcut(s fyne.Shortcut) {
	if ac.readonly {
		switch s.(type) {
		case *fyne.ShortcutPaste:
			return
		case *fyne.ShortcutCut:
			s = &fyne.ShortcutCopy{Clipboard: s.(*fyne.ShortcutCut).Clipboard}
		}
	}
	if ac.OnTypedShortcut != nil && ac.OnTypedShortcut(s) {
		return
	}
	ac.Entry.TypedShortcut(s)
}

func (ac *AutoComplete) MouseIn(me *desktop.MouseEvent)    { ac.ToolTipable.MouseIn(me) }
func (ac *AutoComplete) MouseMoved(me *desktop.MouseEvent) { ac.ToolTipable.MouseMoved(me) }
func (ac *AutoComplete) MouseOut()                         { ac.ToolTipable.MouseOut() }

// ---

func (ac *AutoComplete) ListShow() {
	if ac.pause || ac.readonly {
		return
	}

	if ac.data_length() <= 0 {
		ac.ListHide()
		return
	}

	cnv := fyne.CurrentApp().Driver().CanvasForObject(ac)
	if cnv == nil {
		return // not show
	}

	if ac.list == nil {
		ac.list = newAutoCompleteList(ac)
	}
	if ac.popup == nil {
		ac.popup = widget.NewPopUp(ac.list, cnv)
	}

	ac.popup.ShowAtPosition(ac.popupPos())
	ac.popup.Resize(ac.popupMaxSize())

	ac.list.Select(0)
	cnv.Focus(ac.list)
}

func (ac *AutoComplete) ListHide() {
	if ac.popup != nil {
		ac.list.UnselectAll()
		ac.popup.Hide()
	}
}

func (ac *AutoComplete) ListVisible() bool {
	return ac.popup != nil && ac.popup.Visible()
}

func (ac *AutoComplete) SetText(s string) {
	ac.pause = true
	/*ac.Entry.CursorColumn = 0
	ac.Entry.CursorRow = 0*/ // really needed ???
	ac.Entry.SetText(s)
	ac.pause = false
}

func (ac *AutoComplete) Move(pos fyne.Position) {
	ac.Entry.Move(pos)
	if ac.popup != nil && ac.popup.Visible() {
		ac.popup.Move(ac.popupPos())
		ac.popup.Resize(ac.popupMaxSize())
	}
}

func (ac *AutoComplete) setTextFromList(id widget.ListItemID) {
	ac.popup.Hide()

	ac.pause = true

	ac.Entry.Text = ac.data_complete(id)
	ac.Entry.CursorColumn = len(ac.Entry.Text)
	ac.Entry.Refresh()

	ac.pause = false

	if ac.SubmitOnCompleted && ac.OnSubmitted != nil {
		ac.OnSubmitted(ac.Entry.Text)
	}
}

func (ac *AutoComplete) popupPos() fyne.Position {
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(ac)
	return pos.Add(fyne.NewPos(0, ac.Size().Height+theme.Padding()))
}

func (ac *AutoComplete) popupMaxSize() fyne.Size {
	cnv := fyne.CurrentApp().Driver().CanvasForObject(ac)
	if cnv == nil {
		return fyne.Size{}
	}

	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(ac)

	// define size boundaries
	minWidth := ac.Size().Width
	maxWidth := cnv.Size().Width - pos.X - theme.Padding()
	maxHeight := cnv.Size().Height - pos.Y - ac.MinSize().Height - 2*theme.Padding()

	// iterating items until the end or we reach maxHeight
	var width, height float32
	for i := 0; i < len(ac.Options); i++ {
		item := ac.list.CreateItem()
		ac.list.UpdateItem(i, item)
		sz := item.MinSize()
		if sz.Width > width {
			width = sz.Width
		}
		height += sz.Height + theme.Padding()
		if height > maxHeight {
			height = maxHeight
			break
		}
	}
	height += theme.Padding() // popup padding

	width += 2 * theme.Padding() // let some padding on the triling end of the longest item
	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	return fyne.NewSize(width, height)
}

// ------------------------------------------------------------------------------------------------

type autoCompleteList struct {
	widget.List
	parent *AutoComplete
}

func newAutoCompleteList(parent *AutoComplete) *autoCompleteList {
	list := &autoCompleteList{parent: parent}
	list.ExtendBaseWidget(list)

	list.List.Length = func() int {
		return parent.data_length()
	}
	list.List.CreateItem = func() fyne.CanvasObject {
		return newAutoCompleteListItem(parent, parent.data_create())
	}
	list.List.UpdateItem = func(id widget.ListItemID, co fyne.CanvasObject) {
		co.(*autoCompleteListItem).id = id
		parent.data_update(id, co.(*autoCompleteListItem).co)
		parent.list.SetItemHeight(id, co.MinSize().Height)
		co.Refresh()
	}

	list.List.OnSelected = func(id widget.ListItemID) {
		parent.selected = id
	}
	list.List.OnUnselected = func(_ widget.ListItemID) {
		parent.selected = -1
	}

	return list
}

func (list *autoCompleteList) AcceptsTab() bool {
	return true
}

func (list *autoCompleteList) FocusGained() {}
func (list *autoCompleteList) FocusLost()   {}

func (list *autoCompleteList) TypedRune(r rune) {
	list.parent.TypedRune(r)
}
func (list *autoCompleteList) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyDown:
		if list.parent.selected < len(list.parent.Options)-1 {
			list.parent.list.Select(list.parent.selected + 1)
		} else {
			list.parent.list.Select(0)
		}
	case fyne.KeyUp:
		if list.parent.selected > 0 {
			list.parent.list.Select(list.parent.selected - 1)
		} else {
			list.parent.list.Select(len(list.parent.Options) - 1)
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if list.parent.selected >= 0 {
			list.parent.setTextFromList(list.parent.selected)
		} else {
			list.parent.ListHide()
			list.parent.Entry.TypedKey(k)
		}
	case fyne.KeyTab, fyne.KeyEscape:
		list.parent.ListHide()
	default:
		list.parent.TypedKey(k)
	}
}

func (list *autoCompleteList) TypedShortcut(s fyne.Shortcut) { list.parent.TypedShortcut(s) }

// ---

type autoCompleteListItem struct {
	widget.BaseWidget
	parent *AutoComplete
	co     fyne.CanvasObject
	id     widget.ListItemID
}

func newAutoCompleteListItem(parent *AutoComplete, co fyne.CanvasObject) *autoCompleteListItem {
	item := &autoCompleteListItem{parent: parent, id: -1, co: co}
	item.ExtendBaseWidget(item)
	return item
}

func (item *autoCompleteListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.co)
}

func (item *autoCompleteListItem) Tapped(_ *fyne.PointEvent) {
	item.parent.setTextFromList(item.id)
}

func (item *autoCompleteListItem) MouseIn(_ *desktop.MouseEvent)    { item.parent.list.Select(item.id) }
func (item *autoCompleteListItem) MouseMoved(_ *desktop.MouseEvent) {}
func (item *autoCompleteListItem) MouseOut()                        {}
