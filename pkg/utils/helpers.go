package helpers

import "github.com/martinlindhe/notify"

func ShowInfo(title, message string) {
	notify.Notify(title, title, message, "")
}

func ShowError(title, message string) {
	notify.Alert(title, title, message, "")
}

func AskUser(message string) bool {
	return true
}
