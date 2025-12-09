import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { secureHeaders } from 'hono/secure-headers';
import { trimTrailingSlash } from 'hono/trailing-slash';

// Import handlers
import { apiRouter } from './api/routes';
import { handleEmail } from './email/worker';
import { handleScheduled } from './utils/cron';

// Types
export interface Env {
  // Database
  DB: D1Database;
  
  // Storage
  KV_CACHE: KVNamespace;
  KV_SESSIONS: KVNamespace;
  R2_STORAGE: R2Bucket;
  
  // Analytics
  ANALYTICS: AnalyticsEngineDataset;
  
  // Queues
  BACKGROUND_QUEUE: Queue;
  
  // Email
  EMAIL_ROUTING: EmailRouting;
  
  // Environment variables
  ENVIRONMENT: 'development' | 'staging' | 'production';
  JWT_SECRET: string;
  DOMAIN: string;
  APP_VERSION: string;
}

// App initialization
const app = new Hono<{ Bindings: Env }>();

// Global middleware
app.use('*', logger());
app.use('*', secureHeaders());
app.use('*', trimTrailingSlash());
app.use('*', cors({
  origin: ['https://nawthtech.com', 'https://www.nawthtech.com', 'http://localhost:3000'],
  allowHeaders: ['Content-Type', 'Authorization'],
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  exposeHeaders: ['Content-Length'],
  maxAge: 600,
  credentials: true,
}));

// Health check
app.get('/health', async (c) => {
  const startTime = Date.now();
  
  // Check database connection
  let dbStatus = 'healthy';
  try {
    await c.env.DB.prepare('SELECT 1').run();
  } catch (error) {
    dbStatus = 'unhealthy';
    console.error('Database health check failed:', error);
  }
  
  const health = {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime?.(),
    version: c.env.APP_VERSION || '1.0.0',
    environment: c.env.ENVIRONMENT || 'production',
    services: {
      database: dbStatus,
      cache: 'healthy',
      storage: 'healthy',
    },
    responseTime: Date.now() - startTime,
  };
  
  return c.json(health);
});

// API routes
app.route('/api', apiRouter);

// Static assets (if serving frontend)
app.get('/assets/*', async (c) => {
  const path = c.req.path.replace('/assets/', '');
  const object = await c.env.R2_STORAGE.get(path);
  
  if (!object) {
    return c.notFound();
  }
  
  const headers = new Headers();
  object.writeHttpMetadata(headers);
  headers.set('etag', object.httpEtag);
  
  return new Response(object.body, { headers });
});

// Frontend SPA fallback
app.get('*', async (c) => {
  // Serve index.html for SPA routing
  const object = await c.env.R2_STORAGE.get('index.html');
  if (object) {
    return new Response(object.body, {
      headers: { 'Content-Type': 'text/html; charset=utf-8' },
    });
  }
  
  return c.text('Welcome to NawthTech API', 200);
});

// Error handling
app.onError((err, c) => {
  console.error('Unhandled error:', err);
  
  // Log to Analytics Engine
  if (c.env.ANALYTICS) {
    c.env.ANALYTICS.writeDataPoint({
      indexes: ['error'],
      blobs: [err.message, c.req.method, c.req.path],
      doubles: [1],
    });
  }
  
  return c.json(
    {
      success: false,
      error: 'Internal server error',
      message: c.env.ENVIRONMENT === 'production' ? 'Something went wrong' : err.message,
    },
    500
  );
});

// Email worker handler
export const email = {
  async email(message: ForwardableEmailMessage, env: Env, ctx: ExecutionContext): Promise<void> {
    await handleEmail(message, env, ctx);
  },
};

// Scheduled tasks handler
export const scheduled = {
  async scheduled(event: ScheduledEvent, env: Env, ctx: ExecutionContext): Promise<void> {
    await handleScheduled(event, env, ctx);
  },
};

// Export the app
export default {
  fetch: app.fetch,
  email: email.email,
  scheduled: scheduled.scheduled,
};