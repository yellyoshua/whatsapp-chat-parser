import "tailwindcss/tailwind.css";
import { FC } from "react"

type Props = { Component: FC, pageProps: any }

function MyApp({ Component, pageProps }: Props) {
  return <Component {...pageProps} />
}

export default MyApp