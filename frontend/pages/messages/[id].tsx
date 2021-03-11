import { GetStaticProps, GetStaticPaths } from 'next'
import { FC } from "react"
import { Message } from '@/interfaces/index'
import Layout from '@/components/Layout'
import { APP_NAME } from '@/config/app'

type Props = { messages?: Message[], errors?: string, id: string }

const StaticPropsDetail: FC<Props> = ({ errors, messages, id }) => {
  const title = `${APP_NAME} | ${id ? id : "Buscando"}`
  const notFound = `${APP_NAME} | ${errors ? errors : "Buscando"}`

  if (errors) {
    return (
      <Layout title={notFound} showCreateButton>
        <p>You try find this - {id} messages</p>
      </Layout>
    )
  }

  return (
    <Layout title={title} showCreateButton>
      <h1>Hello world</h1>
      <div>
        {messages?.map(({ author, date }, i) => (
          <div key={i}>
            <h2>{author}</h2>
            <h3>{date}</h3>
          </div>
        ))}
      </div>
    </Layout>
  )
}

export default StaticPropsDetail

export const getStaticPaths: GetStaticPaths = async () => {
  // Get the paths we want to pre-render based on users

  // We'll pre-render only these paths at build time.
  // { fallback: false } means other routes should 404.

  return { paths: [], fallback: true }
}

// This function gets called at build time on server-side.
// It won't be called on client-side, so you can even do
// direct database queries.
export const getStaticProps: GetStaticProps = async ({ params }) => {
  const id = params?.id

  try {

    const API_URL = process.env.NEXT_PUBLIC_API_URI

    const request = await fetch(`${API_URL}/api/messages/${id}`)

    if (request.status === 200) {
      const messages: Message[] = await request.json()
      return { props: { messages, id } }
    }

    return { props: { messages: [], id, errors: "No encontrado" } }
  } catch (err) {
    return { props: { messages: [], id, errors: err.message } }
  }
}
