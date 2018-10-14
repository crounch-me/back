import { Router } from "express"
import { HealthController } from "./controllers/HealthController";
import { UserController } from './controllers/UserController';
import { UserManagement } from './domain/user/UserManagement';
import { UserRepository } from './infra/UserRepository';


export function configureRouter(): Router {
  const router = Router()
  const healthController = new HealthController()
  router.use(healthController.basePath, healthController.getRoutes())

  const userRepository = new UserRepository()
  const userManagement = new UserManagement(userRepository)
  const userController = new UserController(userManagement)
  router.use(userController.basePath, userController.getRoutes())

  return router
}