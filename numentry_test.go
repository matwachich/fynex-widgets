package wx

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestNumEntryType(t *testing.T) {
	n := NewNumEntry()
	test.Type(n, "some text with number in it 123\nother line without numbers...")

	if n.Text != "123" {
		t.Fatal("n.Text != 123")
	}
	if n.GetInt() != 123 {
		t.Fatal("n.GetInt() != 123")
	}
	if n.GetFloat() != 123 {
		t.Fatal("n.GetFloat() != 123")
	}
}

func TestNumEntrySetText(t *testing.T) {
	n := NewNumEntry()
	n.SetText("some text with number in it 123\nother line without numbers...")

	if n.Text != "123" {
		t.Fatal("n.Text != 123")
	}
	if n.GetInt() != 123 {
		t.Fatal("n.GetInt() != 123")
	}
	if n.GetFloat() != 123 {
		t.Fatal("n.GetFloat() != 123")
	}
}

func TestNumEntrySetInt(t *testing.T) {
	n := NewNumEntry()
	n.SetInt(123)

	if n.Text != "123" {
		t.Fatal("n.Text != 123")
	}
	if n.GetInt() != 123 {
		t.Fatal("n.GetInt() != 123")
	}
	if n.GetFloat() != 123 {
		t.Fatal("n.GetFloat() != 123")
	}
}

func TestNumEntryOnChanged(t *testing.T) {
	n := NewNumEntry()
	n.Float = true
	n.Signed = true

	chkStr := ""
	chkInt := 0
	chkFlt := 0.0

	n.OnChanged = func(s string) {
		if s != chkStr {
			t.Error("n.OnChanged != " + chkStr)
		}
	}
	n.OnChangedInt = func(i int) {
		if i != chkInt {
			t.Errorf("n.OnChangedInt != %v", i)
		}
	}
	n.OnChangedFloat = func(f float64) {
		if f != chkFlt {
			t.Errorf("n.OnChangedFloat != %v", f)
		}
	}

	chkStr = "1"
	chkInt = 1
	chkFlt = 1.0
	test.Type(n, "1")

	chkStr = "12"
	chkInt = 12
	chkFlt = 12.0
	test.Type(n, "2")

	chkStr = "123"
	chkInt = 123
	chkFlt = 123.0
	test.Type(n, "3")

	chkStr = "123,"
	chkInt = 123
	chkFlt = 123.0
	test.Type(n, ".")

	chkStr = "123,5"
	chkInt = 123
	chkFlt = 123.5
	test.Type(n, "5")

	chkStr = "-123,5"
	chkInt = -123
	chkFlt = -123.5
	test.Type(n, "-")

	chkStr = "+123,5"
	chkInt = 123
	chkFlt = 123.5
	test.Type(n, "+")

	chkStr = "10"
	chkInt = 10
	chkFlt = 10
	n.SetText("10")

	chkStr = "-20"
	chkInt = -20
	chkFlt = -20.0
	n.SetText("-20")

	chkStr = "15"
	chkInt = 15
	chkFlt = 15
	n.SetInt(15)

	chkStr = "-30,5"
	chkInt = -30
	chkFlt = -30.5
	n.SetFloat(-30.5)
}
