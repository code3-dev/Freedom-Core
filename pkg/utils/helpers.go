package helpers

import (
	"github.com/sqweek/dialog"
)

func AllowDialog(text string) bool {
	ok := dialog.Message(text).Title("پیام").YesNo()
	return ok
}
