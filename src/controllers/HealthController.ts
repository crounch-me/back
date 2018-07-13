import { Controller, Get, Route } from "tsoa"

@Route("_health")
export class HealthController extends Controller {
  
  constructor() {
    super()
  }

  @Get("")
  public handleHealthCheck(): Promise<any> {
    this.setStatus(200)
    return Promise.resolve({status: 'ok'})
  }

  @Get("private")
  public handlePrivateHealtchCheck(): Promise<any> {
    this.setStatus(200)
    return Promise.resolve({status: 'ok'})
  }

}
