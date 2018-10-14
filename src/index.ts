import Logger from './Logger';
import { server } from './Server'

server
  .launch()
  .catch(err => {
    Logger.error(`An error occured while launching server ${err}`)
  })