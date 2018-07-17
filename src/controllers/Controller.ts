import * as express from 'express'

export abstract class Controller {
  public abstract basePath: string
  public abstract getRoutes(): express.Router;
}