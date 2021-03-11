import { NextApiRequest, NextApiResponse } from 'next'
import { sampleMessagesData, sampleMessagesId } from "@/utils/sample-data"

const handler = (req: NextApiRequest, res: NextApiResponse) => {
  const notFound = 404, ok = 200
  const { id } = req.query

  if (id === sampleMessagesId) {
    return res.status(ok).json(sampleMessagesData)
  } else {
    return res.status(notFound).end(`Post: ${id}`)
  }
}

export default handler