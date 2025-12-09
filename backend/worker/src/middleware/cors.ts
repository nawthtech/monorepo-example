import type { IRequest } from 'itty-router';
import type { Env } from '../types/database';

/**
 * CORS middleware
 */
export function handleCORS(request: IRequest, env: Env): Response | void {
  // Handle preflight requests
  if (request.method === 'OPTIONS') {
    return handlePreflight(request, env);
  }

  // Add CORS headers to response
  // The actual headers will be added in the response handler
}

/**
 * Handle CORS preflight requests
 */
function handlePreflight(request: IRequest, env: Env): Response {
  const origin = request.headers.get('Origin') || '';
  const allowedOrigin = getAllowedOrigin(origin, env);

  return new Response(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': allowedOrigin,
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-API-Key',
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Max-Age': '86400',
      'Vary': 'Origin',
    },
  });
}

/**
 * Get allowed origin for CORS
 */
export function getAllowedOrigin(requestOrigin: string, env: Env): string {
  const allowedOrigins = env.ALLOWED_ORIGINS 
    ? env.ALLOWED_ORIGINS.split(',')
    : ['https://nawthtech.com', 'https://www.nawthtech.com'];

  // Allow all origins in development
  if (env.ENVIRONMENT === 'development') {
    return requestOrigin || '*';
  }

  // Check if request origin is allowed
  if (allowedOrigins.includes(requestOrigin)) {
    return requestOrigin;
  }

  // Default to first allowed origin
  return allowedOrigins[0] || '*';
}

/**
 * Add CORS headers to response
 */
export function addCORSHeaders(response: Response, request: Request, env: Env): Response {
  const origin = request.headers.get('Origin') || '';
  const allowedOrigin = getAllowedOrigin(origin, env);

  const headers = new Headers(response.headers);
  headers.set('Access-Control-Allow-Origin', allowedOrigin);
  headers.set('Access-Control-Allow-Credentials', 'true');
  headers.set('Vary', 'Origin');

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}