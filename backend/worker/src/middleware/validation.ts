import type { IRequest } from 'itty-router';
import { z } from 'zod';
import { errorResponse } from '../utils/responses';

// Validation schemas
const schemas = {
  register: z.object({
    email: z.string().email('Invalid email address'),
    username: z.string()
      .min(3, 'Username must be at least 3 characters')
      .max(30, 'Username cannot exceed 30 characters')
      .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers and underscores'),
    password: z.string()
      .min(8, 'Password must be at least 8 characters')
      .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
      .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
      .regex(/[0-9]/, 'Password must contain at least one number')
      .regex(/[^A-Za-z0-9]/, 'Password must contain at least one special character'),
    full_name: z.string().optional(),
  }),

  login: z.object({
    email: z.string().email('Invalid email address'),
    password: z.string().min(1, 'Password is required'),
  }),

  updateProfile: z.object({
    full_name: z.string().min(1, 'Full name is required').optional(),
    avatar_url: z.string().url('Invalid URL').optional(),
    bio: z.string().max(500, 'Bio cannot exceed 500 characters').optional(),
    settings: z.record(z.any()).optional(),
  }),

  createService: z.object({
    name: z.string()
      .min(3, 'Name must be at least 3 characters')
      .max(100, 'Name cannot exceed 100 characters'),
    description: z.string().max(500, 'Description cannot exceed 500 characters').optional(),
    category: z.string().optional(),
    tags: z.array(z.string()).optional(),
    config: z.record(z.any()).optional(),
  }),

  updateService: z.object({
    name: z.string()
      .min(3, 'Name must be at least 3 characters')
      .max(100, 'Name cannot exceed 100 characters')
      .optional(),
    description: z.string().max(500, 'Description cannot exceed 500 characters').optional(),
    category: z.string().optional(),
    tags: z.array(z.string()).optional(),
    config: z.record(z.any()).optional(),
    status: z.enum(['pending', 'active', 'suspended', 'deleted']).optional(),
  }),

  generateAI: z.object({
    prompt: z.string()
      .min(1, 'Prompt is required')
      .max(5000, 'Prompt cannot exceed 5000 characters'),
    provider: z.enum(['gemini', 'openai', 'huggingface']).optional(),
    model: z.string().optional(),
    type: z.enum(['text', 'image']).optional(),
    options: z.record(z.any()).optional(),
  }),

  updateUser: z.object({
    email: z.string().email('Invalid email address').optional(),
    username: z.string()
      .min(3, 'Username must be at least 3 characters')
      .max(30, 'Username cannot exceed 30 characters')
      .optional(),
    role: z.enum(['user', 'admin', 'moderator']).optional(),
    email_verified: z.boolean().optional(),
    quota_text_tokens: z.number().int().min(0).optional(),
    quota_images: z.number().int().min(0).optional(),
    quota_videos: z.number().int().min(0).optional(),
    quota_audio_minutes: z.number().int().min(0).optional(),
  }),
};

/**
 * Validate request against schema
 */
export async function validateRequest(schemaName: keyof typeof schemas) {
  return async (request: IRequest): Promise<Response | void> => {
    try {
      const schema = schemas[schemaName];
      const body = await request.json?.();
      
      const result = schema.safeParse(body);
      
      if (!result.success) {
        const errors = result.error.errors.map(err => ({
          field: err.path.join('.'),
          message: err.message,
        }));
        
        return errorResponse('Validation failed', {
          errors,
          schema: schemaName,
        }, 400);
      }
      
      // Replace request body with validated data
      request.parsedBody = result.data;
    } catch (error) {
      return errorResponse('Invalid request body', 400);
    }
  };
}

/**
 * Validate query parameters
 */
export function validateQuery(schema: z.ZodSchema) {
  return async (request: IRequest): Promise<Response | void> => {
    try {
      const url = new URL(request.url);
      const queryParams = Object.fromEntries(url.searchParams);
      
      const result = schema.safeParse(queryParams);
      
      if (!result.success) {
        return errorResponse('Invalid query parameters', 400);
      }
      
      request.query = result.data;
    } catch (error) {
      return errorResponse('Invalid query parameters', 400);
    }
  };
}

/**
 * Validate path parameters
 */
export function validateParams(schema: z.ZodSchema) {
  return async (request: IRequest): Promise<Response | void> => {
    try {
      const result = schema.safeParse(request.params);
      
      if (!result.success) {
        return errorResponse('Invalid path parameters', 400);
      }
      
      request.params = result.data;
    } catch (error) {
      return errorResponse('Invalid path parameters', 400);
    }
  };
}

// Extend Request type
declare global {
  interface IRequest {
    parsedBody?: any;
  }
}