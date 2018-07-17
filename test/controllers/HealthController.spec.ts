import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { server } from '../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('HealthController', () => {
  it('should return ok status', done => {
    chai.request(server.app)
      .get('/_health')
      .then(res => {
        expect(res.status).to.equals(200)
        expect(res.body).to.deep.equals({status: "ok"})
        done()
      })
      .catch(done)
  })

  it('should return ok private status', done => {
    chai.request(server.app)
      .get('/_health/private')
      .then(res => {
        expect(res.status).to.equals(200)
        expect(res.body).to.deep.equals({status: "ok private"})
        done()
      })
      .catch(done)
  })
})
