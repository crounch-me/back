import { Request, Response, Router } from 'express';
import { Auth } from '../Auth';
import { Product } from '../domain/product/Product';
import { ProductManagement } from '../domain/product/ProductManagement';
import { Controller } from './Controller';

export class ProductController extends Controller {

  public basePath: string = '/products'

  public constructor(
    public productManagement: ProductManagement
  ) {
    super()
  }

  public getRoutes(): Router {
    const jwtCheck = Auth.getInstance().getJwtCheck()
    const router = Router()
    router.post('/', jwtCheck, this.create.bind(this))
    return router
  }

  public create(req: Request, res: Response) {
    const product = new Product(req.body.product.name)
    this.productManagement
      .create(product)
      .then(result => res.json(result))
      .catch(err => res.status(400).json(err))
  }

}