import * as chai from 'chai'
import * as chaiAsPromised from 'chai-as-promised'
import * as sinon from 'sinon'
import { HealthController } from '../../../src/controllers/HealthController';

chai.should()
chai.use(chaiAsPromised)

describe('HealthController', () => {

  describe('Handle health check', () => {
    it('should return a simple ok response', () => {
      const controller: HealthController = new HealthController()
      return controller.handleHealthCheck().should.eventually.deep.equal({status: 'ok'})
    })
  })

})