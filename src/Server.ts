import * as bodyParser from "body-parser"
import * as cors from 'cors'
import * as express from "express"
import { launchDriver } from "./Database";
import Logger from './Logger'
import { configureRouter } from "./Router";

class Server {
  public static getInstance(): Server {
    return this.instance || (this.instance = new this())
  }

  private static instance: Server

  public app: express.Express

  private constructor() {
    this.app = express()

    this.app.use(bodyParser.urlencoded({ extended: true }))
    this.app.use(bodyParser.json())

    this.app.use(cors())

    this.app.use(configureRouter())
    this.app.use(Logger.logMiddleware)
  }

  public launch(): Promise<void> {
    return new Promise((resolve, reject) => {
      launchDriver()
        .then(() => {
          const port = process.env.PORT || 3000;
          this.app.listen(port, () => {
            Logger.debug('Application listening on port ' + port)
            resolve()
          })
        })
        .catch(reject)
    })
  }
}

export const server = Server.getInstance()
