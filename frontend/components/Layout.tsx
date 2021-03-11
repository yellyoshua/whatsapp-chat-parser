import React, { ReactNode, Fragment, FC, useContext } from 'react'
import Link from 'next/link'
import Head from 'next/head'
import { APP_NAME } from "@/config/app"
import CreateButton from './button/createButton'
import { AppContext } from '@/context/contextApp'


type Props = { children?: ReactNode, title?: string, showCreateButton?: boolean }

const Header: FC<{ pageTitle: string, showCreateButton: boolean }> = ({ showCreateButton, pageTitle }) => (
  <header className="bg-black w-full flex relative justify-between items-center mx-auto px-8 h-20">
    <Link href="/">
      <div className="cursor-pointer text-2xl text-white font-semibold inline-flex items-center">
        <span>{pageTitle}</span>
      </div>
    </Link>
    <div>
      <ul className="flex text-white">
        <li className="ml-5 px-2 py-1">
          <Link href="/about">
            <a>Acerca</a>
          </Link>
        </li>
        {showCreateButton && <CreateButton />}
      </ul>
    </div>
  </header>
)

const Footer: FC = () => (
  <footer className="bg-black bg-opacity-95 pt-10 sm:mt-10">
    <div className="max-w-6xl m-auto text-gray-800 flex flex-wrap justify-between">

      <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
        <div className="text-xs uppercase text-gray-400 font-medium mb-6">
          Comenzando
            </div>

        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          Instalaci&oacute;n
            </a>
      </div>

      <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
        <div className="text-xs uppercase text-gray-400 font-medium mb-6">
          Personalizaci&oacute;n
            </div>

        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          Configuraci&oacute;n
            </a>
        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          Configuraci&oacute;n del tema
            </a>
      </div>

      <div className="p-5 w-1/2 sm:w-4/12 md:w-3/12">
        <div className="text-xs uppercase text-gray-400 font-medium mb-6">
          Proyecto
            </div>

        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          GitHub
            </a>
        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          Discord
            </a>
        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          Twitter
            </a>
        <a href="#" className="my-3 block text-gray-300 hover:text-gray-100 text-sm font-medium duration-700">
          YouTube
            </a>
      </div>
    </div>

    <div className="pt-2">
      <div className="flex pb-5 px-3 m-auto pt-5 
            border-t border-gray-500 text-gray-400 text-sm 
            flex-col md:flex-row max-w-6xl">
        <div className="mt-2">
          © Copyright 1998-year. All Rights Reserved.
            </div>

        <div className="md:flex-auto md:flex-row-reverse mt-2 flex-row flex">
          <a href="#" className="w-6 mx-1">
            <i className="uil uil-facebook-f"></i>
          </a>
          <a href="#" className="w-6 mx-1">
            <i className="uil uil-twitter-alt"></i>
          </a>
          <a href="#" className="w-6 mx-1">
            <i className="uil uil-youtube"></i>
          </a>
          <a href="#" className="w-6 mx-1">
            <i className="uil uil-linkedin"></i>
          </a>
          <a href="#" className="w-6 mx-1">
            <i className="uil uil-instagram"></i>
          </a>
        </div>
      </div>
    </div>
  </footer>
)

const Layout: FC<Props> = ({ children, title = APP_NAME, showCreateButton = false }) => {
  const [app] = useContext(AppContext)

  return (
    <Fragment>
      <Head>
        <title>{title}</title>
        <meta charSet="utf-8" />
        <meta name="viewport" content="initial-scale=1.0, width=device-width" />
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@iconscout/unicons@3.0.6/css/line.css"></link>
      </Head>
      {app?.name && (
        <Header pageTitle={app?.name} showCreateButton={showCreateButton} />
      )}
      <div>
        {children}
      </div>
      <Footer />
    </Fragment>
  )
}

export default Layout
