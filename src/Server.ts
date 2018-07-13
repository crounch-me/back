import * as bodyParser from "body-parser"
import * as express from "express"
import Logger from './Logger'
import { configureRouter } from "./Router";

class Server {
  constructor() {
    const app = express()
    const router = express.Router()
    app.use(bodyParser.urlencoded({extended: true}))
    app.use(bodyParser.json())
    
    app.use(configureRouter(router))
    app.use(Logger.logMiddleware)

    const port = process.env.PORT || 3000;
    Logger.debug('Application listening on port ' + port)

    app.listen(port)
  }
}

// tslint:disable:no-unused-expression
new Server()
