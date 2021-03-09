import { Message } from '../interfaces'

const sender = { name: "Frank Sinatra", is: true, no: false };

const receiver = { name: "Frank Sinatra", is: true, no: false };

const noAttachment = undefined;

/** Dummy messages data. */
export const sampleMessagesData: Message[] = [
  {
    attachment: noAttachment,
    author: receiver.name,
    date: new Date(Date.now() + 1).toUTCString(),
    message: "",
    isInfo: receiver.no,
    isReceiver: receiver.is,
    isSender: receiver.no
  },
  {
    attachment: noAttachment,
    author: sender.name,
    date: new Date(Date.now() + 1).toUTCString(),
    message: "",
    isInfo: sender.no,
    isReceiver: sender.no,
    isSender: sender.is
  },
  {
    attachment: noAttachment,
    author: sender.name,
    date: new Date(Date.now() + 1).toUTCString(),
    message: "",
    isInfo: sender.no,
    isReceiver: sender.no,
    isSender: sender.is
  },
  {
    attachment: noAttachment,
    author: receiver.name,
    date: new Date(Date.now() + 1).toUTCString(),
    message: "",
    isInfo: receiver.no,
    isReceiver: receiver.is,
    isSender: receiver.no
  }
]
