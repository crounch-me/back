import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { server } from '../../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('UserController', () => {
  describe('handleConnection', () => {
    it('should handleConnection and create a new user', done => {
      chai.request(server.app)
        .post('/users')
        .then(result => {
          expect(result.status).to.equals(200)
          done()
        })
        .catch(done)
    })
  })
})