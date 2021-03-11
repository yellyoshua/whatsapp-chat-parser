import { createContext, FC, useState } from "react"
import { ContextAppType, ContextAppTypeState } from "@/interfaces/context"
import { APP_NAME } from "@/config/app";

type PropsContextApp = { children?: React.ReactChild, state?: ContextAppTypeState }

const initialStateContextApp: ContextAppTypeState = { name: APP_NAME }

export const AppContext = createContext<ContextAppType>([initialStateContextApp, {}]);

const AppContextProvider: FC<PropsContextApp> = ({ children, state }) => {

  return <AppContext.Provider value={[{ ...initialStateContextApp, ...state }, {}]}>
    {children}
  </AppContext.Provider>
}

export default AppContextProvider