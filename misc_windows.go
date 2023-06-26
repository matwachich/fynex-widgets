package wx

import (
	"math/rand"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	getWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	getWindowTextW       = user32.NewProc("GetWindowTextW")
	getForegroundWindow  = user32.NewProc("GetForegroundWindow")
	enumWindows          = user32.NewProc("EnumWindows")
)

var hwnds map[fyne.Window]uintptr

func init() {
	hwnds = make(map[fyne.Window]uintptr)
}

// GetHWND retreives the HWND (WINAPI Window Handle) of a fyne.Window.
func GetHWND(w fyne.Window) (hwnd uintptr) {
	// check for stored HWND for this fyne.Window
	if h, ok := hwnds[w]; ok {
		return h
	}

	// generate a random 128 chars title
	rndTitle := make([]byte, 128)
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range rndTitle {
		rndTitle[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	// keep original window title
	oldTitle := w.Title()
	defer w.SetTitle(oldTitle)

	// set random title
	w.SetTitle(string(rndTitle))

	// fast check: check foreground window first
	if h, _, _ := getForegroundWindow.Call(); h != 0 {
		if sz, _, _ := getWindowTextLengthW.Call(h); sz == 128 {
			buf := make([]uint16, sz+1)
			getWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), sz+1)

			if string(rndTitle) == syscall.UTF16ToString(buf) {
				hwnds[w] = h
				return h
			}
		}
	}

	// lookup our window by the random title
	enumWindows.Call(syscall.NewCallback(func(h, lparam uintptr) uintptr {
		sz, _, _ := getWindowTextLengthW.Call(h)
		if sz != 128 { // fast check: must be the same length
			return 1
		}

		buf := make([]uint16, sz+1)
		getWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), sz+1)

		if string(rndTitle) == syscall.UTF16ToString(buf) {
			hwnd = h
			return 0
		}
		return 1
	}), 0)

	// store HWND for further lookups
	if hwnd != 0 {
		hwnds[w] = hwnd
	}
	return
}
