import { NotFoundError } from "../errors/NotFoundError";
import { User } from "./User";
import { UserRecords } from "./UserRecords";

export class UserManagement {

  public constructor(
    public userRecords: UserRecords
  ) {
  }

  public findOne(email: string): Promise<User> {
    return this.userRecords.findOne(email)
  }

  public create(user: User): Promise<User> {
    return Promise.resolve(user)
    // return user
    //   .validate()
    //   .then(_ => this.userRecords.findOne(_.email))
    //   .then(foundUser => Promise.resolve(foundUser))
    //   .catch(err => {
    //     if (err instanceof NotFoundError) {
    //       return this.userRecords.create(user)
    //     } else {
    //       return Promise.reject(err)
    //     }
    //   })
    //   .catch(Promise.reject)
  }

}