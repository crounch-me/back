import { Router } from 'express'
import * as jwt from 'express-jwt'

export abstract class Controller {
  public abstract basePath: string
  public abstract getRoutes(checkJwt: jwt.RequestHandler): Router;
}