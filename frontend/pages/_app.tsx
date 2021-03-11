import "tailwindcss/tailwind.css";
import type { AppProps } from 'next/app'
import AppProvider from '@/context/contextApp'
import { APP_NAME } from "@/config/app"


function MyApp({ Component, pageProps }: AppProps) {
  return <AppProvider state={{ name: APP_NAME, API_URL: process.env.NEXT_PUBLIC_API_URI }}>
    <Component {...pageProps} />
  </AppProvider>
}

export default MyApp