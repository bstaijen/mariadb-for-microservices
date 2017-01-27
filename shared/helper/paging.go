package helper

import (
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// PaginationFromRequest return the offset and rows from a request. If not
// available then return the default of 1, 10. Returns in following order: offset, rows.
func PaginationFromRequest(r *http.Request) (int, int) {

	// Keep compiler happy
	var err error

	// offset
	var offsetString = r.URL.Query().Get("offset")
	var offset int
	if len(offsetString) > 0 {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			offset = 0
			logrus.Warnf("offset (%v) is not a number", offsetString)
		}
	}
	// Prevent negative offset number
	if offset < 0 {
		offset = 0
	}

	// rows
	var rowsString = r.URL.Query().Get("rows")
	var rows int
	if len(rowsString) > 0 {
		rows, err = strconv.Atoi(rowsString)
		if err != nil {
			rows = 10
			logrus.Warnf("rows (%v) is not a number", rowsString)
		}
	}
	// Prevent negative rows number
	if rows < 1 {
		rows = 10
	}

	// return
	return offset, rows
}
