import { getSession } from '../Database'
import { NotFoundError } from '../domain/errors/NotFoundError';
import { User } from '../domain/user/User';
import { UserRecords } from "../domain/user/UserRecords";
import Logger from '../Logger';

export class UserRepository implements UserRecords {

  public findOne(email: string): Promise<User> {
    Logger.debug('Hello there')
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
        .run(`MATCH (n:USER {email: {emailParam}}) RETURN n`, { emailParam: email })
        .then(result => {
          if (result.records.length) {
            resolve(new User(email))
          } else {
            reject(new NotFoundError)
          }
          session.close()
        })
        .catch(reject)
    })
  }

  public create(user: User): Promise<User> {
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
      .run(`CREATE (n:USER {email: {emailParam}}) RETURN n`, {emailParam: user.email})
      .then(result => {
        resolve(user)
        session.close()
      })
      .catch(reject)
    })
  }

}