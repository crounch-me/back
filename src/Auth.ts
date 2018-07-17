import * as jwt from 'express-jwt'
import * as jwksRsa from 'jwks-rsa'

export class Auth {
  public static getInstance(): Auth {
    return this.authInstance
  }
  
  private static authInstance = new Auth()

  public getJwtCheck(): jwt.RequestHandler {
    return jwt({
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
  }
}
