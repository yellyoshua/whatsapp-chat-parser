package whatsapp

import "strings"

var (
	TRANSLA_ES_KEY = "es"
	TRANSLA_EN_KEY = "en"

	TRANSLA_MONTH_KEY = "months"
)

var LANGUAGES_AVALIABLE = []string{
	TRANSLA_EN_KEY,
	TRANSLA_ES_KEY,
}

var Translations = map[string]interface{}{
	TRANSLA_ES_KEY: map[string]interface{}{
		TRANSLA_MONTH_KEY: map[int]string{
			0:  "Desconocido",
			1:  "Enero",
			2:  "Febrero",
			3:  "Marzo",
			4:  "Abril",
			5:  "Mayo",
			6:  "Junio",
			7:  "Julio",
			8:  "Agosto",
			9:  "Septiembre",
			10: "Octubre",
			11: "Noviembre",
			12: "Diciembre",
		},
		"date_template": "%s %v, %v", // Enero 2, 2021
	},
	TRANSLA_EN_KEY: map[string]interface{}{
		TRANSLA_MONTH_KEY: map[int]string{
			0:  "Desconocido",
			1:  "January",
			2:  "February",
			3:  "March",
			4:  "April",
			5:  "May",
			6:  "June",
			7:  "July",
			8:  "August",
			9:  "September",
			10: "October",
			11: "November",
			12: "December",
		},
		"date_template": "%s %v, %v", // January 2, 2021
	},
}

// If the language is not supported load ES language as default
func safeLangSupported(lang string) (language string) {
	var isSupported = false
	for _, l := range LANGUAGES_AVALIABLE {
		if support := strings.Contains(lang, l); support {
			isSupported = true
			break
		}
	}

	if isSupported {
		return lang
	} else {
		return TRANSLA_ES_KEY
	}
}
