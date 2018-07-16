import { getSession } from '../Database'
import { User } from '../domain/user/User';
import { UserRecords } from "../domain/user/UserRecords";

export class UserRepository implements UserRecords {

  public getOne(email: string): Promise<User> {
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
        .run(`MATCH (n:USER {email: {emailParam}}) RETURN n`, { emailParam: email })
        .then(result => {
          resolve(new User(email))
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