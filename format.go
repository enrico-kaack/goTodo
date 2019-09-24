package main

func FormatDone(done bool) string {
	if done {
		return "☒"
	} else {
		return "☐"
	}
}
