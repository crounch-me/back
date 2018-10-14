import { isLength } from 'validator'
const nameMinLength = 4
export class Product {
  public constructor(
    public name: string
  ) {
  }

  public validate(): Promise<Product> {
    if (isLength(this.name.trim(), { min: nameMinLength } )) { 
      return Promise.resolve(this)
    } else {
      return Promise.reject(`${this.name} should have length > ${nameMinLength}`)
    }
  }
}