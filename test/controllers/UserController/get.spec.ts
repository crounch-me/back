import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { server } from '../../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('UserController', () => {
  describe('get', () => {
    it('should get a user', done => {
      chai.request(server.app)
        .get('/users/test@test.com')
        .then(result => {
          expect(result.status).to.equals(200)
          done()
        })
        .catch(done)
    })
  })
})