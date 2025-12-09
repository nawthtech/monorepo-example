import type { IRequest } from 'itty-router';
import type { Env } from '../types/database';

/**
 * Static file handler
 */
export async function handleStaticFile(
  request: IRequest,
  env: Env
): Promise<Response> {
  const url = new URL(request.url);
  let pathname = url.pathname;

  // Default to index.html for root path
  if (pathname === '/' || pathname === '') {
    pathname = '/index.html';
  }

  // Remove leading slash
  const filePath = pathname.startsWith('/') ? pathname.slice(1) : pathname;

  try {
    // Try to get from R2
    const file = await env.R2?.get(filePath);
    
    if (!file) {
      // Try with .html extension for SPA routes
      if (!filePath.includes('.')) {
        const htmlFile = await env.R2?.get('index.html');
        if (htmlFile) {
          return new Response(htmlFile.body, {
            headers: {
              'Content-Type': 'text/html; charset=utf-8',
              'Cache-Control': 'public, max-age=3600',
            },
          });
        }
      }
      
      return new Response('File not found', { status: 404 });
    }

    // Determine content type
    const contentType = getContentType(filePath);
    
    // Set caching headers
    const cacheControl = filePath.match(/\.(css|js|png|jpg|jpeg|gif|svg|ico|woff|woff2|ttf|eot)$/)
      ? 'public, max-age=31536000, immutable'
      : 'public, max-age=3600';

    return new Response(file.body, {
      headers: {
        'Content-Type': contentType,
        'Cache-Control': cacheControl,
        'ETag': `"${file.etag}"`,
      },
    });
  } catch (error) {
    console.error('Static file error:', error);
    return new Response('Internal server error', { status: 500 });
  }
}

/**
 * Get content type from file extension
 */
function getContentType(filename: string): string {
  const ext = filename.split('.').pop()?.toLowerCase() || '';
  
  const contentTypes: Record<string, string> = {
    // HTML
    html: 'text/html; charset=utf-8',
    htm: 'text/html; charset=utf-8',
    
    // CSS
    css: 'text/css; charset=utf-8',
    
    // JavaScript
    js: 'application/javascript; charset=utf-8',
    mjs: 'application/javascript; charset=utf-8',
    
    // JSON
    json: 'application/json; charset=utf-8',
    
    // Images
    png: 'image/png',
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    gif: 'image/gif',
    svg: 'image/svg+xml',
    ico: 'image/x-icon',
    webp: 'image/webp',
    
    // Fonts
    woff: 'font/woff',
    woff2: 'font/woff2',
    ttf: 'font/ttf',
    eot: 'application/vnd.ms-fontobject',
    
    // Text
    txt: 'text/plain; charset=utf-8',
    md: 'text/markdown; charset=utf-8',
    
    // Documents
    pdf: 'application/pdf',
    
    // Audio
    mp3: 'audio/mpeg',
    wav: 'audio/wav',
    ogg: 'audio/ogg',
    
    // Video
    mp4: 'video/mp4',
    webm: 'video/webm',
  };
  
  return contentTypes[ext] || 'application/octet-stream';
}