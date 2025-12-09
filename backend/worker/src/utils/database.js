// worker/src/utils/database.js
export class DatabaseManager {
  constructor(env, bindingName = 'NAWTHTECH_DB') {
    this.env = env
    this.bindingName = bindingName
    this.db = null
  }

  // إنشاء الاتصال بقاعدة D1
  connect() {
    if (!this.env?.D1) {
      throw new Error('D1 binding not found in environment')
    }

    this.db = this.env.D1(this.bindingName)
    return this.db
  }

  // تنفيذ استعلام SQL مع معاملات
  async query(sql, params = []) {
    if (!this.db) this.connect()
    try {
      const result = await this.db.prepare(sql).bind(...params).all()
      return result.results
    } catch (err) {
      console.error('D1 query error:', err)
      throw err
    }
  }

  // تنفيذ استعلام واحد وإرجاع صف واحد فقط
  async queryOne(sql, params = []) {
    const results = await this.query(sql, params)
    return results[0] || null
  }
}

// إنشاء مخبأ عالمي للـ DatabaseManager
let cachedDatabaseManager = null

export function getDatabaseManager(env) {
  if (cachedDatabaseManager) return cachedDatabaseManager
  cachedDatabaseManager = new DatabaseManager(env)
  return cachedDatabaseManager
}

// Middleware لتضمين قاعدة البيانات في الطلب
export function withDatabase(handler) {
  return async (request, env, ...args) => {
    try {
      const dbManager = getDatabaseManager(env)
      request.db = dbManager
      return await handler(request, env, ...args)
    } catch (err) {
      console.error('Database middleware error:', err)
      return new Response(
        JSON.stringify({
          success: false,
          error: 'DATABASE_CONNECTION_FAILED',
          message: err.message
        }),
        { status: 503, headers: { 'Content-Type': 'application/json' } }
      )
    }
  }
}