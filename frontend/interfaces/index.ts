// You can include shared interfaces/types in a separate file
// and then use them in any component by importing them. For
// example, to import the interface below do:
//
// import { Message } from 'path/to/interfaces';

export type Message = {
  date: string | Date,
  author: string,
  isSender: boolean,
  isInfo: boolean,
  isReceiver: boolean,
  message: string,
  attachment?: Attachment
}

export type Attachment = {
  fileName: string,
}