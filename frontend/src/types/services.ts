export interface Service {
  id: string;
  name: string;
  description: string;
  category: string;
  basePrice: number;
  features: string[];
  platform: string;
  aiPowered: boolean;
  trending: boolean;
  featured: boolean;
  inventory: {
    stock: number;
    lowStock: boolean;
  };
  stats: {
    orders: number;
    successRate: number;
    deliveryTime: string;
  };
}

export interface CartItem {
  id: string;
  serviceId: string;
  name: string;
  price: number;
  quantity: number;
  unit: string;
  platform: string;
}

export interface Order {
  id: string;
  orderNumber: string;
  userId: string;
  items: CartItem[];
  totalAmount: number;
  status: 'pending' | 'confirmed' | 'processing' | 'completed' | 'cancelled';
  paymentStatus: 'pending' | 'paid' | 'failed' | 'refunded';
  createdAt: string;
  updatedAt: string;
}

export interface AIServiceRequest {
  prompt: string;
  model?: string;
  options?: Record<string, any>;
}

export interface AIServiceResponse {
  success: boolean;
  data: any;
  usage?: {
    tokens: number;
    cost: number;
  };
  error?: string;
}

export interface CloudinaryUploadResponse {
  asset_id: string;
  public_id: string;
  version: number;
  version_id: string;
  signature: string;
  width: number;
  height: number;
  format: string;
  resource_type: string;
  created_at: string;
  tags: string[];
  bytes: number;
  type: string;
  etag: string;
  placeholder: boolean;
  url: string;
  secure_url: string;
  folder: string;
  original_filename: string;
}

export interface EmailTemplate {
  to: string;
  subject: string;
  body: string;
  template?: string;
  variables?: Record<string, any>;
}

export interface SlackMessage {
  channel: string;
  text: string;
  attachments?: any[];
  blocks?: any[];
}

export interface AnalyticsEvent {
  event: string;
  properties?: Record<string, any>;
  userId?: string;
  timestamp?: string;
}

export interface StrategyData {
  title: string;
  description: string;
  type: string;
  goals: string[];
  targetAudience: string;
  budget: number;
  timeline: string;
  metrics: string[];
}

export interface ReportData {
  type: string;
  title: string;
  data: any;
  format: 'pdf' | 'excel' | 'json';
  generatedAt: string;
}