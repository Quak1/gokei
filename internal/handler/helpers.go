package handler

import (
	"errors"
	"net/http"
	"strconv"
)

func readIntParam(r *http.Request, key string) (int, error) {
	idString := r.PathValue(key)

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("Invalid ID")
	}

	return id, nil
}
