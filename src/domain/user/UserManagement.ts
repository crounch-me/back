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
    return user
      .validate()
      .then(validatedUser => this.userRecords.create(validatedUser))
  }

}