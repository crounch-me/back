import { expect } from 'chai'
import { getOK, getPrivateOK } from '../../src/domain/health/HealthStatus';

describe('HealthStatus', () => {
  describe('getOK', () => {
    it('should return an ok status', done => {
      getOK()
        .then(result => {
          expect(result).to.deep.equals({ status: 'ok', version: process.env.npm_package_version })
          done()
        })
        .catch(done)
    })
  })

  describe('getPrivateOK', () => {
    it('should return an ok private status', done => {
      getPrivateOK()
        .then(result => {
          expect(result).to.deep.equals({ status: 'ok private', version: process.env.npm_package_version })
          done()
        })
        .catch(done)
    })
  })
})