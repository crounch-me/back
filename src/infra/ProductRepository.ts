import { getSession } from '../Database'
import { Product } from '../domain/product/Product';
import { ProductRecords } from '../domain/product/ProductRecords';
import Logger from '../Logger';

export class ProductRepository implements ProductRecords {

  public create(product: Product): Promise<Product> {
    Logger.debug(`create product ${JSON.stringify(product)}`)
    return new Promise((resolve, reject) => {
      const session = getSession()
      return session
        .run(`MERGE (n:PRODUCT {name: {nameParam}})`, { nameParam: product.name })
        .then(() => {
          resolve(product)
          session.close()
        })
        .catch(err => {
          session.close()
          reject(err)
        })
    })
  }

}