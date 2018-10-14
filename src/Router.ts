import { Router } from "express"
import { HealthController } from "./controllers/HealthController";
import { ProductController } from "./controllers/ProductController";
import { UserController } from './controllers/UserController';
import { ProductManagement } from "./domain/product/ProductManagement";
import { UserManagement } from './domain/user/UserManagement';
import { ProductRepository } from "./infra/ProductRepository";
import { UserRepository } from './infra/UserRepository';
import Logger from "./Logger";


export function configureRouter(): Router {
  const router = Router()

  router.use('/', (req, res, next) => {
    Logger.debug(`${new Date()} ${req.method} ${req.url}`)
    next()
  })

  const healthController = new HealthController()
  router.use(healthController.basePath, healthController.getRoutes())

  const userRepository = new UserRepository()
  const userManagement = new UserManagement(userRepository)
  const userController = new UserController(userManagement)
  router.use(userController.basePath, userController.getRoutes())

  const productRepository = new ProductRepository()
  const productManagement = new ProductManagement(productRepository)
  const productController = new ProductController(productManagement)

  router.use(productController.basePath, productController.getRoutes())

  Logger.debug('router configured')

  return router
}