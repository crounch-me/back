import { expect } from 'chai'
import * as sinon from 'ts-sinon'
import { User } from "../../../../src/domain/user/User";
import { UserManagement } from "../../../../src/domain/user/UserManagement";
import { UserRepository } from "../../../../src/infra/UserRepository";

const userRepository: UserRepository = new UserRepository()
const userManagement: UserManagement = new UserManagement(userRepository)
const sinonSandbox = sinon.default.sandbox.create()
const defaultUser = new User("test@test.com")

describe('UserManagement.create', () => {
  it('should return user', done => {
    const createStub = sinonSandbox.stub(userRepository, "create").returns(Promise.resolve(defaultUser))
    userManagement
      .create(defaultUser)
      .then(result => {
        expect(result).to.deep.equals(defaultUser)
        createStub.restore()
        done()
      })
      .catch(done)
  })

  it('should call validate method', done => {
    const validateStub = sinonSandbox.stub(defaultUser, "validate").returns(Promise.resolve(defaultUser))
    const createStub = sinonSandbox.stub(userRepository, "create").returns(Promise.resolve(defaultUser))
    userManagement
      .create(defaultUser)
      .then(() => {
        expect(validateStub.calledOnce).to.equals(true)
        done()
      })
      .catch(done)
  })

  it('should call create on user if validated', done => {
    const validateStub = sinonSandbox.stub(defaultUser, "validate").returns(Promise.resolve(defaultUser))
    const createStub = sinonSandbox.stub(userRepository, "create").returns(Promise.resolve(defaultUser))
    userManagement
      .create(defaultUser)
      .then(() => {
        expect(createStub.calledOnce).to.equals(true)
        done()
      })
      .catch(done)
  })

  afterEach(() => {
    sinonSandbox.restore()
  })

})