export function cors(request) {
  const origin = request.headers.get('Origin')
  const allowedOrigins = request.env.CORS_ALLOWED_ORIGINS?.split(',') || []
  
  const isAllowedOrigin = allowedOrigins.includes(origin) || 
                         allowedOrigins.includes('*') || 
                         !origin

  if (request.method === 'OPTIONS') {
    return new Response(null, {
      status: 204,
      headers: {
        'Access-Control-Allow-Origin': isAllowedOrigin ? origin : allowedOrigins[0],
        'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With',
        'Access-Control-Allow-Credentials': 'true',
        'Access-Control-Max-Age': '86400'
      }
    })
  }

  // إضافة رؤوس CORS للاستجابات العادية
  const response = await next()
  response.headers.set('Access-Control-Allow-Origin', isAllowedOrigin ? origin : allowedOrigins[0])
  response.headers.set('Access-Control-Allow-Credentials', 'true')
  response.headers.set('Vary', 'Origin')
  
  return response
}