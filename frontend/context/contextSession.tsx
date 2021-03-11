import { createContext, FC } from "react"
import { ContextSessionType, ContextSessionTypeState } from "@/interfaces/context"

type PropsContextSession = { children: React.ReactChild, state: ContextSessionTypeState }

const initialStateContextSession = {}

const ContextSession = createContext<ContextSessionType>([initialStateContextSession, {}]);

const ContextSessionProvider: FC<PropsContextSession> = ({ children, state }) => {

  return <ContextSession.Provider value={[state, {}]}>
    {children}
  </ContextSession.Provider>
}

export default ContextSessionProvider