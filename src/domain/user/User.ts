import { isEmail }  from 'validator'

export class User {
  public constructor(
    public email: string
  ) {
  }

  public validate(): Promise<User> {
    if (isEmail(this.email)) {
      return Promise.resolve(this)
    } else {
      return Promise.reject(`${this.email} is not an email`)
    }
  }
}