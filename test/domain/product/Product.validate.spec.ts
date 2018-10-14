import { expect } from 'chai'
import { Product } from '../../../src/domain/product/Product';

describe('User.validate', () => {
  it('should return user if all is correct', done => {
    const product = new Product('Carrot')
    product
      .validate()
      .then(result => {
        expect(result).to.deep.equals(product)
        done()
      })
      .catch(done)
  })

  it('should return error if length of name is lower than 4', done => {
    const user = new Product('a')
    user
      .validate()
      .catch(err => {
        expect(err).to.deep.equals('a should have length > 4')
        done()
      })
      .catch(done)
  })

})