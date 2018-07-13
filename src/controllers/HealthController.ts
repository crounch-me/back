import * as express from 'express'

export class HealthController {
  
  public handleHealthCheck(req: express.Request, res: express.Response) {
    res.json({status: 'ok'})
  }

  public handlePrivateHealtchCheck(req: express.Request, res: express.Response) {
    res.json({status: 'ok private'})
  }

}
