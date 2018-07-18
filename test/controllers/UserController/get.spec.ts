import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { User } from '../../../src/domain/user/User';
import { UserRepository } from '../../../src/infra/UserRepository';
import { server } from '../../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('UserController', () => {
  before(done => {
    const userRepository = new UserRepository()
    userRepository.create(new User("test@test.com")).then(() => done())
  })

  describe('get', () => {
    it('should get a user', done => {
      chai.request(server.app)
        .get('/users/test@test.com')
        .then(result => {
          expect(result.status).to.equals(200)
          expect(result.body).to.deep.equals({ email: "test@test.com" })
          done()
        })
        .catch(done)
    })
  })
})