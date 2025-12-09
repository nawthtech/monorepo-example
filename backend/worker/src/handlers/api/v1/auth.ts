import type { IRequest } from 'itty-router';
import type { Env, User, RegisterRequest, LoginRequest, UpdateProfileRequest } from '../../../types/database';
import { successResponse, errorResponse, paginatedResponse } from '../../../utils/responses';
import { hashPassword, verifyPassword, generateToken } from '../../../utils/crypto';

/**
 * Register new user
 */
export async function registerUser(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const body = request.parsedBody as RegisterRequest;
    
    // Check if user already exists
    const existingUser = await env.DB.prepare(
      'SELECT id FROM users WHERE email = ? OR username = ?'
    )
      .bind(body.email, body.username)
      .first();

    if (existingUser) {
      return errorResponse('User already exists', 409);
    }

    // Hash password
    const passwordHash = await hashPassword(body.password);
    const userId = crypto.randomUUID();

    // Create user
    await env.DB.prepare(`
      INSERT INTO users (id, email, username, password_hash, full_name)
      VALUES (?, ?, ?, ?, ?)
    `)
      .bind(
        userId,
        body.email,
        body.username,
        passwordHash,
        body.full_name || null
      )
      .run();

    // Generate token
    const token = generateToken({ sub: userId }, env.JWT_SECRET);

    const user = {
      id: userId,
      email: body.email,
      username: body.username,
      full_name: body.full_name,
      role: 'user' as const,
      email_verified: false,
    };

    return new Response(
      JSON.stringify(
        successResponse({
          user,
          token,
        }, 'User registered successfully'),
        null,
        2
      ),
      {
        status: 201,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('Registration error:', error);
    return errorResponse('Failed to register user', 500);
  }
}

/**
 * Login user
 */
export async function loginUser(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const body = request.parsedBody as LoginRequest;
    
    // Get user by email
    const user = await env.DB.prepare(
      'SELECT * FROM users WHERE email = ? AND deleted_at IS NULL'
    )
      .bind(body.email)
      .first<User>();

    if (!user) {
      return errorResponse('Invalid credentials', 401);
    }

    // Verify password
    const isValid = await verifyPassword(body.password, user.password_hash);
    if (!isValid) {
      return errorResponse('Invalid credentials', 401);
    }

    // Update last login
    await env.DB.prepare(
      'UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?'
    ).bind(user.id).run();

    // Generate token
    const token = generateToken({ sub: user.id }, env.JWT_SECRET);

    const userResponse = {
      id: user.id,
      email: user.email,
      username: user.username,
      role: user.role,
      email_verified: user.email_verified,
      full_name: user.full_name,
      avatar_url: user.avatar_url,
      bio: user.bio,
      created_at: user.created_at,
    };

    return new Response(
      JSON.stringify(
        successResponse({
          user: userResponse,
          token,
        }, 'Login successful'),
        null,
        2
      ),
      {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('Login error:', error);
    return errorResponse('Failed to login', 500);
  }
}

/**
 * Get current user
 */
export async function getCurrentUser(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const user = request.user!;
    
    const userResponse = {
      id: user.id,
      email: user.email,
      username: user.username,
      role: user.role,
      email_verified: user.email_verified,
      full_name: user.full_name,
      avatar_url: user.avatar_url,
      bio: user.bio,
      quota: {
        text_tokens: user.quota_text_tokens,
        images: user.quota_images,
        videos: user.quota_videos,
        audio_minutes: user.quota_audio_minutes,
      },
      settings: JSON.parse(user.settings),
      created_at: user.created_at,
      updated_at: user.updated_at,
      last_login_at: user.last_login_at,
    };

    return new Response(
      JSON.stringify(successResponse(userResponse), null, 2),
      {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('Get current user error:', error);
    return errorResponse('Failed to get user', 500);
  }
}

/**
 * Update user profile
 */
export async function updateUser(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const user = request.user!;
    const body = request.parsedBody as UpdateProfileRequest;
    
    // Build update query dynamically
    const updates: string[] = [];
    const values: any[] = [];
    
    if (body.full_name !== undefined) {
      updates.push('full_name = ?');
      values.push(body.full_name || null);
    }
    
    if (body.avatar_url !== undefined) {
      updates.push('avatar_url = ?');
      values.push(body.avatar_url || null);
    }
    
    if (body.bio !== undefined) {
      updates.push('bio = ?');
      values.push(body.bio || null);
    }
    
    if (body.settings !== undefined) {
      updates.push('settings = ?');
      values.push(JSON.stringify(body.settings));
    }
    
    if (updates.length === 0) {
      return errorResponse('No fields to update', 400);
    }
    
    updates.push('updated_at = CURRENT_TIMESTAMP');
    values.push(user.id);
    
    const query = `UPDATE users SET ${updates.join(', ')} WHERE id = ?`;
    
    await env.DB.prepare(query)
      .bind(...values)
      .run();

    // Get updated user
    const updatedUser = await env.DB.prepare(
      'SELECT * FROM users WHERE id = ?'
    )
      .bind(user.id)
      .first<User>();

    const userResponse = {
      id: updatedUser!.id,
      email: updatedUser!.email,
      username: updatedUser!.username,
      role: updatedUser!.role,
      email_verified: updatedUser!.email_verified,
      full_name: updatedUser!.full_name,
      avatar_url: updatedUser!.avatar_url,
      bio: updatedUser!.bio,
      settings: JSON.parse(updatedUser!.settings),
      updated_at: updatedUser!.updated_at,
    };

    return new Response(
      JSON.stringify(
        successResponse(userResponse, 'Profile updated successfully'),
        null,
        2
      ),
      {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('Update profile error:', error);
    return errorResponse('Failed to update profile', 500);
  }
}

/**
 * List users (admin only)
 */
export async function listUsers(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const url = new URL(request.url);
    const page = parseInt(url.searchParams.get('page') || '1');
    const limit = parseInt(url.searchParams.get('limit') || '20');
    const offset = (page - 1) * limit;
    
    // Get total count
    const countResult = await env.DB.prepare(
      'SELECT COUNT(*) as total FROM users WHERE deleted_at IS NULL'
    ).first<{ total: number }>();
    
    // Get users
    const users = await env.DB.prepare(`
      SELECT id, email, username, role, email_verified, full_name, 
             avatar_url, created_at, updated_at, last_login_at
      FROM users 
      WHERE deleted_at IS NULL
      ORDER BY created_at DESC
      LIMIT ? OFFSET ?
    `)
      .bind(limit, offset)
      .all<User>();

    const usersResponse = users.results.map(user => ({
      id: user.id,
      email: user.email,
      username: user.username,
      role: user.role,
      email_verified: user.email_verified,
      full_name: user.full_name,
      avatar_url: user.avatar_url,
      created_at: user.created_at,
      updated_at: user.updated_at,
      last_login_at: user.last_login_at,
    }));

    return new Response(
      JSON.stringify(
        paginatedResponse(usersResponse, countResult!.total, page, limit),
        null,
        2
      ),
      {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('List users error:', error);
    return errorResponse('Failed to list users', 500);
  }
}

/**
 * Get user by ID
 */
export async function getUserById(
  request: IRequest,
  env: Env
): Promise<Response> {
  try {
    const userId = request.params!.id;
    const currentUser = request.user!;
    
    // Check permissions (admin or self)
    if (currentUser.role !== 'admin' && currentUser.id !== userId) {
      return errorResponse('Access denied', 403);
    }
    
    const user = await env.DB.prepare(`
      SELECT id, email, username, role, email_verified, full_name, 
             avatar_url, bio, created_at, updated_at, last_login_at
      FROM users 
      WHERE id = ? AND deleted_at IS NULL
    `)
      .bind(userId)
      .first<User>();

    if (!user) {
      return errorResponse('User not found', 404);
    }

    const userResponse = {
      id: user.id,
      email: user.email,
      username: user.username,
      role: user.role,
      email_verified: user.email_verified,
      full_name: user.full_name,
      avatar_url: user.avatar_url,
      bio: user.bio,
      created_at: user.created_at,
      updated_at: user.updated_at,
      last_login_at: user.last_login_at,
    };

    return new Response(
      JSON.stringify(successResponse(userResponse), null, 2),
      {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }
    );
  } catch (error) {
    console.error('Get user error:', error);
    return errorResponse('Failed to get user', 500);
  }
}