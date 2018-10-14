import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { server } from '../../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('ProductController', () => {
  describe('create', () => {
    it('should create a new product', done => {
      const product = { name: 'Carrot' }
      chai.request(server.app)
        .post('/products')
        .send({ product })
        .then(result => {
          expect(result.status).to.equals(200)
          expect(result.body).to.deep.equal(product)
          done()
        })
        .catch(done)
    })

    it('should return 400 if name is an empty string', done => {
      const product = { name: '' }
      chai.request(server.app)
        .post('/products')
        .send({ product })
        .then(result => {
          expect(result.status).to.equals(400)
          done()
        })
        .catch(done)
    })

  })
})