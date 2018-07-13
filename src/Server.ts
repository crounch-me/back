import * as bodyParser from "body-parser"
import * as express from "express"
import * as Crawler from './Crawler'
import Logger from './Logger'
import { RegisterRoutes } from './routes'

class Server {
  constructor() {
    const app = express()
    app.use(bodyParser.urlencoded({extended: true}))
    app.use(bodyParser.json())

    app.use(Logger.logMiddleware)

    const port = process.env.PORT || 3000;
    Logger.debug('Application listening on port ' + port)

    RegisterRoutes(app)
    app.listen(port)
  }
}

// tslint:disable:no-unused-expression
new Server()
