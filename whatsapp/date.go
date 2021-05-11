package whatsapp

import (
	"fmt"
	"regexp"
	"time"

	"github.com/yellyoshua/whatsapp-chat-parser/utils"
)

func formatDate(date string, day_time string, day_time_ampm string) DateFormat {
	regexDate, _ := regexp.Compile("[-/.]")
	regexTime, _ := regexp.Compile("[:.]")

	date_split := utils.SafeStringArray(regexDate.Split(date, -1), 3)
	day_time_split := utils.SafeStringArray(regexTime.Split(day_time, -1), 2)

	month := utils.PadStart(date_split[0], "01", 2)
	day := utils.PadStart(date_split[1], "01", 2)
	year := utils.PadStart(date_split[2], "2000", 4)

	hours := utils.PadStart(day_time_split[0], "0", 2)
	minutes := utils.PadStart(day_time_split[1], "0", 2)

	return DateFormat{
		Hours:  hours,
		Mins:   minutes,
		Month:  utils.StringToInt(month),
		Day:    utils.StringToInt(day),
		Year:   utils.StringToInt(year),
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
	var language = safeLangSupported(lang)

	t, _ := Translations[language].(map[string]interface{})

	months, _ := t[TRANSLA_MONTH_KEY].(map[int]string)
	date_template, _ := t["date_template"].(string)

	yearFormated := utils.PadStart(utils.IntToString(year), "2000", 4)

	return fmt.Sprintf(date_template, months[month], day, yearFormated)
}
