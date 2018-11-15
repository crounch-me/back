import * as chai from 'chai'
import chaiHttp = require('chai-http')
import { server } from '../../../src/Server'

chai.use(chaiHttp)
const expect = chai.expect

describe('HealthController', () => {
  describe('handleHealthCheck', () => {
    it('should return ok status', done => {
      chai.request(server.app)
        .get('/_health')
        .then(res => {
          expect(res.status).to.equals(200)
          expect(res.body).to.deep.equals({ status: "ok", version: process.env.npm_package_version })
          done()
        })
        .catch(done)
    })

  })
}) 