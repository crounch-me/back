import { Request, Response, Router } from 'express'
import { Auth } from '../Auth';
import { getOK, getPrivateOK } from '../domain/health/HealthStatus';
import { Controller } from './Controller';

export class HealthController extends Controller {

  public basePath: string = '/_health'

  public getRoutes(): Router {
    const jwtCheck = Auth.getInstance().getJwtCheck()
    const router = Router()
    router.get('/', this.handleHealthCheck)
    router.get('/private', jwtCheck, this.handlePrivateHealtchCheck)
    return router
  }

  public handleHealthCheck(req: Request, res: Response) {
    getOK()
      .then(result => {
        res.json(result)
      })
  }

  public handlePrivateHealtchCheck(req: Request, res: Response) {
    getPrivateOK()
      .then(result => {
        res.json(result)
      })
  }

}
