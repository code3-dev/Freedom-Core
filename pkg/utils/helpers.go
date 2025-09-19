package helpers

import "github.com/martinlindhe/notify"

func ShowInfo(title, message string) {
	notify.Notify("Freedom Core", title, message, "")
}

func ShowError(title, message string) {
	notify.Alert("Freedom Core", title, message, "")
}

func AskUser(message string) bool {
	return true
}
