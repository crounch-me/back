import { NextFunction, Request, Response } from 'express';
import * as sinon from 'sinon'
import { Auth } from '../src/Auth';

sinon.stub(Auth.getInstance(), "getJwtCheck").returns((req: Request, res: Response, next: NextFunction) => next())
