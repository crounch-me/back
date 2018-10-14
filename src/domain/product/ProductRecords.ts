import { Product } from "./Product";

export interface ProductRecords {
  create(user: Product): Promise<Product>
}