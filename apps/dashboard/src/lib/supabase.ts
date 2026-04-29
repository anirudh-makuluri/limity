import { Pool } from 'pg'

const databaseUrl = process.env.DATABASE_URL || ''

let pool: Pool | null = null

export const getDbPool = () => {
  if (!pool) {
    pool = new Pool({
      connectionString: databaseUrl,
    })
  }
  return pool
}

export const createServerClient = () => {
  return getDbPool()
}
