package wx

import "time"

type DoubleClickable struct {
	lastTap time.Time
}

func (dbl *DoubleClickable) IsDoubleClick() (ret bool) {
	if time.Since(dbl.lastTap)/time.Millisecond <= time.Duration(doubleClickTime()) {
		dbl.lastTap = time.Time{}
		ret = true
	} else {
		dbl.lastTap = time.Now()
	}
	return
}
