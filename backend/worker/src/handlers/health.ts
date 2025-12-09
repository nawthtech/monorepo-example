import type { IRequest } from 'itty-router';
import type { Env } from '../types/database';
import { successResponse } from '../utils/responses';

/**
 * Health check handler
 */
export async function handleHealthCheck(
  request: IRequest,
  env: Env,
  ctx: ExecutionContext
): Promise<Response> {
  const startTime = Date.now();
  
  try {
    // Check database connection
    const dbCheck = await env.DB.prepare('SELECT 1 as status').first();
    
    // Check KV
    const kvCheck = await env.KV?.put('health-check', startTime.toString(), {
      expirationTtl: 60,
    });
    
    // Check R2 (if configured)
    let r2Check = false;
    if (env.R2) {
      try {
        await env.R2.list({ limit: 1 });
        r2Check = true;
      } catch (error) {
        r2Check = false;
      }
    }

    const health = {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime?.(),
      environment: env.ENVIRONMENT || 'production',
      version: '1.0.0',
      checks: {
        database: !!dbCheck,
        kv: !!env.KV,
        r2: r2Check,
      },
      response_time: Date.now() - startTime,
    };

    return new Response(JSON.stringify(successResponse(health), null, 2), {
      status: 200,
      headers: {
        'Content-Type': 'application/json',
        'Cache-Control': 'no-cache',
      },
    });
  } catch (error) {
    const errorHealth = {
      status: 'unhealthy',
      timestamp: new Date().toISOString(),
      environment: env.ENVIRONMENT || 'production',
      error: error instanceof Error ? error.message : 'Unknown error',
      response_time: Date.now() - startTime,
    };

    return new Response(JSON.stringify(errorHealth, null, 2), {
      status: 503,
      headers: {
        'Content-Type': 'application/json',
        'Cache-Control': 'no-cache',
      },
    });
  }
}