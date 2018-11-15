import { expect } from 'chai'
import * as sinon from 'sinon'
import { User } from "../../../src/domain/user/User";
import { UserManagement } from "../../../src/domain/user/UserManagement";
import { UserRepository } from "../../../src/infra/UserRepository";

const userRepository: UserRepository = new UserRepository()
const userManagement: UserManagement = new UserManagement(userRepository)
const sinonSandbox = sinon.createSandbox()
const defaultUser = new User("test@test.com")

describe('UserManagement.create', () => {
  it('should create and return user', done => {
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

  afterEach(() => {
    sinonSandbox.restore()
  })

  it('should return error if email not valid', done => {
    const wrongUser = new User("a")

    userManagement
      .create(wrongUser)
      .then(errors => {
        done()
      })
      .catch(errors => {
        expect(errors[0].value).equals("a")
        expect(errors[0].property).equals("email")
        done()
      })
  })

}) 