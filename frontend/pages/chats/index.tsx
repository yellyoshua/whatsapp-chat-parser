import { GetStaticProps } from 'next'
import { Message } from '../../interfaces'
import Layout from '@/components/Layout'

type Props = {
  messages: Message[]
}

const WithStaticProps = ({ }: Props) => (
  <Layout title="Users List | Next.js + TypeScript Example">
  </Layout>
)

export const getStaticProps: GetStaticProps = async () => {
  // Example for including static props in a Next.js function component page.
  // Don't forget to include the respective types for any props passed into
  // the component.
  const messages: Message[] = []
  return { props: { messages } }
}

export default WithStaticProps
