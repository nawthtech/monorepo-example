/**
 * Configuration Module Export
 */

export * from './settings'
export * from './constants'

// Re-export commonly used functions
export { 
  getSettings, 
  getSetting, 
  getApiEndpoint, 
  isFeatureEnabled,
  isDevelopment,
  isProduction,
  getEnvSetting 
} from './settings'

// Export default settings instance
export { default as settings } from './settings'