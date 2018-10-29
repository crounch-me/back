import { Router } from "express"
import { HealthController } from "./controllers/HealthController";
import { UserController } from './controllers/UserController';
import { UserManagement } from './domain/user/UserManagement';
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

  Logger.debug('router configured')
  return router
}