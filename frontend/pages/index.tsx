import Link from 'next/link'
import Layout from '@/components/Layout'

const IndexPage = () => (
  <Layout title="WhatsApp Book">
    <h1>Hello everyone this is a solution can u export chats from whatsapp and transform to html, pdf or JSON</h1>
    <p>
      <Link href="/about">
        <a>About</a>
      </Link>
    </p>
  </Layout>
)

export default IndexPage
