package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDate(t *testing.T) {
	expected := DateFormat{Hours: "00", Mins: "00", Format: "", Day: 1, Month: 1, Year: 2000, UTC: "2000-01-01 00:00:00 +0000 UTC"}
	date := formatDate("", "", "")
	assert.Equal(t, expected, date)

	expected1 := DateFormat{Hours: "14", Mins: "10", Format: "", Day: 20, Month: 12, Year: 2021, UTC: "2021-12-20 14:10:00 +0000 UTC"}
	date1 := formatDate("12/20/21", "14:10", "")
	assert.Equal(t, expected1, date1)

	expected2 := DateFormat{Hours: "04", Mins: "10", Format: "am", Day: 20, Month: 12, Year: 2021, UTC: "2021-12-20 04:10:00 +0000 UTC"}
	date2 := formatDate("12/20/2021", "4:10", "am")
	assert.Equal(t, expected2, date2)

}

func TestGetTranslateDate(t *testing.T) {
	expected := "Marzo 1, 2001"
	dateParsed := getTranslateDate("es", 3, 1, 2001)
	assert.Equal(t, expected, dateParsed)

	expected1 := "Marzo 1, 2001"
	dateParsed1 := getTranslateDate("es", 3, 1, 1)
	assert.Equal(t, expected1, dateParsed1)
}
