//go:build !windows

package wx

import "fyne.io/fyne/v2"

// GetHWND retreives the HWND (WINAPI Window Handle) of a fyne.Window.
func GetHWND(w fyne.Window) (hwnd uintptr) { return }

func MaximizeWindow(hwnd uintptr) {}

func doubleClickTime() uint {
	return 500 // ms
}
