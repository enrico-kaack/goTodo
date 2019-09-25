package main

import "time"

func FormatDone(done bool) string {
	if done {
		return "☒"
	} else {
		return "☐"
	}
}

func FormatDoneInString(done bool) string {
	if done {
		return "DONE"
	} else {
		return "OPEN"
	}
}

func FormatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	} else {
		return date.Format("2006-Jan-02")
	}
}
