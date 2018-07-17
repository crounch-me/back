import * as express from 'express'
import { Auth } from '../Auth';
import { getOK, getPrivateOK } from '../domain/health/HealthStatus';
import { Controller } from './Controller';

export class HealthController extends Controller {

  public basePath: string = '/_health'

  public getRoutes(): express.Router {
    const jwtCheck = Auth.getInstance().getJwtCheck()
    const router = express.Router()
    router.get('/', this.handleHealthCheck)
    router.get('/private', jwtCheck, this.handlePrivateHealtchCheck)
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
