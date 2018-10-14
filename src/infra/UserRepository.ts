import { getSession } from '../Database'
import { NotFoundError } from '../domain/errors/NotFoundError';
import { User } from '../domain/user/User';
import { UserRecords } from "../domain/user/UserRecords";
import Logger from '../Logger';

export class UserRepository implements UserRecords {

  public findOne(email: string): Promise<User> {
    Logger.debug(`find user with email: ${email}`)
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
        .run(`MATCH (n:USER {email: {emailParam}}) RETURN n`, { emailParam: email })
        .then(result => {
          if (result.records.length) {
            resolve(new User(email))
          } else {
            reject(new NotFoundError())
          }
          session.close()
        })
        .catch(err => {
          session.close()
          reject(err)
        })
    })
  }

  public create(user: User): Promise<User> {
    Logger.debug(`create user ${JSON.stringify(user)}`)
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
        .run(`MERGE (n:USER {email: {emailParam}})`, { emailParam: user.email })
        .then(() => {
          resolve(user)
          session.close()
        })
        .catch(err => {
          session.close()
          reject(err)
        })
    })
  }

}