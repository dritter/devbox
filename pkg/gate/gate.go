package gate

import (
	"os"
	"strconv"
)

func Flakes() bool {
	val, err := strconv.ParseBool(os.Getenv("DEVBOX_FLAKES"))
	if err != nil {
		return false
	}
	return val
}
