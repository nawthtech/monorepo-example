import { Router } from 'itty-router'
import { cors } from './middleware/cors.js'
import { auth } from './middleware/auth.js'
import { healthHandlers } from './handlers/health.js'
import { authHandlers } from './handlers/auth.js'
import { userHandlers } from './handlers/users.js'
import { serviceHandlers } from './handlers/services.js'

const router = Router()

// ✅ وسائط عامة
router.all('*', cors)

// ✅ مسارات الصحة
router.get('/health', healthHandlers.check)
router.get('/health/live', healthHandlers.live)
router.get('/health/ready', healthHandlers.ready)

// ✅ المصادقة
router.post('/auth/register', authHandlers.register)
router.post('/auth/login', authHandlers.login)
router.post('/auth/refresh', authHandlers.refresh)
router.post('/auth/forgot-password', authHandlers.forgotPassword)

// ✅ المسارات المحمية
router.get('/user/profile', auth, userHandlers.getProfile)
router.put('/user/profile', auth, userHandlers.updateProfile)

// ✅ الخدمات
router.get('/services', serviceHandlers.getServices)
router.get('/services/:id', serviceHandlers.getServiceById)

// ✅ معالج للمسارات غير المعروفة
router.all('*', () => new Response('Not Found', { status: 404 }))

export default {
  async fetch(request, env, ctx) {
    try {
      // إضافة env إلى request للوصول للأسرار
      request.env = env
      return router.handle(request)
    } catch (error) {
      console.error('Worker Error:', error)
      return new Response(
        JSON.stringify({
          success: false,
          error: 'Internal Server Error',
          message: 'Something went wrong'
        }),
        {
          status: 500,
          headers: { 'Content-Type': 'application/json' }
        }
      )
    }
  }
}