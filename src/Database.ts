import * as dotenv from 'dotenv'
import neo4j from 'neo4j-driver'
import { Driver, Neo4jError, Session } from '../node_modules/neo4j-driver/types/v1';
import Logger from './Logger';

dotenv.config()

if (!process.env.DB_HOST || !process.env.DB_USER || !process.env.DB_PASSWORD) {
  throw new Error('Make sure you have DB_URL and DB_USER (and DB_PASSWORD) in your .env file');
}

let driver: Driver

export function launchDriver() {
  driver = neo4j.default.driver(`bolt://${process.env.DB_HOST}`, neo4j.default.auth.basic(`${process.env.DB_USER}`, `${process.env.DB_PASSWORD}`))
  driver.onError = (error: Neo4jError) => {
    throw error
  }
  
  driver.onCompleted = () => {
    Logger.debug('Driver connected')
  }
}

export function getSession(): Session {
  return driver.session()
}