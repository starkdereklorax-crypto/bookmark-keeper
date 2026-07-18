package handlers

import (
	"strconv"
	"strings"
)

func validURL(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	if strings.Contains(url, " ") {
		return false
	}

	return true
}

func parseID(id string) (int, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return idInt, err
	}

	return idInt, nil
}
