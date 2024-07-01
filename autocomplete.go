package wx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// AutoComplete is an extend widget.Entry that allows to display a list of suggestions
// beneath it.
//
// The suggestions can be (by default) a simple list of strings (.Options), or anything
// you want (complexe data structures).
//
// The CustomXxx callbacks are used for custom suggestion data.
type AutoComplete struct {
	widget.Entry // AutoComplete extends widget.Entry

	//
	Options []string // List of suggestions.

	//
	CustomLength   func() int                            // Returns the length of custom data source (widget.List like)
	CustomCreate   func() fyne.CanvasObject              // Creates a fyne.CanvasObject to display a custom data source item (widget.List like)
	CustomUpdate   func(id int, co fyne.CanvasObject)    // Updates a fyne.CanvasObject to display a custom data source item (widget.List like)
	CustomComplete func(id int) (ret string, close bool) // Called when a custom data source is used, to match a (complexe) item with its textual representation (that will be filled in the Entry)

	//
	SubmitOnCompleted bool // if true, completing from list (either with Enter key or with click) triggers Entry.OnSubmited if set

	// custom callbacks
	OnFocusGained   func()                              // Called when widget gains focus
	OnFocusLost     func()                              // Called when widget loses focus
	OnTypedRune     func(r rune) (block bool)           // Called when a rune is typed ; block = true will prevent the event from reaching the widget
	OnTypedKey      func(k *fyne.KeyEvent) (block bool) // Called when a key is typed ; block = true will prevent the event from reaching the widget (this will also prevent TypedRune to be called)
	OnTypedShortcut func(s fyne.Shortcut) (block bool)  // Called when a shortut is typed ; block = true will prevent the event from reaching the widget

	// tooltips
	ToolTipable

	// internals
	popup    *widget.PopUp
	list     *autoCompleteList
	selected widget.ListItemID
	pause    bool
	readonly bool
}

// NewAutoComplete creates an AutoComplete widget.
// minLines > 1 will create a multiline WordWrapped widget by default.
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

// ReadOnly returns read-only status.
//
// Read-Only widget will display like a normal non focused widget,
// but it will be impossible to modify its content (either by keyboard or from clipboard).
func (ac *AutoComplete) ReadOnly() bool { return ac.readonly }

// SetReadOnly sets read-only status.
//
// Read-Only widget will display like a normal non focused widget,
// but it will be impossible to modify its content (either by keyboard or from clipboard).
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
		switch ss := s.(type) {
		case *fyne.ShortcutPaste:
			return
		case *fyne.ShortcutCut:
			s = &fyne.ShortcutCopy{Clipboard: ss.Clipboard}
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

func (ac *AutoComplete) data_length() int {
	if ac.CustomLength == nil {
		return len(ac.Options)
	} else {
		return ac.CustomLength()
	}
}

func (ac *AutoComplete) data_create() fyne.CanvasObject {
	if ac.CustomCreate == nil {
		return &widget.Label{}
	} else {
		return ac.CustomCreate()
	}
}

func (ac *AutoComplete) data_update(id int, co fyne.CanvasObject) {
	if ac.CustomUpdate == nil {
		co.(*widget.Label).SetText(ac.Options[id])
		ac.list.SetItemHeight(id, co.MinSize().Height)
	} else {
		ac.CustomUpdate(id, co)
	}
}

func (ac *AutoComplete) data_complete(id int) (ret string, close bool) {
	if ac.CustomComplete == nil {
		return ac.Options[id], true
	} else {
		return ac.CustomComplete(id)
	}
}

// ----------------------------------------------

// ListShow will display the auto-completion list (if there is data to display).
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

// ListHide hides auto-complete list.
func (ac *AutoComplete) ListHide() {
	if ac.popup != nil {
		ac.list.UnselectAll()
		ac.popup.Hide()
	}
}

// ListVisible returns wether the auto-complete list is visible.
func (ac *AutoComplete) ListVisible() bool {
	return ac.popup != nil && ac.popup.Visible()
}

// SetText sets the text in the Entry without triggering OnChanged.
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

func (ac *AutoComplete) RefreshItem(id int) {
	ac.list.RefreshItem(id)
}

func (ac *AutoComplete) SetItemHeight(id int, height float32) {
	ac.list.SetItemHeight(id, height)
}

func (ac *AutoComplete) setTextFromList(id widget.ListItemID) {
	ac.pause = true

	var close bool
	ac.Entry.Text, close = ac.data_complete(id)
	ac.Entry.CursorColumn = len(ac.Entry.Text)
	ac.Entry.Refresh()

	ac.pause = false

	if ac.SubmitOnCompleted && ac.OnSubmitted != nil {
		ac.OnSubmitted(ac.Entry.Text)
	}

	if close {
		ac.popup.Hide()
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
	maxHeight := cnv.Size().Height - pos.Y - ac.MinSize().Height - 2*theme.Padding() - theme.Padding()

	// iterating items until the end or we reach maxHeight
	var width, height float32
	for i := 0; i < ac.data_length(); i++ {
		item := ac.data_create()
		ac.data_update(i, item)
		sz := item.MinSize()
		if sz.Width > width {
			width = sz.Width
		}
		height += sz.Height + theme.Padding() // FIXME when height is different from minHeight (wrapped content)
		if height > maxHeight {
			height = maxHeight
			break
		}
	}
	height += theme.Padding() // popup padding

	width += 2 * theme.Padding() // let some padding on the trailing end of the longest item
	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	return fyne.NewSize(width, height)
}

// ----------------------------------------------

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
		/*parent.list.SetItemHeight(id, co.MinSize().Height)
		co.Refresh()*/
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
		if list.parent.selected < list.parent.data_length()-1 {
			list.parent.list.Select(list.parent.selected + 1)
		} else {
			list.parent.list.Select(0)
		}
	case fyne.KeyUp:
		if list.parent.selected > 0 {
			list.parent.list.Select(list.parent.selected - 1)
		} else {
			list.parent.list.Select(list.parent.data_length() - 1)
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
