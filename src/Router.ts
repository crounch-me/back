import * as dotenv from 'dotenv'
import * as express from "express"
import * as jwt from 'express-jwt'
import * as jwksRsa from 'jwks-rsa'
import { HealthController } from "./controllers/HealthController";

dotenv.config()

if (!process.env.AUTH0_DOMAIN || !process.env.AUTH0_AUDIENCE) {
  throw new Error('Make sure you have AUTH0_DOMAIN, and AUTH0_AUDIENCE in your .env file');
}

const checkJwt = jwt({
  // Validate the audience and the issuer.
  algorithms: ['RS256'],
  audience: process.env.AUTH0_AUDIENCE,
  issuer: `https://${process.env.AUTH0_DOMAIN}/`,

  // Dynamically provide a signing key based on the kid in the header and the singing keys provided by the JWKS endpoint.
  secret: jwksRsa.expressJwtSecret({
    cache: true,
    jwksRequestsPerMinute: 5,
    jwksUri: `https://${process.env.AUTH0_DOMAIN}/.well-known/jwks.json`,
    rateLimit: true
  })

});

export function configureRouter(router: express.Router): express.Router {
  const healthController =  new HealthController()
  
  router.get('/_health', healthController.handleHealthCheck)
  router.get('/_health/private', checkJwt, healthController.handlePrivateHealtchCheck)

  return router
}