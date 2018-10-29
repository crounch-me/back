import { validate } from 'class-validator'
import { Request, Response, Router } from 'express';
import { BAD_REQUEST } from 'http-status-codes'
import { Auth } from '../Auth';
import { User } from '../domain/user/User';
import { UserManagement } from '../domain/user/UserManagement';
import Logger from '../Logger';
import { Controller } from './Controller';

export class UserController extends Controller {

  public basePath: string = '/users'

  public constructor(
    public userManagement: UserManagement
  ) {
    super()
  }

  public getRoutes(): Router {
    const jwtCheck = Auth.getInstance().getJwtCheck()
    const router = Router()
    router.get('/:email', this.get.bind(this))
    router.post('/', jwtCheck, this.handleConnection.bind(this))
    return router
  }

  public get(req: Request, res: Response) {
    this.userManagement
      .findOne(req.params.email)
      .then(result => {
        res.json(result)
      })
      .catch(err => {
        res.json(err)
      })
  }

  public handleConnection(req: Request, res: Response) {
    this.userManagement.create(new User(req.user.email))
      .then(result => {
        res.json(result)
      })
      .catch(errors => {
        res.status(BAD_REQUEST).json(errors)
      })
  }

}