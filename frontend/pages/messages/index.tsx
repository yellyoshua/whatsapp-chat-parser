import { GetStaticProps } from 'next'
import { Message } from '@/interfaces/index'
import Layout from '@/components/Layout'
import SearchBackup from '@/components/search/searchBackup'

type Props = {
  messages: Message[]
}

// TODO: create a component with a input that find chats backup

const WithStaticProps = ({ }: Props) => {
  const title = "Crear un pdf del chat"

  return <Layout title={title}>
    <SearchBackup />
  </Layout>
}

export const getStaticProps: GetStaticProps = async () => {
  // Example for including static props in a Next.js function component page.
  // Don't forget to include the respective types for any props passed into
  // the component.
  const messages: Message[] = []
  return { props: { messages } }
}

export default WithStaticProps
