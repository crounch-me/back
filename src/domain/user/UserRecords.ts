import { User } from "./User";

export interface UserRecords {
  findOne(email: string): Promise<User>
  create(user: User): Promise<User>
}