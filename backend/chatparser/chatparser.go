package chatparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"rogchap.com/v8go"
)

var jsparser string = "var __defProp = Object.defineProperty;\nvar __name = (target, value) => __defProp(target, \"name\", {value, configurable: true});\nvar __commonJS = (callback, module2) => () => {\n  if (!module2) {\n    module2 = {exports: {}};\n    callback(module2.exports, module2);\n  }\n  return module2.exports;\n};\n" +
	" var require_utils = __commonJS((exports2, module2) => {\n  function indexAboveValue(index, value) {\n    return (array) => array[index] > value;\n  }\n  __name(indexAboveValue, \"indexAboveValue\");\n  function isNegative(number) {\n    return number < 0;\n  }\n  __name(isNegative, \"isNegative\");\n  function groupArrayByValueAtIndex(array, index) {" +
	"return Object.values(array.reduce((obj, item) => {\n      const key = String(\"key_\" + item[index]);\n      if (!obj[key]) {\n        obj[key] = [];\n      }\n      obj[key].push(item);\n      return obj;\n    }, {}));\n  }\n  __name(groupArrayByValueAtIndex, \"groupArrayByValueAtIndex\");\n  module2.exports = {indexAboveValue, isNegative, groupArrayByValueAtIndex};\n});\n" +
	"var require_date = __commonJS((exports2, module2) => {\n  var {\n    indexAboveValue,\n    isNegative,\n    groupArrayByValueAtIndex\n  } = require_utils();\n  function checkAbove12(numericDates) {\n    const daysFirst = numericDates.some(indexAboveValue(0, 12));\n    if (daysFirst)\n      return true;\n    const daysSecond = numericDates.some(indexAboveValue(1, 12));\n    if (daysSecond)\n      return false;\n    return null;\n  }" +
	"__name(checkAbove12, \"checkAbove12\");\n  function checkDecreasing(numericDates) {\n    const datesByYear = groupArrayByValueAtIndex(numericDates, 2);\n    const results = datesByYear.map((dates) => {\n      const daysFirst = dates.slice(1).some((date, i) => {\n        const [first1] = dates[i];\n        const [first2] = date;\n        return isNegative(first2 - first1);\n      });\n      if (daysFirst)\n        return true;" +
	"const daysSecond = dates.slice(1).some((date, i) => {\n        const [, second1] = dates[i];\n        const [, second2] = date;\n        return isNegative(second2 - second1);\n      });\n      if (daysSecond)\n        return false;\n      return null;\n    });\n    const anyTrue = results.some((value) => value === true);\n    if (anyTrue)\n      return true;\n    const anyFalse = results.some((value) => value === false);" +
	"if (anyFalse)\n      return false;\n    return null;\n  }\n  __name(checkDecreasing, \"checkDecreasing\");\n  function changeFrequencyAnalysis(numericDates) {\n    const diffs = numericDates.slice(1).map((date, i) => date.map((num, j) => Math.abs(numericDates[i][j] - num)));\n    const [first, second] = diffs.reduce((total, diff) => {\n      const [first1, second1] = total;\n      const [first2, second2] = diff;" +
	"return [first1 + first2, second1 + second2];\n    }, [0, 0]);\n    if (first > second)\n      return true;\n    if (first < second)\n      return false;\n    return null;\n  }\n  __name(changeFrequencyAnalysis, \"changeFrequencyAnalysis\");\n  function daysBeforeMonths(numericDates) {\n    const firstCheck = checkAbove12(numericDates);\n    if (firstCheck !== null)\n      return firstCheck;" +
	"const secondCheck = checkDecreasing(numericDates);\n    if (secondCheck !== null)\n      return secondCheck;\n    return changeFrequencyAnalysis(numericDates);\n  }\n  __name(daysBeforeMonths, \"daysBeforeMonths\");\n  function normalizeDate(year, month, day) {\n    return [\n      year.padStart(4, \"2000\"),\n      month.padStart(2, \"0\"),\n      day.padStart(2, \"0\")\n    ];\n  }" +
	"__name(normalizeDate, \"normalizeDate\");\n  module2.exports = {\n    checkAbove12,\n    checkDecreasing,\n    changeFrequencyAnalysis,\n    daysBeforeMonths,\n    normalizeDate\n  };\n});\nvar require_time = __commonJS((exports2, module2) => {\n  var regexSplitTime = /[:.]/;\n  function convertTime12to24(time, ampm) {\n    let [hours, minutes, seconds] = time.split(regexSplitTime);" +
	"if (hours === \"12\")\n      hours = \"00\";\n    if (ampm === \"PM\")\n      hours = parseInt(hours, 10) + 12;\n    var fullseconds = seconds ? String(\":\" + seconds) : \"\";\n    var fulltime = hours + \":\" + minutes + fullseconds;\n    return String(fulltime);\n  }\n  __name(convertTime12to24, \"convertTime12to24\");\n  function normalizeTime(time) {\n    const [hours, minutes, seconds] = time.split(regexSplitTime);" +
	"var fullseconds = seconds || \"00\";\n    if (hours.length == 1) {\n      return String(\"0\" + hours + \":\" + minutes + \":\" + fullseconds);\n    }\n    return String(hours + \":\" + minutes + \":\" + fullseconds);\n  }\n  __name(normalizeTime, \"normalizeTime\");\n  function normalizeAMPM(ampm) {\n    return ampm.replace(/[^apm]/gi, \"\").toUpperCase();\n  }\n  __name(normalizeAMPM, \"normalizeAMPM\");" +
	"module2.exports = {\n    regexSplitTime,\n    convertTime12to24,\n    normalizeTime,\n    normalizeAMPM\n  };\n});\nvar require_parser = __commonJS((exports2, module2) => {\n  var {daysBeforeMonths, normalizeDate} = require_date();\n  var {\n    regexSplitTime,\n    convertTime12to24,\n    normalizeAMPM,\n    normalizeTime\n  } = require_time();" +
	"var regexParser = /^(?:\\u200E|\\u200F)*\\[?(\\d{1,4}[-/.] ?\\d{1,4}[-/.] ?\\d{1,4})[,.]? \\D*?(\\d{1,2}[.:]\\d{1,2}(?:[.:]\\d{1,2})?)(?: ([ap]\\.? ?m\\.?))?\\]?(?: -|:)? (.+?): ([^]*)/i;\n  var regexParserSystem = /^(?:\\u200E|\\u200F)*\\[?(\\d{1,4}[-/.] ?\\d{1,4}[-/.] ?\\d{1,4})[,.]? \\D*?(\\d{1,2}[.:]\\d{1,2}(?:[.:]\\d{1,2})?)(?: ([ap]\\.? ?m\\.?))?\\]?(?: -|:)? ([^]+)/i;\n  var regexSplitDate = /[-/.] ?/;\n  var regexAttachment = /<.+:(.+)>/;" +
	"function makeArrayOfMessages2(lines) {\n    return lines.reduce((acc, line) => {\n      if (!regexParser.test(line)) {\n        if (regexParserSystem.test(line)) {\n          acc.push({system: true, msg: line});\n        } else if (typeof acc[acc.length - 1] !== \"undefined\") {\n          const prevMessage = acc.pop();\n          acc.push({\n            system: prevMessage.system,\n            msg: String(prevMessage.msg + \"\\n\" + line)\n          });\n        }\n      }" +
	"else {\n        acc.push({system: false, msg: line});\n      }\n      return acc;\n    }, []);\n  }\n  __name(makeArrayOfMessages2, \"makeArrayOfMessages\");\n  function parseMessageAttachment(message) {\n    const attachmentMatch = message.match(regexAttachment);\n    if (attachmentMatch)\n      return {fileName: attachmentMatch[1].trim()};\n    return null;\n  }\n  __name(parseMessageAttachment, \"parseMessageAttachment\");" +
	"function parseMessages2(messages, options = {daysFirst: void 0, parseAttachments: false}) {\n    const sortByLengthAsc = /* @__PURE__ */ __name((a, b) => a.length - b.length, \"sortByLengthAsc\");\n    let {daysFirst} = options;\n    const {parseAttachments} = options;\n    const parsed = messages.map((obj) => {\n      const {system, msg} = obj;\n      if (system) {\n        const [, date2, time2, ampm2, message2] = regexParserSystem.exec(msg);" +
	"return {date: date2, time: time2, ampm: ampm2 || null, author: \"System\", message: message2};\n      }\n      const [, date, time, ampm, author, message] = regexParser.exec(msg);\n      return {date, time, ampm: ampm || null, author, message};\n    });\n    if (typeof daysFirst !== \"boolean\") {\n      const numericDates = Array.from(new Set(parsed.map(({date}) => date)), (date) => date.split(regexSplitDate).sort(sortByLengthAsc).map(Number));" +
	"daysFirst = daysBeforeMonths(numericDates);\n    }\n    return parsed.map(({date, time, ampm, author, message}) => {\n      let day;\n      let month;\n      let year;\n      const splitDate = date.split(regexSplitDate).sort(sortByLengthAsc);\n      if (daysFirst === false) {\n        [month, day, year] = splitDate;\n      } else {\n        [day, month, year] = splitDate;\n      }\n      [year, month, day] = normalizeDate(year, month, day);" +
	"const [hours, minutes, seconds] = normalizeTime(ampm ? convertTime12to24(time, normalizeAMPM(ampm)) : time).split(regexSplitTime);\n      const finalObject = {\n        date: new Date(year, month - 1, day, hours, minutes, seconds),\n        author,\n        message\n      };\n      if (parseAttachments) {\n        const attachment = parseMessageAttachment(message);\n        if (attachment)\n          finalObject.attachment = attachment;\n      }" +
	"return finalObject;\n    });\n  }\n  __name(parseMessages2, \"parseMessages\");\n  module2.exports = {\n    makeArrayOfMessages: makeArrayOfMessages2,\n    parseMessages: parseMessages2\n  };\n});\nvar {makeArrayOfMessages, parseMessages} = require_parser();\nfunction parseString(plainMessages, options) {\n  const messages = plainMessages.split(/(?:\\r\\n|\\r|\\n)/);\n  const fullMessages = makeArrayOfMessages(messages);" +
	"const parsedMessages = parseMessages(fullMessages, options);\n  return parsedMessages;\n}\n__name(parseString, \"parseString\");"

// RegexContact format input `$date - Carlos perez: $message`
var RegexContact string = `(\d{1,2}/\d{1,2}/\d{2,4})+(, )[0-9:]+(.+?)(: )`

// RegexAttachment format input `$date - $contact: IMG-20200319-WA0011.jpg (file attached)`
var RegexAttachment string = `(: )+[\S\s]+(\.\w{2,4}\s)+\(+(file attached)+\)`

// RegexTextAttachment format input `IMG-20200319-WA0011.jpg (file attached)`
var RegexTextAttachment string = `\(file attached\)`

// Parser _
type Parser interface {
	ParserMessages(data []byte, outputMessages *string) error
}

type parserstruct struct {
	v8 *v8go.Context
}

// New __
func New() Parser {
	ctx, err := v8go.NewContext(nil)
	if err != nil {
		logger.Fatal("Error creating parser class -> %s", err)
	}

	if _, err := regexp.Compile(RegexTextAttachment); err != nil {
		logger.Fatal("Error regex TextAttachment -> %s", err)
	}
	if _, err := regexp.Compile(RegexContact); err != nil {
		logger.Fatal("Error regex Contact -> %s", err)
	}
	if _, err := regexp.Compile(RegexAttachment); err != nil {
		logger.Fatal("Error regex Attachment -> %s", err)
	}

	return &parserstruct{v8: ctx}
}

// ParserMessages _
func (p *parserstruct) ParserMessages(data []byte, outputMessages *string) error {
	v8 := p.v8

	plainMessages := byteToStringMessages(data)

	if _, err := v8.RunScript(jsparser, "loader.js"); err != nil {
		return fmt.Errorf("Error RunScript from (jsparser) -> %s", err)
	}

	scriptMessages := fmt.Sprintf("var plainMessages=`%s`", plainMessages)
	if _, err := v8.RunScript(scriptMessages, "index.js"); err != nil {
		return fmt.Errorf("Error RunScript from (plainMessages) -> %s", err)
	}

	scriptParseMessages := `
		var options = {parseAttachments: true, daysFirst: false};
		var messages = parseString(plainMessages, options); 
		JSON.stringify(messages, null, ' ')
	`
	val, err := v8.RunScript(scriptParseMessages, "parser.js")
	if err != nil {
		return fmt.Errorf("Error RunScript from (parseString) -> %s", err)
	}

	*outputMessages = val.String()

	return nil
}

func byteToStringMessages(data []byte) string {
	var messages string

	regexContact, _ := regexp.Compile(RegexContact)

	var whatsappMessages []string

	plainMessages := strings.TrimSpace(string(data))
	bytesOfMessages := []byte(plainMessages)
	messagesIndexes := regexContact.FindAllStringIndex(plainMessages, -1)

	for i := 0; i < len(messagesIndexes); i++ {
		axis := messagesIndexes[i]
		nextIndex := i + 1
		existMessage := len(axis) == 2

		if existMessage {
			start, _ := axis[0], axis[1]

			if nextIndex < len(messagesIndexes) {
				nextAxis := messagesIndexes[nextIndex]
				existNextAxis := len(nextAxis) == 2

				if existNextAxis {
					message := string(bytesOfMessages[start:nextAxis[0]])
					message = strings.TrimSpace(replaceAttachment(message))
					whatsappMessages = append(whatsappMessages, message)
				}

			} else {
				message := string(bytesOfMessages[start:])
				message = strings.TrimSpace(replaceAttachment(message))
				whatsappMessages = append(whatsappMessages, message)
			}

		}
	}

	messages = strings.Join(whatsappMessages, "\n")

	return messages
}

func replaceAttachment(message string) string {
	regexAttachment, _ := regexp.Compile(RegexAttachment)
	regexTextAttachment, _ := regexp.Compile(RegexTextAttachment)

	attachment := regexTextAttachment.ReplaceAllString(regexAttachment.FindString(message), "${1}$2")

	attachmentBytes := []byte(attachment)

	if len(attachmentBytes) == 0 {
		return message
	}

	fileName := strings.TrimSpace(string(attachmentBytes[1:]))

	fileName = strings.ReplaceAll(fileName, " ", "%20")

	repl := fmt.Sprintf(": <attached: %s>", fileName)
	result := regexAttachment.ReplaceAllString(message, repl)
	return result
}
