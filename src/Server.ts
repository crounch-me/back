import * as bodyParser from "body-parser"
import * as cors from 'cors'
import * as express from "express"
import { launchDriver } from "./Database";
import Logger from './Logger'
import { configureRouter } from "./Router";

class Server {
  constructor() {
    const app = express()
    app.use(bodyParser.urlencoded({ extended: true }))
    app.use(bodyParser.json())

    app.use(cors())
    app.use(Logger.logMiddleware)

    app.use(configureRouter())

    const port = process.env.PORT || 3000;
    Logger.debug('Application listening on port ' + port)
    launchDriver()

    app.listen(port)
  }
}

// tslint:disable:no-unused-expression
new Server()
