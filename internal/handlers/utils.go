package handlers

import "strings"

func IsBlankOrInvisible(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
