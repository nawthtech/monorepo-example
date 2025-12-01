/**
 * NawthTech Application Settings Configuration
 * Compatible with React + Vite + Go Backend Monorepo
 */

// ==================== TYPES ====================
export interface AppConfig {
  name: string;
  version: string;
  description: string;
  baseUrl: string;
  environment: 'development' | 'test' | 'production';
  supportEmail: string;
  contactPhone: string;
}

export interface ApiConfig {
  baseURL: string;
  timeout: number;
  retryAttempts: number;
  endpoints: {
    auth: AuthEndpoints;
    users: UserEndpoints;
    media: MediaEndpoints;
    website: WebsiteEndpoints;
    ai: AIEndpoints;
    storage: StorageEndpoints;
  };
}

export interface AuthEndpoints {
  login: string;
  register: string;
  logout: string;
  verify: string;
  refresh: string;
  forgotPassword: string;
  resetPassword: string;
}

export interface UserEndpoints {
  profile: string;
  update: string;
  changePassword: string;
  list: string;
  roles: string;
}

export interface MediaEndpoints {
  upload: string;
  list: string;
  delete: string;
  folders: string;
}

export interface WebsiteEndpoints {
  settings: string;
  sections: string;
  themes: string;
}

export interface AIEndpoints {
  generate: string;
  analyze: string;
  chat: string;
}

export interface StorageEndpoints {
  upload: string;
  download: string;
  list: string;
}

export interface StorageConfig {
  baseURL: string;
  uploadEndpoint: string;
  maxFileSize: number;
  allowedTypes: string[];
  chunkSize: number;
  timeout: number;
}

export interface FeatureFlags {
  aiAssistant: boolean;
  socialMediaIntegration: boolean;
  analytics: boolean;
  multiLanguage: boolean;
  darkMode: boolean;
  offlineMode: boolean;
  pushNotifications: boolean;
  voiceCommands: boolean;
  fileStorage: boolean;
  apiAccess: boolean;
}

export interface LocalizationConfig {
  defaultLanguage: string;
  supportedLanguages: Array<{
    code: string;
    name: string;
    dir: 'rtl' | 'ltr';
    flag: string;
  }>;
  fallbackLanguage: string;
  dateFormat: string;
  timeFormat: string;
  currency: string;
  timezone: string;
}

export interface ThemeConfig {
  defaultTheme: string;
  availableThemes: string[];
  colorSchemes: string[];
  direction: 'rtl' | 'ltr';
  fontScale: number;
}

export interface UploadConfig {
  maxFileSize: number;
  allowedImageTypes: string[];
  allowedDocumentTypes: string[];
  allowedVideoTypes: string[];
  maxFilesPerUpload: number;
  chunkSize: number;
  timeout: number;
}

export interface SecurityConfig {
  password: {
    minLength: number;
    requireUppercase: boolean;
    requireLowercase: boolean;
    requireNumbers: boolean;
    requireSpecialChars: boolean;
  };
  session: {
    timeout: number;
    refreshInterval: number;
  };
  csrf: {
    enabled: boolean;
    headerName: string;
  };
  cors: {
    allowedOrigins: string[];
    allowedMethods: string[];
  };
}

export interface AnalyticsConfig {
  googleAnalytics: string;
  facebookPixel: string;
  hotjar: string;
  mixpanel: string;
  enabled: boolean;
}

export interface SocialConfig {
  platforms: {
    instagram: SocialPlatform;
    twitter: SocialPlatform;
    facebook: SocialPlatform;
    linkedin: SocialPlatform;
  };
  share: {
    enabled: boolean;
    platforms: string[];
  };
}

export interface SocialPlatform {
  enabled: boolean;
  appId?: string;
  apiKey?: string;
  apiSecret?: string;
  clientId?: string;
  redirectUri?: string;
}

export interface AIConfig {
  openai: {
    apiKey: string;
    model: string;
    maxTokens: number;
    temperature: number;
  };
  stabilityai: {
    apiKey: string;
    engine: string;
  };
  huggingface: {
    apiKey: string;
  };
  rateLimit: {
    requestsPerMinute: number;
    requestsPerHour: number;
  };
}

export interface PerformanceConfig {
  lazyLoading: boolean;
  imageOptimization: boolean;
  codeSplitting: boolean;
  caching: {
    enabled: boolean;
    duration: number;
  };
  debounce: {
    search: number;
    resize: number;
  };
}

export interface StorageConfigDetail {
  localStorage: {
    prefix: string;
    version: string;
  };
  sessionStorage: {
    prefix: string;
  };
  indexedDB: {
    name: string;
    version: number;
  };
}

export interface NotificationConfig {
  desktop: {
    enabled: boolean;
    permission: 'default' | 'granted' | 'denied';
  };
  push: {
    enabled: boolean;
    publicKey: string;
  };
  inApp: {
    enabled: boolean;
    duration: number;
    position: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
  };
  sounds: {
    enabled: boolean;
    volume: number;
  };
}

export interface SEOConfig {
  defaultTitle: string;
  defaultDescription: string;
  defaultKeywords: string;
  defaultImage: string;
  twitterHandle: string;
  facebookAppId: string;
  canonicalUrl: string;
}

export interface PaymentConfig {
  currency: string;
  taxRate: number;
  stripe: {
    publishableKey: string;
    secretKey: string;
  };
  plans: {
    basic: Plan;
    pro: Plan;
    enterprise: Plan;
  };
}

export interface Plan {
  price: number;
  features: string[];
}

export interface DevelopmentConfig {
  debug: boolean;
  logLevel: 'debug' | 'info' | 'warn' | 'error';
  reduxDevTools: boolean;
  mockData: boolean;
}

export interface Settings {
  app: AppConfig;
  api: ApiConfig;
  storage: StorageConfig;
  features: FeatureFlags;
  localization: LocalizationConfig;
  theme: ThemeConfig;
  upload: UploadConfig;
  security: SecurityConfig;
  analytics: AnalyticsConfig;
  social: SocialConfig;
  ai: AIConfig;
  performance: PerformanceConfig;
  storageConfig: StorageConfigDetail;
  notifications: NotificationConfig;
  seo: SEOConfig;
  payments: PaymentConfig;
  development: DevelopmentConfig;
}

// ==================== DEFAULT SETTINGS ====================
export const defaultSettings: Settings = {
  app: {
    name: 'NawthTech',
    version: '2.1.0',
    description: 'منصة الذكاء الاصطناعي المتكاملة لتطوير الأعمال الرقمية',
    baseUrl: import.meta.env.VITE_APP_BASE_URL || 'https://nawthtech.com',
    environment: (import.meta.env.MODE as 'development' | 'test' | 'production') || 'development',
    supportEmail: 'support@nawthtech.com',
    contactPhone: '+966 12 345 6789'
  },

  api: {
    baseURL: import.meta.env.VITE_API_URL || 'https://api.nawthtech.com/api/v1',
    timeout: 30000,
    retryAttempts: 3,
    endpoints: {
      auth: {
        login: '/auth/login',
        register: '/auth/register',
        logout: '/auth/logout',
        verify: '/auth/verify',
        refresh: '/auth/refresh',
        forgotPassword: '/auth/forgot-password',
        resetPassword: '/auth/reset-password'
      },
      users: {
        profile: '/users/profile',
        update: '/users/update',
        changePassword: '/users/change-password',
        list: '/users',
        roles: '/users/roles'
      },
      media: {
        upload: '/uploads/upload',
        list: '/uploads',
        delete: '/uploads/delete',
        folders: '/uploads/folders'
      },
      website: {
        settings: '/website/settings',
        sections: '/website/sections',
        themes: '/website/themes'
      },
      ai: {
        generate: '/ai/generate',
        analyze: '/ai/analyze',
        chat: '/ai/chat'
      },
      storage: {
        upload: '/storage/upload',
        download: '/storage/download',
        list: '/storage/files'
      }
    }
  },

  storage: {
    baseURL: import.meta.env.VITE_STORAGE_URL || 'https://storage.nawthtech.com',
    uploadEndpoint: '/api/v1/uploads/upload',
    maxFileSize: 100 * 1024 * 1024, // 100MB
    allowedTypes: [
      'image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif',
      'application/pdf', 'video/mp4', 'audio/mpeg', 'application/zip'
    ],
    chunkSize: 5 * 1024 * 1024, // 5MB chunks
    timeout: 60000
  },

  features: {
    aiAssistant: true,
    socialMediaIntegration: true,
    analytics: true,
    multiLanguage: true,
    darkMode: true,
    offlineMode: false,
    pushNotifications: true,
    voiceCommands: false,
    fileStorage: true,
    apiAccess: true
  },

  localization: {
    defaultLanguage: 'ar',
    supportedLanguages: [
      { code: 'ar', name: 'العربية', dir: 'rtl', flag: '🇸🇦' },
      { code: 'en', name: 'English', dir: 'ltr', flag: '🇺🇸' }
    ],
    fallbackLanguage: 'ar',
    dateFormat: 'DD/MM/YYYY',
    timeFormat: 'HH:mm',
    currency: 'SAR',
    timezone: 'Asia/Riyadh'
  },

  theme: {
    defaultTheme: 'light',
    availableThemes: ['light', 'dark', 'auto'],
    colorSchemes: ['blue', 'green', 'purple', 'orange'],
    direction: 'rtl',
    fontScale: 1.0
  },

  upload: {
    maxFileSize: 100 * 1024 * 1024, // 100MB
    allowedImageTypes: ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'],
    allowedDocumentTypes: ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'],
    allowedVideoTypes: ['video/mp4', 'video/mpeg', 'video/quicktime'],
    maxFilesPerUpload: 10,
    chunkSize: 5 * 1024 * 1024, // 5MB
    timeout: 60000
  },

  security: {
    password: {
      minLength: 6,
      requireUppercase: false,
      requireLowercase: false,
      requireNumbers: true,
      requireSpecialChars: false
    },
    session: {
      timeout: 24 * 60 * 60 * 1000, // 24 hours
      refreshInterval: 60 * 60 * 1000 // 1 hour
    },
    csrf: {
      enabled: true,
      headerName: 'X-CSRF-Token'
    },
    cors: {
      allowedOrigins: ['https://nawthtech.com', 'https://www.nawthtech.com'],
      allowedMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS']
    }
  },

  analytics: {
    googleAnalytics: import.meta.env.VITE_GA_TRACKING_ID || '',
    facebookPixel: import.meta.env.VITE_FB_PIXEL_ID || '',
    hotjar: import.meta.env.VITE_HOTJAR_ID || '',
    mixpanel: import.meta.env.VITE_MIXPANEL_TOKEN || '',
    enabled: import.meta.env.MODE === 'production'
  },

  social: {
    platforms: {
      instagram: {
        enabled: true,
        appId: import.meta.env.VITE_INSTAGRAM_APP_ID || '',
        redirectUri: import.meta.env.VITE_INSTAGRAM_REDIRECT_URI || 'https://nawthtech.com/auth/instagram/callback'
      },
      twitter: {
        enabled: true,
        apiKey: import.meta.env.VITE_TWITTER_API_KEY || '',
        apiSecret: import.meta.env.VITE_TWITTER_API_SECRET || ''
      },
      facebook: {
        enabled: false,
        appId: import.meta.env.VITE_FACEBOOK_APP_ID || ''
      },
      linkedin: {
        enabled: true,
        clientId: import.meta.env.VITE_LINKEDIN_CLIENT_ID || ''
      }
    },
    share: {
      enabled: true,
      platforms: ['twitter', 'linkedin', 'whatsapp', 'telegram']
    }
  },

  ai: {
    openai: {
      apiKey: import.meta.env.VITE_OPENAI_API_KEY || '',
      model: 'gpt-4',
      maxTokens: 2000,
      temperature: 0.7
    },
    stabilityai: {
      apiKey: import.meta.env.VITE_STABILITYAI_API_KEY || '',
      engine: 'stable-diffusion-xl-1024-v1-0'
    },
    huggingface: {
      apiKey: import.meta.env.VITE_HUGGINGFACE_API_KEY || ''
    },
    rateLimit: {
      requestsPerMinute: 60,
      requestsPerHour: 1000
    }
  },

  performance: {
    lazyLoading: true,
    imageOptimization: true,
    codeSplitting: true,
    caching: {
      enabled: true,
      duration: 5 * 60 * 1000 // 5 minutes
    },
    debounce: {
      search: 300,
      resize: 150
    }
  },

  storageConfig: {
    localStorage: {
      prefix: 'nawthtech_',
      version: 'v1'
    },
    sessionStorage: {
      prefix: 'nawthtech_session_'
    },
    indexedDB: {
      name: 'NawthTechDB',
      version: 1
    }
  },

  notifications: {
    desktop: {
      enabled: true,
      permission: 'default'
    },
    push: {
      enabled: true,
      publicKey: import.meta.env.VITE_VAPID_PUBLIC_KEY || ''
    },
    inApp: {
      enabled: true,
      duration: 5000,
      position: 'top-right'
    },
    sounds: {
      enabled: true,
      volume: 0.5
    }
  },

  seo: {
    defaultTitle: 'NawthTech - منصة الذكاء الاصطناعي المتكاملة',
    defaultDescription: 'منصة متكاملة لإدارة وتطوير أعمال وسائل التواصل الاجتماعي باستخدام الذكاء الاصطناعي',
    defaultKeywords: 'ذكاء اصطناعي, وسائل تواصل, متابعين, مشاهدات, تحليل, نمو, NawthTech, أعمال رقمية',
    defaultImage: '/assets/og-image.jpg',
    twitterHandle: '@nawthtech',
    facebookAppId: import.meta.env.VITE_FACEBOOK_APP_ID || '',
    canonicalUrl: 'https://nawthtech.com'
  },

  payments: {
    currency: 'SAR',
    taxRate: 0.15, // 15% VAT
    stripe: {
      publishableKey: import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY || '',
      secretKey: import.meta.env.VITE_STRIPE_SECRET_KEY || ''
    },
    plans: {
      basic: {
        price: 99,
        features: ['ai_assistant', 'social_analytics']
      },
      pro: {
        price: 199,
        features: ['ai_assistant', 'social_analytics', 'premium_support']
      },
      enterprise: {
        price: 499,
        features: ['all_features', 'custom_integration']
      }
    }
  },

  development: {
    debug: import.meta.env.MODE === 'development',
    logLevel: import.meta.env.MODE === 'development' ? 'debug' : 'error',
    reduxDevTools: import.meta.env.MODE === 'development',
    mockData: import.meta.env.VITE_USE_MOCK_DATA === 'true'
  }
};

// ==================== ENVIRONMENT OVERRIDES ====================
const environmentOverrides = {
  development: {
    app: {
      baseUrl: 'http://localhost:5173',
      environment: 'development' as const
    },
    api: {
      baseURL: 'http://localhost:8080/api/v1'
    },
    storage: {
      baseURL: 'http://localhost:8080'
    },
    features: {
      mockData: true
    },
    development: {
      debug: true,
      logLevel: 'debug' as const,
      reduxDevTools: true,
      mockData: true
    }
  },
  test: {
    app: {
      environment: 'test' as const
    },
    api: {
      baseURL: 'https://test-api.nawthtech.com/api/v1'
    },
    storage: {
      baseURL: 'https://test-storage.nawthtech.com'
    },
    development: {
      debug: true,
      logLevel: 'info' as const
    }
  },
  production: {
    app: {
      environment: 'production' as const,
      baseUrl: 'https://nawthtech.com'
    },
    api: {
      baseURL: 'https://api.nawthtech.com/api/v1'
    },
    storage: {
      baseURL: 'https://storage.nawthtech.com'
    },
    features: {
      mockData: false
    },
    development: {
      debug: false,
      logLevel: 'error' as const,
      reduxDevTools: false,
      mockData: false
    }
  }
};

// ==================== UTILITY FUNCTIONS ====================
/**
 * Get merged settings based on current environment
 */
export function getSettings(): Settings {
  const env = import.meta.env.MODE as keyof typeof environmentOverrides;
  const overrides = environmentOverrides[env] || environmentOverrides.development;
  
  return deepMerge(defaultSettings, overrides);
}

/**
 * Get a specific setting by dot notation path
 * @example getSetting('app.name') // returns 'NawthTech'
 */
export function getSetting<T = any>(path: string): T | undefined {
  const keys = path.split('.');
  let value: any = getSettings();
  
  for (const key of keys) {
    if (value && typeof value === 'object' && key in value) {
      value = value[key];
    } else {
      return undefined;
    }
  }
  
  return value as T;
}

/**
 * Check if a feature is enabled
 */
export function isFeatureEnabled(featureName: keyof FeatureFlags): boolean {
  return getSettings().features[featureName];
}

/**
 * Get API endpoint by path
 */
export function getApiEndpoint(
  category: keyof ApiConfig['endpoints'],
  endpoint: string
): string {
  const settings = getSettings();
  const baseURL = settings.api.baseURL.replace(/\/$/, '');
  const endpointPath = (settings.api.endpoints[category] as any)[endpoint];
  
  if (!endpointPath) {
    throw new Error(`Endpoint not found: ${category}.${endpoint}`);
  }
  
  return `${baseURL}${endpointPath}`;
}

/**
 * Get environment-specific setting
 */
export function getEnvSetting(): 'development' | 'test' | 'production' {
  return getSettings().app.environment;
}

/**
 * Check if running in development mode
 */
export const isDevelopment = (): boolean => getEnvSetting() === 'development';

/**
 * Check if running in production mode
 */
export const isProduction = (): boolean => getEnvSetting() === 'production';

// ==================== HELPER FUNCTIONS ====================
function deepMerge(target: any, source: any): any {
  const output = { ...target };
  
  if (isObject(target) && isObject(source)) {
    Object.keys(source).forEach(key => {
      if (isObject(source[key])) {
        if (!(key in target)) {
          Object.assign(output, { [key]: source[key] });
        } else {
          output[key] = deepMerge(target[key], source[key]);
        }
      } else {
        Object.assign(output, { [key]: source[key] });
      }
    });
  }
  
  return output;
}

function isObject(item: any): boolean {
  return item && typeof item === 'object' && !Array.isArray(item);
}

// ==================== EXPORT SETTINGS INSTANCE ====================
export const settings = getSettings();

export default settings;