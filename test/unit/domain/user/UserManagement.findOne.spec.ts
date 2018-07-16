import { expect } from 'chai'
import * as sinon from 'ts-sinon'
import { User } from '../../../../src/domain/user/User';
import { UserManagement } from "../../../../src/domain/user/UserManagement";
import { UserRepository } from "../../../../src/infra/UserRepository";

let userRepository: UserRepository
let userManagement: UserManagement
const sinonSandbox = sinon.default.sandbox.create()
const defaultEmail = "test@test.com"
const defaultUser = new User(defaultEmail)

describe('UserManagement.findOne', () => {

  before(() => {
    userRepository = new UserRepository()
    userManagement = new UserManagement(userRepository)
  })

  it('should call UserRecords findOne with email', done => {
    const mock = sinonSandbox
      .stub(userRepository, "findOne")
      .withArgs(defaultEmail)
      .returns(Promise.resolve(defaultUser))

    userManagement
      .findOne('test@test.com')
      .then(result => {
        expect(result).to.deep.equals(defaultUser)
        expect(mock.calledOnce).equals(true)
        done()
      })
      .catch(done)
  })

})