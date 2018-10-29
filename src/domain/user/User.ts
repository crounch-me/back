import { IsEmail } from 'class-validator'

export class User {

  @IsEmail()
  public email: string

  public constructor(
    email:string
  ) {
    this.email = email
  }
}