import es from "./dayjs/locale/es.js"
import relativeTime from "./dayjs/plugin/relativeTime/index.js"
import { makeArrayOfMessages, parseMessages } from './parser.ts'
import { ParseStringOptions } from './types.ts'
import { process } from "https://deno.land/std@0.95.0/node/process.ts"
import { readFileSync } from "https://deno.land/std@0.95.0/node/fs.ts"
import * as emoji from "https://deno.land/x/emoji/mod.ts"
import dayjs from "./dayjs/index.js"

// dayjs.extend(relativeTime)

/**
 * Parses a string containing a WhatsApp chat log.
 *
 * Returns a promise that will contain the parsed messages.
 *
 * @since 1.2.0
 */
const args = process.argv.slice(2);

const pathChatFile = args[0];

dayjs.locale("es", es)

if (pathChatFile === "--is-ok") {
  console.log("ok")
} else {
  const options: ParseStringOptions = { parseAttachments: true, daysFirst: false };

  const plainChat = readFileSync(pathChatFile).toString("UTF-8")

  const data = await Promise.resolve(plainChat);

  const splitedMessages = await data.split(/(?:\r\n|\r|\n)/);

  const rawMessages = await makeArrayOfMessages(splitedMessages);

  const messages = await parseMessages(rawMessages, options);

  const parsedMessages = messages.map(m => ({
    ...m,
    message: emoji.emojify(emoji.unemojify(m.message)),
    date: dayjs(m.date).format("MM_DD_YYYY=HH:mm")
  }))

  const JSON_MESSAGES = JSON.stringify(parsedMessages)

  console.log(JSON_MESSAGES)
}