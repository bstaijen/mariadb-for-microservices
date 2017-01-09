package helper

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// PaginationFromRequest return the offset and rows from a request. If not
// available then return the default of 1, 10. Return in order: offset, rows.
func PaginationFromRequest(r *http.Request) (int, int) {
	// get query param offset
	var offsetString = r.URL.Query().Get("offset")

	// get query param rows
	var rowsString = r.URL.Query().Get("rows")

	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		offset = 1
		logrus.Warn("offset (" + offsetString + ") is not a number")
	}
	rows, err := strconv.Atoi(rowsString)
	if err != nil {
		rows = 10
		logrus.Warn("rows (" + rowsString + ") is not a number")
	}
	return offset, rows
}
