import { getSession } from '../src/Database'

export function emptyDatabase(): Promise<void> {
  const session = getSession()
  return session
    .run('MATCH (n) DETACH DELETE n')
    .then(() => {
      session.close()
      Promise.resolve()
    })
    .catch(err => {
      session.close()
      Promise.reject(err)
    })
}