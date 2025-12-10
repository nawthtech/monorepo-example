/**
 * Admin API Service for NawthTech Dashboard
 * Integrates with Go backend and Cloudflare services
 */

import { api } from './api';
import type { PaginationParams, ApiResponse } from './api';

// ==================== TYPES ====================
export interface DashboardStats {
  // Core Metrics
  totalUsers: number;
  totalOrders: number;
  totalRevenue: number;
  activeServices: number;
  pendingOrders: number;
  supportTickets: number;
  conversionRate: number;
  bounceRate: number;
  storeVisits: number;
  newCustomers: number;
  growthRate: number;
  averageOrderValue: number;
  customerSatisfaction: number;
  
  // AI Metrics
  totalAIRequests: number;
  successfulAIRequests: number;
  aiProcessingTime: number;
  videoGenerationRequests: number;
  imageGenerationRequests: number;
  textGenerationRequests: number;
  
  // Platform Metrics
  cloudflareRequests: number;
  cloudflareBandwidth: number;
  cloudinaryUploads: number;
  slackNotifications: number;
  sentryErrors: number;
  
  // Business Metrics
  monthlyRecurringRevenue: number;
  churnRate: number;
  lifetimeValue: number;
  customerAcquisitionCost: number;
  returnOnInvestment: number;
}

export interface ServiceMetrics {
  serviceId: string;
  name: string;
  category: string;
  platform: string;
  
  // Usage Metrics
  totalOrders: number;
  completedOrders: number;
  pendingOrders: number;
  revenueGenerated: number;
  averageRating: number;
  
  // AI Metrics
  aiProcessingTime: number;
  successRate: number;
  errorRate: number;
  
  // Inventory
  stockLevel: number;
  lowStock: boolean;
  restockNeeded: boolean;
  
  // Performance
  conversionRate: number;
  customerSatisfaction: number;
  returnRate: number;
  
  // Platform Specific
  platformMetrics?: {
    instagram?: {
      followers: number;
      engagement: number;
      reach: number;
    };
    tiktok?: {
      views: number;
      likes: number;
      shares: number;
    };
    twitter?: {
      tweets: number;
      retweets: number;
      mentions: number;
    };
    youtube?: {
      views: number;
      subscribers: number;
      watchTime: number;
    };
  };
}

export interface AIDashboardMetrics {
  // General AI Stats
  totalRequests: number;
  totalSuccessfulRequests: number;
  totalFailedRequests: number;
  averageProcessingTime: number;
  totalCost: number;
  
  // Model-Specific Stats
  gemini: {
    requests: number;
    successRate: number;
    averageTime: number;
    cost: number;
  };
  openai: {
    requests: number;
    successRate: number;
    averageTime: number;
    cost: number;
  };
  ollama: {
    requests: number;
    successRate: number;
    averageTime: number;
    cost: number;
  };
  huggingface: {
    requests: number;
    successRate: number;
    averageTime: number;
    cost: number;
  };
  
  // Stability AI Stats
  stability: {
    imageRequests: number;
    videoRequests: number;
    successRate: number;
    averageTime: number;
    cost: number;
  };
  
  // Video Generation Stats
  videoGeneration: {
    totalRequests: number;
    lumaRequests: number;
    runwayRequests: number;
    pikaRequests: number;
    geminiVeoRequests: number;
    successRate: number;
    averageDuration: number;
    cost: number;
  };
  
  // Usage by Category
  usageByCategory: {
    contentGeneration: number;
    imageProcessing: number;
    videoGeneration: number;
    analytics: number;
    validation: number;
  };
}

export interface Order {
  id: string;
  orderNumber: string;
  user: {
    id: string;
    name: string;
    email: string;
    avatar?: string;
  };
  service: {
    id: string;
    name: string;
    category: string;
    platform: string;
    aiPowered: boolean;
  };
  amount: number;
  status: 'pending' | 'processing' | 'completed' | 'cancelled' | 'refunded';
  date: string;
  type: string;
  category: string;
  paymentMethod: string;
  paymentStatus: 'paid' | 'pending' | 'failed' | 'refunded';
  notes?: string;
  attachments?: string[];
  
  // AI Processing Details
  aiProcessing?: {
    model?: string;
    processingTime?: number;
    cost?: number;
    qualityScore?: number;
    validationStatus?: 'pending' | 'approved' | 'rejected';
    validationNotes?: string;
  };
  
  // Video Generation Specific
  videoGeneration?: {
    platform: 'luma' | 'runway' | 'pika' | 'gemini-veo' | 'stability';
    duration: number;
    resolution: string;
    format: string;
    outputUrl?: string;
    previewUrl?: string;
  };
}

export interface UserActivity {
  id: string;
  user: {
    id: string;
    name: string;
    email: string;
    avatar?: string;
    role: string;
  };
  action: string;
  service?: {
    id: string;
    name: string;
    category: string;
  };
  time: string;
  ip: string;
  type: 'login' | 'logout' | 'purchase' | 'view' | 'update' | 'delete' | 'create';
  details?: Record<string, any>;
  userAgent?: string;
  location?: {
    city?: string;
    country?: string;
    coordinates?: {
      lat: number;
      lng: number;
    };
  };
}

export interface SystemAlert {
  id: string;
  type: 'error' | 'warning' | 'info' | 'success';
  title: string;
  message: string;
  timestamp: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  resolved: boolean;
  actionRequired: boolean;
  service?: string;
  metadata?: Record<string, any>;
  
  // Service Specific
  source?: 'cloudflare' | 'cloudinary' | 'sentry' | 'slack' | 'ai' | 'payment' | 'database';
  assignedTo?: string;
  priority?: number;
}

export interface UserReport {
  id: string;
  name: string;
  email: string;
  role: string;
  status: 'active' | 'inactive' | 'suspended' | 'pending';
  joinedDate: string;
  lastLogin: string;
  totalOrders: number;
  totalSpent: number;
  subscriptionPlan?: string;
  tags?: string[];
  
  // AI Usage
  aiUsage?: {
    totalRequests: number;
    videoGenerationRequests: number;
    imageGenerationRequests: number;
    textGenerationRequests: number;
    totalCost: number;
  };
  
  // Platform Preferences
  preferredPlatforms?: string[];
  
  // Billing
  billingCycle?: 'monthly' | 'quarterly' | 'yearly';
  nextBillingDate?: string;
}

export interface RevenueReport {
  period: string;
  revenue: number;
  orders: number;
  averageOrderValue: number;
  growth: number;
  expenses: number;
  profit: number;
  margin: number;
  
  // AI Revenue Breakdown
  aiRevenue: number;
  serviceRevenue: number;
  subscriptionRevenue: number;
  
  // Platform Revenue
  platformRevenue?: {
    instagram?: number;
    tiktok?: number;
    twitter?: number;
    youtube?: number;
    facebook?: number;
  };
}

export interface DashboardData {
  stats: DashboardStats;
  serviceMetrics: ServiceMetrics[];
  recentOrders: Order[];
  userActivity: UserActivity[];
  systemAlerts: SystemAlert[];
  revenueTrend: RevenueReport[];
  topUsers: UserReport[];
  
  // AI Dashboard
  aiMetrics: AIDashboardMetrics;
  
  // Performance Metrics
  performanceMetrics: {
    responseTime: number;
    uptime: number;
    errorRate: number;
    serverLoad: number;
    
    // Cloudflare Metrics
    cloudflare: {
      requests: number;
      bandwidth: number;
      cacheHitRate: number;
      securityEvents: number;
    };
    
    // Sentry Metrics
    sentry: {
      totalErrors: number;
      unresolvedErrors: number;
      usersAffected: number;
    };
  };
}

export interface AnalyticsFilters {
  timeRange: 'today' | 'yesterday' | 'week' | 'month' | 'quarter' | 'year' | 'custom';
  startDate?: string;
  endDate?: string;
  category?: string;
  service?: string;
  userGroup?: string;
  status?: string;
  platform?: string;
  aiModel?: string;
}

export interface ExportOptions {
  format: 'csv' | 'excel' | 'pdf' | 'json';
  includeCharts: boolean;
  includeDetails: boolean;
  timeZone: string;
  language: string;
}

export interface AdminSettings {
  // General Settings
  siteMaintenance: boolean;
  registrationEnabled: boolean;
  emailNotifications: boolean;
  autoBackup: boolean;
  backupFrequency: 'daily' | 'weekly' | 'monthly';
  
  // Security Settings
  maxLoginAttempts: number;
  sessionTimeout: number;
  apiRateLimit: number;
  cacheEnabled: boolean;
  cacheDuration: number;
  securityLevel: 'low' | 'medium' | 'high' | 'strict';
  
  // AI Settings
  aiEnabled: boolean;
  defaultAIModel: 'gemini' | 'openai' | 'ollama' | 'huggingface';
  aiRateLimit: number;
  videoGenerationEnabled: boolean;
  imageGenerationEnabled: boolean;
  
  // Cloudflare Settings
  cloudflare: {
    cacheEnabled: boolean;
    workersEnabled: boolean;
    d1Enabled: boolean;
    kvEnabled: boolean;
  };
  
  // Cloudinary Settings
  cloudinary: {
    uploadEnabled: boolean;
    autoOptimize: boolean;
    transformationEnabled: boolean;
  };
  
  // Service Settings
  serviceLimits: {
    maxOrdersPerDay: number;
    maxAiRequestsPerDay: number;
    maxVideoDuration: number;
    maxFileSize: number;
  };
}

// ==================== ADMIN API SERVICE ====================
export const adminAPI = {
  // ==================== DASHBOARD ====================
  getDashboardData: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<DashboardData>> => {
    return api.get<DashboardData>(
      { category: 'admin', endpoint: 'dashboard' },
      {
        params: filters,
      }
    );
  },

  getDashboardStats: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<DashboardStats>> => {
    return api.get<DashboardStats>(
      { category: 'admin', endpoint: 'dashboard/stats' },
      {
        params: filters,
      }
    );
  },

  getAIDashboardMetrics: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<AIDashboardMetrics>> => {
    return api.get<AIDashboardMetrics>(
      { category: 'admin', endpoint: 'dashboard/ai-metrics' },
      {
        params: filters,
      }
    );
  },

  // ==================== ORDERS ====================
  getRecentOrders: async (
    limit: number = 10,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<Order[]>> => {
    return api.get<Order[]>(
      { category: 'admin', endpoint: 'orders/recent' },
      {
        params: {
          limit,
          ...filters,
        },
      }
    );
  },

  getAllOrders: async (
    params: PaginationParams & {
      status?: string;
      dateFrom?: string;
      dateTo?: string;
      userId?: string;
      serviceId?: string;
      platform?: string;
      aiModel?: string;
    } = {}
  ): Promise<ApiResponse<Order[]>> => {
    return api.getPaginated<Order>(
      { category: 'admin', endpoint: 'orders' },
      params
    );
  },

  getOrderDetails: async (orderId: string): Promise<ApiResponse<Order>> => {
    return api.get<Order>(
      { category: 'admin', endpoint: `orders/${orderId}` }
    );
  },

  updateOrderStatus: async (
    orderId: string,
    status: Order['status'],
    notes?: string
  ): Promise<ApiResponse<Order>> => {
    return api.put<Order>(
      { category: 'admin', endpoint: `orders/${orderId}/status` },
      { status, notes }
    );
  },

  refundOrder: async (
    orderId: string,
    amount?: number,
    reason?: string
  ): Promise<ApiResponse<Order>> => {
    return api.post<Order>(
      { category: 'admin', endpoint: `orders/${orderId}/refund` },
      { amount, reason }
    );
  },

  // ==================== AI MANAGEMENT ====================
  getAIRequests: async (
    params: PaginationParams & {
      model?: string;
      type?: string;
      status?: string;
      userId?: string;
      dateFrom?: string;
      dateTo?: string;
    } = {}
  ): Promise<ApiResponse<any>> => {
    return api.getPaginated(
      { category: 'admin', endpoint: 'ai/requests' },
      params
    );
  },

  getVideoGenerationRequests: async (
    params: PaginationParams & {
      platform?: string;
      status?: string;
      userId?: string;
      dateFrom?: string;
      dateTo?: string;
    } = {}
  ): Promise<ApiResponse<any>> => {
    return api.getPaginated(
      { category: 'admin', endpoint: 'ai/video-requests' },
      params
    );
  },

  retryAIRequest: async (requestId: string): Promise<ApiResponse<any>> => {
    return api.post<any>(
      { category: 'admin', endpoint: `ai/requests/${requestId}/retry` }
    );
  },

  cancelAIRequest: async (requestId: string): Promise<ApiResponse<any>> => {
    return api.delete<any>(
      { category: 'admin', endpoint: `ai/requests/${requestId}` }
    );
  },

  // ==================== USERS ====================
  getUserActivity: async (
    limit: number = 10,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<UserActivity[]>> => {
    return api.get<UserActivity[]>(
      { category: 'admin', endpoint: 'users/activity' },
      {
        params: {
          limit,
          ...filters,
        },
      }
    );
  },

  getAllUsers: async (
    params: PaginationParams & {
      status?: string;
      role?: string;
      dateFrom?: string;
      dateTo?: string;
      platform?: string;
    } = {}
  ): Promise<ApiResponse<UserReport[]>> => {
    return api.getPaginated<UserReport>(
      { category: 'admin', endpoint: 'users' },
      params
    );
  },

  getUserDetails: async (userId: string): Promise<ApiResponse<UserReport>> => {
    return api.get<UserReport>(
      { category: 'admin', endpoint: `users/${userId}` }
    );
  },

  updateUserStatus: async (
    userId: string,
    status: UserReport['status'],
    reason?: string
  ): Promise<ApiResponse<UserReport>> => {
    return api.put<UserReport>(
      { category: 'admin', endpoint: `users/${userId}/status` },
      { status, reason }
    );
  },

  updateUserRole: async (
    userId: string,
    role: string,
    permissions?: string[]
  ): Promise<ApiResponse<UserReport>> => {
    return api.put<UserReport>(
      { category: 'admin', endpoint: `users/${userId}/role` },
      { role, permissions }
    );
  },

  impersonateUser: async (userId: string): Promise<ApiResponse<{ token: string }>> => {
    return api.post<{ token: string }>(
      { category: 'admin', endpoint: `users/${userId}/impersonate` }
    );
  },

  // ==================== ANALYTICS & REPORTS ====================
  getRevenueReport: async (
    filters: AnalyticsFilters
  ): Promise<ApiResponse<RevenueReport[]>> => {
    return api.get<RevenueReport[]>(
      { category: 'admin', endpoint: 'analytics/revenue' },
      {
        params: filters,
      }
    );
  },

  getAIAnalytics: async (
    filters: AnalyticsFilters
  ): Promise<ApiResponse<any>> => {
    return api.get<any>(
      { category: 'admin', endpoint: 'analytics/ai' },
      {
        params: filters,
      }
    );
  },

  getPlatformAnalytics: async (
    platform: string,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<any>> => {
    return api.get<any>(
      { category: 'admin', endpoint: `analytics/platform/${platform}` },
      {
        params: filters,
      }
    );
  },

  // ==================== SYSTEM & ALERTS ====================
  getSystemAlerts: async (
    params: PaginationParams & {
      severity?: SystemAlert['severity'];
      resolved?: boolean;
      type?: SystemAlert['type'];
      source?: SystemAlert['source'];
    } = {}
  ): Promise<ApiResponse<SystemAlert[]>> => {
    return api.getPaginated<SystemAlert>(
      { category: 'admin', endpoint: 'system/alerts' },
      params
    );
  },

  resolveAlert: async (alertId: string, notes?: string): Promise<ApiResponse<SystemAlert>> => {
    return api.put<SystemAlert>(
      { category: 'admin', endpoint: `system/alerts/${alertId}/resolve` },
      { notes }
    );
  },

  acknowledgeAlert: async (alertId: string): Promise<ApiResponse<SystemAlert>> => {
    return api.put<SystemAlert>(
      { category: 'admin', endpoint: `system/alerts/${alertId}/acknowledge` }
    );
  },

  getSystemMetrics: async (): Promise<ApiResponse<{
    // Server Metrics
    cpuUsage: number;
    memoryUsage: number;
    diskUsage: number;
    activeConnections: number;
    responseTime: number;
    uptime: number;
    
    // Cloudflare Metrics
    cloudflare: {
      requests: number;
      bandwidth: number;
      cacheHitRate: number;
      securityEvents: number;
      workerInvocations: number;
    };
    
    // AI Service Status
    aiServices: {
      gemini: boolean;
      openai: boolean;
      ollama: boolean;
      huggingface: boolean;
      stability: boolean;
      luma: boolean;
      runway: boolean;
      pika: boolean;
      geminiVeo: boolean;
    };
    
    // External Services
    externalServices: {
      cloudinary: boolean;
      slack: boolean;
      sentry: boolean;
      email: boolean;
      database: boolean;
    };
  }>> => {
    return api.get(
      { category: 'admin', endpoint: 'system/metrics' }
    );
  },

  // ==================== CLOUDFLARE MANAGEMENT ====================
  getCloudflareMetrics: async (): Promise<ApiResponse<{
    analytics: any;
    workers: any;
    d1: any;
    kv: any;
    r2: any;
  }>> => {
    return api.get(
      { category: 'admin', endpoint: 'cloudflare/metrics' }
    );
  },

  purgeCloudflareCache: async (): Promise<ApiResponse<{ purged: boolean }>> => {
    return api.post<{ purged: boolean }>(
      { category: 'admin', endpoint: 'cloudflare/cache/purge' }
    );
  },

  deployWorker: async (workerName: string, script: string): Promise<ApiResponse<{ deployed: boolean }>> => {
    return api.post<{ deployed: boolean }>(
      { category: 'admin', endpoint: 'cloudflare/workers/deploy' },
      { workerName, script }
    );
  },

  // ==================== SETTINGS ====================
  getAdminSettings: async (): Promise<ApiResponse<AdminSettings>> => {
    return api.get<AdminSettings>(
      { category: 'admin', endpoint: 'settings' }
    );
  },

  updateAdminSettings: async (
    settings: Partial<AdminSettings>
  ): Promise<ApiResponse<AdminSettings>> => {
    return api.put<AdminSettings>(
      { category: 'admin', endpoint: 'settings' },
      settings
    );
  },

  updateAISettings: async (
    settings: {
      enabled?: boolean;
      defaultModel?: string;
      rateLimit?: number;
      videoGeneration?: boolean;
      imageGeneration?: boolean;
    }
  ): Promise<ApiResponse<AdminSettings>> => {
    return api.put<AdminSettings>(
      { category: 'admin', endpoint: 'settings/ai' },
      settings
    );
  },

  // ==================== EXPORT & BACKUP ====================
  exportReport: async (
    type: 'orders' | 'users' | 'revenue' | 'analytics' | 'ai' | 'all',
    options: ExportOptions,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<{ url: string; filename: string; expiresAt: string }>> => {
    return api.post<{ url: string; filename: string; expiresAt: string }>(
      { category: 'admin', endpoint: 'export' },
      {
        type,
        options,
        filters,
      }
    );
  },

  exportAIData: async (
    model: string,
    options: ExportOptions
  ): Promise<ApiResponse<{ url: string; filename: string; expiresAt: string }>> => {
    return api.post<{ url: string; filename: string; expiresAt: string }>(
      { category: 'admin', endpoint: 'export/ai' },
      {
        model,
        options,
      }
    );
  },

  createBackup: async (): Promise<ApiResponse<{ 
    backupId: string; 
    createdAt: string; 
    size: number;
    location: string;
    type: 'full' | 'incremental';
  }>> => {
    return api.post<{ 
      backupId: string; 
      createdAt: string; 
      size: number;
      location: string;
      type: 'full' | 'incremental';
    }>(
      { category: 'admin', endpoint: 'backup' }
    );
  },

  restoreBackup: async (backupId: string): Promise<ApiResponse<{ message: string }>> => {
    return api.post<{ message: string }>(
      { category: 'admin', endpoint: `backup/${backupId}/restore` }
    );
  },

  // ==================== BULK OPERATIONS ====================
  bulkUpdateOrders: async (
    orderIds: string[],
    updates: Partial<{
      status: Order['status'];
      category: string;
      assignedTo: string;
      priority: number;
    }>
  ): Promise<ApiResponse<{ updated: number; failed: number; details: any[] }>> => {
    return api.post<{ updated: number; failed: number; details: any[] }>(
      { category: 'admin', endpoint: 'orders/bulk-update' },
      { orderIds, updates }
    );
  },

  bulkUpdateUsers: async (
    userIds: string[],
    updates: Partial<{
      status: UserReport['status'];
      role: string;
      subscriptionPlan: string;
      aiRateLimit: number;
    }>
  ): Promise<ApiResponse<{ updated: number; failed: number; details: any[] }>> => {
    return api.post<{ updated: number; failed: number; details: any[] }>(
      { category: 'admin', endpoint: 'users/bulk-update' },
      { userIds, updates }
    );
  },

  sendBulkNotifications: async (
    userIds: string[],
    notification: {
      title: string;
      message: string;
      type: 'email' | 'push' | 'both' | 'slack';
      data?: Record<string, any>;
    }
  ): Promise<ApiResponse<{ sent: number; failed: number }>> => {
    return api.post<{ sent: number; failed: number }>(
      { category: 'admin', endpoint: 'notifications/bulk' },
      { userIds, notification }
    );
  },

  // ==================== UTILITIES ====================
  clearCache: async (cacheType?: 'all' | 'data' | 'images' | 'api' | 'ai'): Promise<ApiResponse<{ cleared: string[] }>> => {
    return api.post<{ cleared: string[] }>(
      { category: 'admin', endpoint: 'cache/clear' },
      { cacheType }
    );
  },

  sendTestEmail: async (
    email: string,
    template?: string
  ): Promise<ApiResponse<{ message: string }>> => {
    return api.post<{ message: string }>(
      { category: 'admin', endpoint: 'test/email' },
      { email, template }
    );
  },

  sendTestSlack: async (
    channel: string,
    message: string
  ): Promise<ApiResponse<{ message: string }>> => {
    return api.post<{ message: string }>(
      { category: 'admin', endpoint: 'test/slack' },
      { channel, message }
    );
  },

  testAIService: async (
    service: string,
    prompt?: string
  ): Promise<ApiResponse<any>> => {
    return api.post<any>(
      { category: 'admin', endpoint: 'test/ai' },
      { service, prompt }
    );
  },

  checkSystemHealth: async (): Promise<ApiResponse<{
    status: 'healthy' | 'degraded' | 'unhealthy';
    components: {
      database: boolean;
      redis: boolean;
      storage: boolean;
      email: boolean;
      api: boolean;
      auth: boolean;
      
      // AI Services
      gemini: boolean;
      openai: boolean;
      ollama: boolean;
      huggingface: boolean;
      stability: boolean;
      
      // Video Services
      luma: boolean;
      runway: boolean;
      pika: boolean;
      geminiVeo: boolean;
      
      // External Services
      cloudflare: boolean;
      cloudinary: boolean;
      slack: boolean;
      sentry: boolean;
    };
    issues: Array<{
      component: string;
      issue: string;
      severity: 'low' | 'medium' | 'high' | 'critical';
      recommendedAction: string;
    }>;
    lastCheck: string;
  }>> => {
    return api.get(
      { category: 'admin', endpoint: 'health' }
    );
  },

  // ==================== LOGS & AUDIT ====================
  getAuditLogs: async (
    params: PaginationParams & {
      userId?: string;
      action?: string;
      startDate?: string;
      endDate?: string;
      ip?: string;
      service?: string;
    } = {}
  ): Promise<ApiResponse<UserActivity[]>> => {
    return api.getPaginated<UserActivity>(
      { category: 'admin', endpoint: 'audit-logs' },
      params
    );
  },

  getErrorLogs: async (
    params: PaginationParams & {
      level?: 'error' | 'warning' | 'info' | 'debug';
      startDate?: string;
      endDate?: string;
      service?: string;
      source?: string;
    } = {}
  ): Promise<ApiResponse<Array<{
    id: string;
    timestamp: string;
    level: string;
    message: string;
    service: string;
    source: string;
    stackTrace?: string;
    userId?: string;
    ip?: string;
    metadata?: Record<string, any>;
  }>>> => {
    return api.getPaginated(
      { category: 'admin', endpoint: 'error-logs' },
      params
    );
  },

  getAIErrorLogs: async (
    params: PaginationParams & {
      model?: string;
      type?: string;
      startDate?: string;
      endDate?: string;
    } = {}
  ): Promise<ApiResponse<any>> => {
    return api.getPaginated(
      { category: 'admin', endpoint: 'error-logs/ai' },
      params
    );
  },
};

// ==================== HELPER FUNCTIONS ====================
export const adminHelpers = {
  /**
   * Format dashboard data for charts
   */
  formatChartData: (
    data: DashboardData,
    chartType: 'line' | 'bar' | 'pie' | 'donut' | 'radar'
  ): any[] => {
    switch (chartType) {
      case 'line':
        return data.revenueTrend.map(item => ({
          period: item.period,
          revenue: item.revenue,
          aiRevenue: item.aiRevenue,
          serviceRevenue: item.serviceRevenue,
        }));
      case 'bar':
        return [
          { name: 'Users', value: data.stats.totalUsers },
          { name: 'Orders', value: data.stats.totalOrders },
          { name: 'Revenue', value: data.stats.totalRevenue },
          { name: 'AI Requests', value: data.stats.totalAIRequests },
          { name: 'Video Generations', value: data.stats.videoGenerationRequests },
        ];
      case 'pie':
        return Object.entries(data.aiMetrics.usageByCategory).map(([category, count]) => ({ 
          category, 
          count 
        }));
      case 'radar':
        return Object.entries(data.performanceMetrics.cloudflare).map(([key, value]) => ({
          metric: key,
          value: typeof value === 'number' ? value : 0,
        }));
      default:
        return [];
    }
  },

  /**
   * Calculate dashboard metrics changes
   */
  calculateMetricsChange: (
    current: DashboardStats,
    previous?: DashboardStats
  ): Record<string, { value: number; change: number; trend: 'up' | 'down' | 'stable' }> => {
    if (!previous) return {};

    const changes: Record<string, { value: number; change: number; trend: 'up' | 'down' | 'stable' }> = {};
    const keys = Object.keys(current) as Array<keyof DashboardStats>;

    keys.forEach(key => {
      const currentVal = current[key] as number;
      const previousVal = previous[key] as number;

      if (typeof currentVal === 'number' && typeof previousVal === 'number' && previousVal !== 0) {
        const change = ((currentVal - previousVal) / previousVal) * 100;
        changes[key] = {
          value: currentVal,
          change: parseFloat(change.toFixed(2)),
          trend: change > 5 ? 'up' : change < -5 ? 'down' : 'stable',
        };
      }
    });

    return changes;
  },

  /**
   * Generate export filename
   */
  generateExportFilename: (
    type: string,
    timeRange: string = 'all'
  ): string => {
    const timestamp = new Date().toISOString().split('T')[0];
    return `nawthtech-${type}-${timeRange}-${timestamp}`;
  },

  /**
   * Check if admin has permission
   */
  hasPermission: (
    requiredPermission: string,
    userPermissions?: string[]
  ): boolean => {
    if (!userPermissions) return false;
    
    // Check direct permission
    if (userPermissions.includes(requiredPermission)) {
      return true;
    }

    // Check wildcard permissions
    if (userPermissions.includes('*') || userPermissions.includes('admin.*')) {
      return true;
    }

    // Check category permission
    const [category] = requiredPermission.split('.');
    if (userPermissions.includes(`${category}.*`)) {
      return true;
    }

    // Check specific AI permissions
    if (requiredPermission.startsWith('ai.')) {
      const aiPermission = requiredPermission.split('.')[1];
      if (userPermissions.includes(`ai.${aiPermission}`)) {
        return true;
      }
    }

    return false;
  },

  /**
   * Validate admin filters
   */
  validateFilters: (filters: AnalyticsFilters): string[] => {
    const errors: string[] = [];

    if (filters.timeRange === 'custom') {
      if (!filters.startDate || !filters.endDate) {
        errors.push('Custom time range requires both start and end dates');
      } else {
        const start = new Date(filters.startDate);
        const end = new Date(filters.endDate);

        if (start > end) {
          errors.push('Start date must be before end date');
        }

        // Limit custom range to 2 years for AI analytics
        const maxRange = 2 * 365 * 24 * 60 * 60 * 1000;
        if (end.getTime() - start.getTime() > maxRange) {
          errors.push('Custom time range cannot exceed 2 years for detailed analytics');
        }
      }
    }

    return errors;
  },

  /**
   * Format AI model display name
   */
  formatAIModelName: (model: string): string => {
    const modelNames: Record<string, string> = {
      'gemini': 'Google Gemini',
      'openai': 'OpenAI GPT',
      'ollama': 'Ollama LLM',
      'huggingface': 'Hugging Face',
      'stability': 'Stability AI',
      'luma': 'Luma AI',
      'runway': 'Runway ML',
      'pika': 'Pika Labs',
      'gemini-veo': 'Gemini Veo',
    };

    return modelNames[model] || model;
  },

  /**
   * Calculate AI cost estimate
   */
  estimateAICost: (
    model: string,
    requests: number,
    type: 'text' | 'image' | 'video'
  ): number => {
    const costPerRequest: Record<string, number> = {
      'gemini-text': 0.001,
      'openai-text': 0.002,
      'ollama-text': 0.0001,
      'huggingface-text': 0.0005,
      'stability-image': 0.01,
      'stability-video': 0.05,
      'luma-video': 0.02,
      'runway-video': 0.03,
      'pika-video': 0.015,
      'gemini-veo-video': 0.04,
    };

    const key = `${model}-${type}`;
    return (costPerRequest[key] || 0.001) * requests;
  },

  /**
   * Get service icon based on platform
   */
  getServiceIcon: (platform: string): string => {
    const icons: Record<string, string> = {
      'instagram': '📸',
      'tiktok': '🎵',
      'twitter': '🐦',
      'youtube': '🎥',
      'facebook': '📘',
      'ai': '🤖',
      'video': '🎬',
      'image': '🖼️',
      'analytics': '📊',
    };

    return icons[platform.toLowerCase()] || '📦';
  },

  /**
   * Format currency
   */
  formatCurrency: (amount: number, currency: string = 'SAR'): string => {
    return new Intl.NumberFormat('ar-SA', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    }).format(amount);
  },

  /**
   * Format large numbers
   */
  formatNumber: (num: number): string => {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
  },

  /**
   * Get severity color
   */
  getSeverityColor: (severity: SystemAlert['severity']): string => {
    const colors: Record<SystemAlert['severity'], string> = {
      'low': '#3fb950',
      'medium': '#e3b341',
      'high': '#ff7b72',
      'critical': '#f85149',
    };
    return colors[severity];
  },
};

// ==================== DEFAULT EXPORT ====================
export default adminAPI;