import { NextApiRequest, NextApiResponse } from 'next'
import { sampleMessagesData } from "utils/sample-data"

const handler = (req: NextApiRequest, res: NextApiResponse) => {
  const methodNotAllowed = 405, ok = 200, internalError = 500;

  if (req.method === "GET") {
    try {
      if (!Array.isArray([])) {
        throw new Error('Cannot find messages')
      }

      res.status(ok).json(sampleMessagesData)
    } catch (err) {
      res.status(internalError).json({ statusCode: internalError, message: err.message })
    }
  } else {
    res.status(methodNotAllowed)
  }
}

export default handler
