import neo4j from 'neo4j-driver'
import { Neo4jError, Session } from '../node_modules/neo4j-driver/types/v1';
import Logger from './Logger';

const dbUrl = 'bolt://localhost'
const dbUser = 'neo4j'
const dbPassword = 'test'

const driver = neo4j.default.driver(dbUrl, neo4j.default.auth.basic(dbUser, dbPassword))
driver.onError = (error: Neo4jError) => {
  throw error
}

driver.onCompleted = () => {
  Logger.debug('Driver connected')
}

export function getSession(): Session {
  return driver.session()
}