import { makeArrayOfMessages, parseMessages } from './parser.ts';
import { Message, ParseStringOptions } from './types.ts';
import { spawn } from "https://deno.land/std@0.95.0/node/child_process.ts"
import { process } from "https://deno.land/std@0.95.0/node/process.ts"
import { readFileSync } from "https://deno.land/std@0.95.0/node/fs.ts"

/**
 * Parses a string containing a WhatsApp chat log.
 *
 * Returns a promise that will contain the parsed messages.
 *
 * @since 1.2.0
 */
const args = process.argv.slice(2)

const pathChatFile = args[0];

const options: ParseStringOptions = { parseAttachments: false, daysFirst: false };

const plainChat = readFileSync(pathChatFile).toString("UTF-8")

const data = await Promise.resolve(plainChat);

const splitedMessages = await data.split(/(?:\r\n|\r|\n)/);

const rawMessages = await makeArrayOfMessages(splitedMessages);

const messages = await parseMessages(rawMessages, options);

const JSON_MESSAGES = JSON.stringify(messages)

console.log(JSON_MESSAGES)