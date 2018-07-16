import { User } from "./User";
import { UserRecords } from "./UserRecords";


export class UserManagement {

  public constructor(
    public userRecords: UserRecords
  ) {

  }

  public getOne(email: string): Promise<User> {
    return this.userRecords.getOne(email)
  }

  public create(user: User): Promise<User> {
    return user
      .validate()
      .then(this.userRecords.create)
  }

}