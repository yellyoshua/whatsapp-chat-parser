package whatsapp

import (
	"fmt"
	"regexp"
	"time"

	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func formatDate(date string, day_time string, day_time_ampm string) DateFormat {
	var date_split []string = make([]string, 3)
	var day_time_split []string = make([]string, 2)

	regexDate, _ := regexp.Compile("[-/.]")
	date_split = regexDate.Split(date, -1)

	regexTime, _ := regexp.Compile("[:.]")
	day_time_split = regexTime.Split(day_time, -1)
	var hours string   // index 0
	var minutes string // index 1
	var month string   // index 0
	var day string     // index 1
	var year string    // index 2

	if len(day_time_split)-1 >= 0 {
		hours = day_time_split[0]
	}
	if len(day_time_split)-1 >= 1 {
		minutes = day_time_split[1]
	}

	if len(date_split)-1 >= 0 {
		month = date_split[0]
	}
	if len(date_split)-1 >= 1 {
		day = date_split[1]
	}
	if len(date_split)-1 >= 2 {
		year = date_split[2]
	}

	return DateFormat{
		Hours:  utils.PadStart(hours, "0", 2),
		Mins:   utils.PadStart(minutes, "0", 2),
		Month:  utils.StringToInt(month),
		Day:    utils.StringToInt(day),
		Year:   utils.StringToInt(utils.PadStart(year, "2000", 4)),
		Format: day_time_ampm,
		UTC: time.Date(
			utils.StringToInt(year),
			time.Month(utils.StringToInt(month)),
			utils.StringToInt(day),
			utils.StringToInt(hours),
			utils.StringToInt(minutes),
			0,
			0,
			time.UTC,
		).String(),
	}
}

func getTranslateDate(lang string, month int, day int, year int) string {
	t, _ := Translations[lang].(map[string]interface{})

	months, _ := t["months"].(map[int]string)
	date_template, _ := t["date_template"].(string)

	yearFormated := utils.PadStart(utils.IntToString(year), "2000", 4)

	return fmt.Sprintf(date_template, months[month], day, yearFormated)
}
