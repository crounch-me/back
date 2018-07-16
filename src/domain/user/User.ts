
export class User {
  public constructor(
    public email: string
  ) {
    
  }

  public validate(): Promise<User> {
    return Promise.resolve(this)
  }
}