import axios from 'axios';
import * as Sentry from '@sentry/react';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_role');
      window.location.href = '/login';
    }
    
    Sentry.captureException(error);
    return Promise.reject(error);
  }
);

// AI Services
export const aiService = {
  // Gemini
  async generateWithGemini(prompt: string, options?: any) {
    const response = await api.post('/ai/gemini/generate', { prompt, ...options });
    return response.data;
  },

  // OpenAI
  async generateWithOpenAI(prompt: string, model = 'gpt-4') {
    const response = await api.post('/ai/openai/generate', { prompt, model });
    return response.data;
  },

  // Ollama
  async generateWithOllama(prompt: string, model = 'llama2') {
    const response = await api.post('/ai/ollama/generate', { prompt, model });
    return response.data;
  },

  // Hugging Face
  async generateWithHuggingFace(prompt: string, model: string) {
    const response = await api.post('/ai/huggingface/generate', { prompt, model });
    return response.data;
  },

  // Stability AI
  async generateImageWithStability(prompt: string, options?: any) {
    const response = await api.post('/ai/stability/generate-image', { prompt, ...options });
    return response.data;
  },

  async generateVideoWithStability(prompt: string, options?: any) {
    const response = await api.post('/ai/stability/generate-video', { prompt, ...options });
    return response.data;
  },

  // Video Generation Services
  async generateVideoWithGeminiVeo(prompt: string, options?: any) {
    const response = await api.post('/ai/video/gemini-veo', { prompt, ...options });
    return response.data;
  },

  async generateVideoWithLuma(prompt: string, options?: any) {
    const response = await api.post('/ai/video/luma', { prompt, ...options });
    return response.data;
  },

  async generateVideoWithRunway(prompt: string, options?: any) {
    const response = await api.post('/ai/video/runway', { prompt, ...options });
    return response.data;
  },

  async generateVideoWithPika(prompt: string, options?: any) {
    const response = await api.post('/ai/video/pika', { prompt, ...options });
    return response.data;
  },

  // AI Validation
  async validateContent(content: string, type: string) {
    const response = await api.post('/ai/validation/content', { content, type });
    return response.data;
  },

  async validateOrder(orderData: any) {
    const response = await api.post('/ai/validation/order', orderData);
    return response.data;
  },
};

// Store Services
export const storeService = {
  async getServices(filters?: any) {
    const response = await api.get('/store/services', { params: filters });
    return response.data;
  },

  async getServiceById(id: string) {
    const response = await api.get(`/store/services/${id}`);
    return response.data;
  },

  async createOrder(orderData: any) {
    const response = await api.post('/store/orders', orderData);
    return response.data;
  },

  async getCart() {
    const response = await api.get('/store/cart');
    return response.data;
  },

  async addToCart(item: any) {
    const response = await api.post('/store/cart/items', item);
    return response.data;
  },

  async updateCartItem(id: string, quantity: number) {
    const response = await api.patch(`/store/cart/items/${id}`, { quantity });
    return response.data;
  },

  async removeCartItem(id: string) {
    const response = await api.delete(`/store/cart/items/${id}`);
    return response.data;
  },
};

// Cloudinary Services
export const cloudinaryService = {
  async uploadImage(file: File, options?: any) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('upload_preset', import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET);
    
    if (options) {
      Object.entries(options).forEach(([key, value]) => {
        formData.append(key, value as string);
      });
    }

    const response = await api.post('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  async uploadVideo(file: File, options?: any) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('upload_preset', import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET);
    
    if (options) {
      Object.entries(options).forEach(([key, value]) => {
        formData.append(key, value as string);
      });
    }

    const response = await api.post('/upload/video', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
};

// Payment Services
export const paymentService = {
  async createPaymentIntent(amount: number, currency: string = 'SAR') {
    const response = await api.post('/payments/create-intent', { amount, currency });
    return response.data;
  },

  async confirmPayment(paymentId: string) {
    const response = await api.post('/payments/confirm', { paymentId });
    return response.data;
  },

  async getPaymentMethods() {
    const response = await api.get('/payments/methods');
    return response.data;
  },
};

// Email Services (Cloudflare Workers)
export const emailService = {
  async sendWelcomeEmail(email: string, name: string) {
    const response = await api.post('/email/welcome', { email, name });
    return response.data;
  },

  async sendOrderConfirmation(orderId: string) {
    const response = await api.post('/email/order-confirmation', { orderId });
    return response.data;
  },

  async sendPasswordReset(email: string) {
    const response = await api.post('/email/password-reset', { email });
    return response.data;
  },

  async sendNotification(email: string, subject: string, body: string) {
    const response = await api.post('/email/notification', { email, subject, body });
    return response.data;
  },
};

// Slack Integration
export const slackService = {
  async sendMessage(channel: string, message: string) {
    const response = await api.post('/slack/message', { channel, message });
    return response.data;
  },

  async sendOrderNotification(orderData: any) {
    const response = await api.post('/slack/order-notification', orderData);
    return response.data;
  },

  async sendErrorNotification(error: Error, context?: any) {
    const response = await api.post('/slack/error-notification', { error, context });
    return response.data;
  },
};

// Analytics Services
export const analyticsService = {
  async trackEvent(event: string, properties?: any) {
    const response = await api.post('/analytics/track', { event, properties });
    return response.data;
  },

  async getDashboardStats(period: string = '30d') {
    const response = await api.get('/analytics/dashboard', { params: { period } });
    return response.data;
  },

  async getUserAnalytics(userId: string) {
    const response = await api.get(`/analytics/user/${userId}`);
    return response.data;
  },

  async getStoreAnalytics() {
    const response = await api.get('/analytics/store');
    return response.data;
  },
};

// Content Services
export const contentService = {
  async generateContent(prompt: string, options?: any) {
    const response = await api.post('/content/generate', { prompt, ...options });
    return response.data;
  },

  async analyzeContent(content: string) {
    const response = await api.post('/content/analyze', { content });
    return response.data;
  },

  async optimizeContent(content: string, target: string) {
    const response = await api.post('/content/optimize', { content, target });
    return response.data;
  },

  async getContentHistory() {
    const response = await api.get('/content/history');
    return response.data;
  },
};

// Strategy Services
export const strategyService = {
  async createStrategy(data: any) {
    const response = await api.post('/strategies/create', data);
    return response.data;
  },

  async analyzeStrategy(data: any) {
    const response = await api.post('/strategies/analyze', data);
    return response.data;
  },

  async updateStrategyProgress(strategyId: string, progress: number) {
    const response = await api.patch(`/strategies/${strategyId}/progress`, { progress });
    return response.data;
  },

  async getStrategies() {
    const response = await api.get('/strategies');
    return response.data;
  },
};

// Report Services
export const reportService = {
  async generateReport(type: string, data: any) {
    const response = await api.post('/reports/generate', { type, data });
    return response.data;
  },

  async getReport(reportId: string) {
    const response = await api.get(`/reports/${reportId}`);
    return response.data;
  },

  async getReportList() {
    const response = await api.get('/reports');
    return response.data;
  },
};

// Service Manager
export const serviceManager = {
  async initializeServices() {
    const response = await api.post('/services/initialize');
    return response.data;
  },

  async getServiceStatus() {
    const response = await api.get('/services/status');
    return response.data;
  },

  async healthCheck() {
    const response = await api.get('/services/health');
    return response.data;
  },
};

export default api;