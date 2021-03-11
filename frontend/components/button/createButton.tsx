import { FC } from "react"
import Link from 'next/link'


const CreateButton: FC = () => {
  return <li className="ml-5 px-3 py-1">
    <Link href="/messages">
      <a className="px-3 py-1 rounded font-semibold bg-gray-100 hover:bg-purple-600 duration-200 hover:text-white text-gray-800">
        Crear
    </a>
    </Link>
  </li>
}

export default CreateButton