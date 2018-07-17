import * as bodyParser from "body-parser"
import * as cors from 'cors'
import * as express from "express"
import Logger from './Logger'
import { configureRouter } from "./Router";

class Server {
  
  public app: express.Express = express()

  constructor() {
    const router = express.Router()
    this.app.use(bodyParser.urlencoded({ extended: true }))
    this.app.use(bodyParser.json())

    this.app.use(cors())

    this.app.use(configureRouter(router))
    this.app.use(Logger.logMiddleware)

    const port = process.env.PORT || 3000;
    Logger.debug('Application listening on port ' + port)

    this.app.listen(port)
  }
}

// tslint:disable:no-unused-expression
export const server = new Server()
