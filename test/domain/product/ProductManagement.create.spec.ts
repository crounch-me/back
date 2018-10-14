import { expect } from 'chai'
import * as sinon from 'sinon'
import { Product } from '../../../src/domain/product/Product';
import { ProductManagement } from '../../../src/domain/product/ProductManagement';
import { ProductRepository } from '../../../src/infra/ProductRepository';

const productRepository: ProductRepository = new ProductRepository()
const productManagement: ProductManagement = new ProductManagement(productRepository)
const sinonSandbox = sinon.createSandbox()
const defaultProduct = new Product('Carrot')

describe('ProductManagement.create', () => {
  it('should return user', done => {
    const createStub = sinonSandbox.stub(productRepository, "create").returns(Promise.resolve(defaultProduct))
    productManagement
      .create(defaultProduct)
      .then(result => {
        expect(result).to.deep.equals(defaultProduct)
        createStub.restore()
        done()
      })
      .catch(done)
  })
 
  it('should call validate method', done => {
    const validateStub = sinonSandbox.stub(defaultProduct, "validate").returns(Promise.resolve(defaultProduct))
    const createStub = sinonSandbox.stub(productRepository, "create").returns(Promise.resolve(defaultProduct))
    productManagement
      .create(defaultProduct)
      .then(() => {
        expect(validateStub.calledOnce).to.equals(true)
        done()
      })
      .catch(done)
  })

  it('should call create on user if validated', done => {
    const validateStub = sinonSandbox.stub(defaultProduct, "validate").returns(Promise.resolve(defaultProduct))
    const createStub = sinonSandbox.stub(productRepository, "create").returns(Promise.resolve(defaultProduct))
    productManagement
      .create(defaultProduct)
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