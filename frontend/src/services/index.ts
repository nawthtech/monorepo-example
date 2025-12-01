/**
 * Services Module - Main Export File
 */

// Export API
export { api, apiHelpers } from './api';
export type { 
  ApiResponse, 
  ErrorResponse, 
  RequestConfig, 
  PaginationParams,
  UploadProgressEvent,
  UploadProgressCallback 
} from './api';

// Export Admin API
export { adminAPI, adminHelpers } from './admin';
export type * from './admin';

// Re-export everything for convenience
export * from './api';