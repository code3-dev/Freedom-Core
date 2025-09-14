package helpers

import "github.com/gen2brain/dlgs"

func AllowDialog(text string) bool {
	ok, _ := dlgs.Question("پیام", text, true)
	return ok
}

func InfoDialog(title, text string) {
	dlgs.Info(title, text)
}

func ErrorDialog(title, text string) {
	dlgs.Error(title, text)
}
