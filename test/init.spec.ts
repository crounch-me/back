import { NextFunction, Request, Response } from 'express';
import * as sinon from 'sinon'
import { Auth } from '../src/Auth';
sinon.stub(Auth.getInstance(), "getJwtCheck").returns((req: Request, res: Response, next: NextFunction) => {
  req.user = { email: 'test@test.com' }
  next()
})
import { server } from '../src/Server'
import { emptyDatabase } from './util';

describe('Launch server', () => {
  it('should launch server', done => {
    server
      .launch()
      .then(() => emptyDatabase())
      .then(() => done())
  })
})