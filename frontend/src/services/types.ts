/**
 * API Types
 */

export interface ApiResponse<T = any> {
  data: T;
  status: number;
  success: boolean;
  message?: string;
  meta?: {
    pagination?: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
      hasNext: boolean;
      hasPrev: boolean;
    };
    [key: string]: any;
  };
}

export interface ErrorResponse {
  message: string;
  status: number;
  errors?: Record<string, string[]>;
  timestamp: string;
  path?: string;
}

export interface RequestConfig {
  headers?: Record<string, string>;
  params?: Record<string, any>;
  timeout?: number;
  signal?: AbortSignal;
  formData?: boolean;
  onUploadProgress?: (progress: {
    loaded: number;
    total: number;
    percentage: number;
  }) => void;
}