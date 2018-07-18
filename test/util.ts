import { getSession } from '../src/Database'

export function emptyDatabase(): Promise<void> {
  const session = getSession()
  return session
    .run('MATCH (n) DETACH DELETE n')
    .then(() => Promise.resolve())
    .catch(err => Promise.reject(err))
}