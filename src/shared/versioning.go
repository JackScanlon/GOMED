package shared

import (
	"fmt"
)

var (
	Version = ""
)

func GetVersion() string {
	return fmt.Sprintf("v%s", Version)
}
