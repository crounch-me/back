import { validate } from 'class-validator';
import Logger from '../../Logger';
import { User } from './User';
import { UserRecords } from './UserRecords';

export class UserManagement {

  public constructor(
    public userRecords: UserRecords
  ) {
  }

  public findOne(email: string): Promise<User> {
    return this.userRecords.findOne(email)
  }

  public create(user: User): Promise<User> {
    Logger.debug(`validate user ${JSON.stringify(user)}`)

    return validate(user, { validationError: { target: false } })
      .then(err => {
        if (err.length > 0) {
          Logger.debug(`validation errors ${JSON.stringify(err)}`)
          return Promise.reject(err)
        }
        return this.userRecords.create(user)
      })
  }

}