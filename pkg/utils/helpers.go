package helpers

import (
	"github.com/gen2brain/dlgs"
)

func AllowDialog(text string) bool {
	ok, _ := dlgs.Question("پیام", text, true)
	return ok
}
