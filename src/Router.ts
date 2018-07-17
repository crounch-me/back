import * as dotenv from 'dotenv'
import * as express from "express"

import { HealthController } from "./controllers/HealthController";

dotenv.config()

if (!process.env.AUTH0_DOMAIN || !process.env.AUTH0_AUDIENCE) {
  throw new Error('Make sure you have AUTH0_DOMAIN, and AUTH0_AUDIENCE in your .env file');
}

export function configureRouter(router: express.Router): express.Router {
  const healthController = new HealthController()
  router.use(healthController.basePath, healthController.getRoutes())

  return router
}