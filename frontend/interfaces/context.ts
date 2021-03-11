export type ContextAppTypeActions = {}

export type ContextAppTypeState = {
  name?: string,
  API_URL?: string
}

export type ContextAppType = [ContextAppTypeState, ContextAppTypeActions]


export type ContextSessionTypeActions = {}

export type ContextSessionTypeState = {}

export type ContextSessionType = [ContextSessionTypeState, ContextSessionTypeActions]