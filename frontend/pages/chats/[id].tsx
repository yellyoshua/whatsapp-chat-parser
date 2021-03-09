import { GetStaticProps, GetStaticPaths } from 'next'

import { Message } from '../../interfaces'
import Layout from '@/components/Layout'
import { API_URL } from 'config/env'

type Props = {
  messages?: Message[]
  errors?: string
}

const StaticPropsDetail = ({ errors }: Props) => {
  if (errors) {
    return (
      <Layout title="Not founded message">
        <p>
          <span style={{ color: 'red' }}>Error:</span> {errors}
        </p>
      </Layout>
    )
  }

  return (
    <Layout
      title="User Detail">
      <h1>Hello world</h1>
    </Layout>
  )
}

export default StaticPropsDetail

export const getStaticPaths: GetStaticPaths = async () => {
  // Get the paths we want to pre-render based on users

  // We'll pre-render only these paths at build time.
  // { fallback: false } means other routes should 404.
  return { paths: [], fallback: false }
}

// This function gets called at build time on server-side.
// It won't be called on client-side, so you can even do
// direct database queries.
export const getStaticProps: GetStaticProps = async ({ params }) => {
  try {
    const id = params?.id

    const request = await fetch(`${API_URL}/chats/${id}`)
    const messages: Message[] = await request.json()

    if (request.status === 200) {
      return { props: { messages } }
    }

    return { props: { messages: [], errors: "Exist error" } }
  } catch (err) {
    return { props: { messages: [], errors: err.message } }
  }
}
