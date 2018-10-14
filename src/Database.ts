import * as dotenv from 'dotenv'
import neo4j from 'neo4j-driver'
import { Driver, Neo4jError, Session } from '../node_modules/neo4j-driver/types/v1';
import Logger from './Logger';

dotenv.config()

if (!process.env.DB_HOST || !process.env.DB_USER || !process.env.DB_PASSWORD) {
  throw new Error('Make sure you have DB_URL and DB_USER (and DB_PASSWORD) in your .env file');
}

let driver: Driver

export function launchDriver(): Promise<void> {
  return new Promise((resolve, reject) => {
    driver = neo4j.default.driver(`bolt://${process.env.DB_HOST}`, neo4j.default.auth.basic(`${process.env.DB_USER}`, `${process.env.DB_PASSWORD}`))
    driver.onError = (err: Neo4jError) => {
      Logger.error(`Error occured while connecting to database: ${err}`)
      reject(err)
    }

    driver.onCompleted = () => {
      Logger.debug('Driver connected')
      resolve()
    }
  })
}

export function getSession(): Session {
  return driver.session()
}