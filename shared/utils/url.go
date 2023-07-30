package utils

import (
	"fmt"
	"strings"
)

// CreateURL will create formatted URL standards
func CreateURL(base, checkUUID string) string {
	url := strings.Trim(base, "")

	if url[len(url)-1] != '/' {
		url += "/"
	}

	return fmt.Sprintf("%s%s", url, checkUUID)
}
