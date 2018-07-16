import { User } from "./User";

export interface UserRecords {
  getOne(email: string): Promise<User>
  create(user: User): Promise<User>
}