import * as express from 'express'
import { getOK, getPrivateOK } from '../domain/health/HealthStatus';
import { checkJwt } from '../Router';
import { Controller } from './Controller';

export class HealthController extends Controller {

  public basePath: string = '/_health'

  public getRoutes(): express.Router {
    const router = express.Router()
    router.get('/', this.handleHealthCheck)
    router.get('/private', checkJwt, this.handlePrivateHealtchCheck)
    return router
  }

  public handleHealthCheck(req: express.Request, res: express.Response) {
    getOK()
      .then(result => {
        res.json(result)
      })
  }

  public handlePrivateHealtchCheck(req: express.Request, res: express.Response) {
    getPrivateOK()
      .then(result => {
        res.json(result)
      })
  }

}
