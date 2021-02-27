var __defProp = Object.defineProperty;
var __name = (target, value) => __defProp(target, "name", {value, configurable: true});
var __commonJS = (callback, module2) => () => {
  if (!module2) {
    module2 = {exports: {}};
    callback(module2.exports, module2);
  }
  return module2.exports;
};

// src/utils.js
var require_utils = __commonJS((exports2, module2) => {
  function indexAboveValue(index, value) {
    return (array) => array[index] > value;
  }
  __name(indexAboveValue, "indexAboveValue");
  function isNegative(number) {
    return number < 0;
  }
  __name(isNegative, "isNegative");
  function groupArrayByValueAtIndex(array, index) {
    return Object.values(array.reduce((obj, item) => {
      const key = String("key_" + item[index]);
      if (!obj[key]) {
        obj[key] = [];
      }
      obj[key].push(item);
      return obj;
    }, {}));
  }
  __name(groupArrayByValueAtIndex, "groupArrayByValueAtIndex");
  module2.exports = {indexAboveValue, isNegative, groupArrayByValueAtIndex};
});

// src/date.js
var require_date = __commonJS((exports2, module2) => {
  var {
    indexAboveValue,
    isNegative,
    groupArrayByValueAtIndex
  } = require_utils();
  function checkAbove12(numericDates) {
    const daysFirst = numericDates.some(indexAboveValue(0, 12));
    if (daysFirst)
      return true;
    const daysSecond = numericDates.some(indexAboveValue(1, 12));
    if (daysSecond)
      return false;
    return null;
  }
  __name(checkAbove12, "checkAbove12");
  function checkDecreasing(numericDates) {
    const datesByYear = groupArrayByValueAtIndex(numericDates, 2);
    const results = datesByYear.map((dates) => {
      const daysFirst = dates.slice(1).some((date, i) => {
        const [first1] = dates[i];
        const [first2] = date;
        return isNegative(first2 - first1);
      });
      if (daysFirst)
        return true;
      const daysSecond = dates.slice(1).some((date, i) => {
        const [, second1] = dates[i];
        const [, second2] = date;
        return isNegative(second2 - second1);
      });
      if (daysSecond)
        return false;
      return null;
    });
    const anyTrue = results.some((value) => value === true);
    if (anyTrue)
      return true;
    const anyFalse = results.some((value) => value === false);
    if (anyFalse)
      return false;
    return null;
  }
  __name(checkDecreasing, "checkDecreasing");
  function changeFrequencyAnalysis(numericDates) {
    const diffs = numericDates.slice(1).map((date, i) => date.map((num, j) => Math.abs(numericDates[i][j] - num)));
    const [first, second] = diffs.reduce((total, diff) => {
      const [first1, second1] = total;
      const [first2, second2] = diff;
      return [first1 + first2, second1 + second2];
    }, [0, 0]);
    if (first > second)
      return true;
    if (first < second)
      return false;
    return null;
  }
  __name(changeFrequencyAnalysis, "changeFrequencyAnalysis");
  function daysBeforeMonths(numericDates) {
    const firstCheck = checkAbove12(numericDates);
    if (firstCheck !== null)
      return firstCheck;
    const secondCheck = checkDecreasing(numericDates);
    if (secondCheck !== null)
      return secondCheck;
    return changeFrequencyAnalysis(numericDates);
  }
  __name(daysBeforeMonths, "daysBeforeMonths");
  function normalizeDate(year, month, day) {
    return [
      year.padStart(4, "2000"),
      month.padStart(2, "0"),
      day.padStart(2, "0")
    ];
  }
  __name(normalizeDate, "normalizeDate");
  module2.exports = {
    checkAbove12,
    checkDecreasing,
    changeFrequencyAnalysis,
    daysBeforeMonths,
    normalizeDate
  };
});

// src/time.js
var require_time = __commonJS((exports2, module2) => {
  var regexSplitTime = /[:.]/;
  function convertTime12to24(time, ampm) {
    let [hours, minutes, seconds] = time.split(regexSplitTime);
    if (hours === "12")
      hours = "00";
    if (ampm === "PM")
      hours = parseInt(hours, 10) + 12;
    var fullseconds = seconds ? String(":" + seconds) : "";
    var fulltime = hours + ":" + minutes + fullseconds;
    return String(fulltime);
  }
  __name(convertTime12to24, "convertTime12to24");
  function normalizeTime(time) {
    const [hours, minutes, seconds] = time.split(regexSplitTime);
    var fullseconds = seconds || "00";
    if (hours.length == 1) {
      return String("0" + hours + ":" + minutes + ":" + fullseconds);
    }
    return String(hours + ":" + minutes + ":" + fullseconds);
  }
  __name(normalizeTime, "normalizeTime");
  function normalizeAMPM(ampm) {
    return ampm.replace(/[^apm]/gi, "").toUpperCase();
  }
  __name(normalizeAMPM, "normalizeAMPM");
  module2.exports = {
    regexSplitTime,
    convertTime12to24,
    normalizeTime,
    normalizeAMPM
  };
});

// src/parser.js
var require_parser = __commonJS((exports2, module2) => {
  var {daysBeforeMonths, normalizeDate} = require_date();
  var {
    regexSplitTime,
    convertTime12to24,
    normalizeAMPM,
    normalizeTime
  } = require_time();
  var regexParser = /^(?:\u200E|\u200F)*\[?(\d{1,4}[-/.] ?\d{1,4}[-/.] ?\d{1,4})[,.]? \D*?(\d{1,2}[.:]\d{1,2}(?:[.:]\d{1,2})?)(?: ([ap]\.? ?m\.?))?\]?(?: -|:)? (.+?): ([^]*)/i;
  var regexParserSystem = /^(?:\u200E|\u200F)*\[?(\d{1,4}[-/.] ?\d{1,4}[-/.] ?\d{1,4})[,.]? \D*?(\d{1,2}[.:]\d{1,2}(?:[.:]\d{1,2})?)(?: ([ap]\.? ?m\.?))?\]?(?: -|:)? ([^]+)/i;
  var regexSplitDate = /[-/.] ?/;
  var regexAttachment = /<.+:(.+)>/;
  function makeArrayOfMessages2(lines) {
    return lines.reduce((acc, line) => {
      if (!regexParser.test(line)) {
        if (regexParserSystem.test(line)) {
          acc.push({system: true, msg: line});
        } else if (typeof acc[acc.length - 1] !== "undefined") {
          const prevMessage = acc.pop();
          acc.push({
            system: prevMessage.system,
            msg: String(prevMessage.msg + "\n" + line)
          });
        }
      } else {
        acc.push({system: false, msg: line});
      }
      return acc;
    }, []);
  }
  __name(makeArrayOfMessages2, "makeArrayOfMessages");
  function parseMessageAttachment(message) {
    const attachmentMatch = message.match(regexAttachment);
    if (attachmentMatch)
      return {fileName: attachmentMatch[1].trim()};
    return null;
  }
  __name(parseMessageAttachment, "parseMessageAttachment");
  function parseMessages2(messages, options = {daysFirst: void 0, parseAttachments: true}) {
    const sortByLengthAsc = /* @__PURE__ */ __name((a, b) => a.length - b.length, "sortByLengthAsc");
    let {daysFirst} = options;
    const {parseAttachments} = options;
    const parsed = messages.map((obj) => {
      const {system, msg} = obj;
      if (system) {
        const [, date2, time2, ampm2, message2] = regexParserSystem.exec(msg);
        return {date: date2, time: time2, ampm: ampm2 || null, author: "System", message: message2};
      }
      const [, date, time, ampm, author, message] = regexParser.exec(msg);
      return {date, time, ampm: ampm || null, author, message};
    });
    if (typeof daysFirst !== "boolean") {
      const numericDates = Array.from(new Set(parsed.map(({date}) => date)), (date) => date.split(regexSplitDate).sort(sortByLengthAsc).map(Number));
      daysFirst = daysBeforeMonths(numericDates);
    }
    return parsed.map(({date, time, ampm, author, message}) => {
      let day;
      let month;
      let year;
      const splitDate = date.split(regexSplitDate).sort(sortByLengthAsc);
      if (daysFirst === false) {
        [month, day, year] = splitDate;
      } else {
        [day, month, year] = splitDate;
      }
      [year, month, day] = normalizeDate(year, month, day);
      const [hours, minutes, seconds] = normalizeTime(ampm ? convertTime12to24(time, normalizeAMPM(ampm)) : time).split(regexSplitTime);
      const finalObject = {
        date: new Date(year, month - 1, day, hours, minutes, seconds),
        author,
        message
      };
      if (parseAttachments) {
        const attachment = parseMessageAttachment(message);
        if (attachment)
          finalObject.attachment = attachment;
      }
      return finalObject;
    });
  }
  __name(parseMessages2, "parseMessages");
  module2.exports = {
    makeArrayOfMessages: makeArrayOfMessages2,
    parseMessages: parseMessages2
  };
});

// src/index.js
var {makeArrayOfMessages, parseMessages} = require_parser();
function parseString(plainMessages, options) {
  const messages = plainMessages.split(/(?:\r\n|\r|\n)/);
  const fullMessages = makeArrayOfMessages(messages);
  const parsedMessages = parseMessages(fullMessages, options);
  return parsedMessages;
}
__name(parseString, "parseString");
module.exports = {parseString};
