package wx

import (
	"github.com/songcser/gingo/pkg/admin"
)

func Admin() {
	var a TWxUsers
	admin.New(a, "app", "应用")
}
