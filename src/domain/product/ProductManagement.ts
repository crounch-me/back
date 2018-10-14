import { Product } from "./Product";
import { ProductRecords } from './ProductRecords'

export class ProductManagement {

  public constructor(
    public productRecords: ProductRecords
  ) {
  }

  public create(product: Product): Promise<Product> {
    return product
      .validate()
      .then(validatedProduct => this.productRecords.create(validatedProduct))
  }
} 