// worker/src/utils/database.js
import postgres from 'postgres'

export function createDatabaseConnection(env) {
  const connectionString = env.DATABASE_URL
  
  if (!connectionString) {
    throw new Error('DATABASE_URL is required')
  }

  return postgres(connectionString, {
    ssl: env.ENVIRONMENT === 'production' ? 'require' : 'allow',
    idle_timeout: 20,
    max_lifetime: 60 * 30
  })
}

export async function withDatabase(handler) {
  return async (request, env, ...args) => {
    const sql = createDatabaseConnection(env)
    
    try {
      // إضافة اتصال DB إلى request
      request.db = sql
      return await handler(request, env, ...args)
    } finally {
      await sql.end()
    }
  }
}