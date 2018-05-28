import { Request, Response } from 'express'
import { Logger as WinstonLogger, LoggerInstance, transports } from 'winston'

export default class Logger {
  
  public static debug(text: string) {
    Logger.winstonLogger.debug(text);
  }

  public static info(text: string) {
    Logger.winstonLogger.info(text);
  }

  public static warn(text: string) {
    Logger.winstonLogger.warn(text);
  }

  public static error(text: string) {
    Logger.winstonLogger.error(text);
  }

  public static usage(wrongLevel: string) {
    Logger.warn("Unrecognized level " + wrongLevel + ". Please pick one from : debug, error, warn and info");
  }

  public static write(text: string) {
    Logger.winstonLogger.info(text);
  }

  public static log(level: string, text: string, ...args: any[]) {
    Logger.winstonLogger.log(level, text, ...args);
  }

  public static logMiddleware(req: Request, res: Response, callback: any) {
    Logger.debug(new Date().toISOString() + ' ' + req.method + ' ' + req.path)
    callback()
  }

  private static winstonLogger: LoggerInstance = new WinstonLogger({
    exitOnError: false,
    level: "debug",
    transports: [
      new (transports.Console)({
        colorize: true,
        handleExceptions: true,
        json: false,
        level: "debug",
        name: "console logs"
      })
    ]
  });
}