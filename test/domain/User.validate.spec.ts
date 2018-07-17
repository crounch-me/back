import { expect } from 'chai'
import { User } from "../../src/domain/user/User";

describe('User.validate', () => {
  it('should return user if all is correct', done => {
    const user = new User('test@test.com')
    user
      .validate()
      .then(result => {
        expect(result).to.deep.equals(user)
        done()
      })
      .catch(done)
  })

  it('should return error if user email is not an email', done => {
    const user = new User('a')
    user
      .validate()
      .catch(err => {
        expect(err).to.deep.equals('a is not an email')
        done()
      })
      .catch(done)
  })

})