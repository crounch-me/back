import * as bodyParser from "body-parser"
import * as express from "express"
import { HealthController } from './controllers/HealthController'
import { RegisterRoutes } from './routes'

class Server {
  constructor() {
    const app = express()
    app.use(bodyParser.urlencoded({extended: true}))
    app.use(bodyParser.json())

    RegisterRoutes(app)
    app.listen(3000)
  }
}

// tslint:disable:no-unused-expression
new Server()
