import axios from 'axios';

export interface AIRequest {
  prompt: string;
  model?: string;
  options?: Record<string, any>;
}

export interface AIResponse {
  success: boolean;
  data?: any;
  error?: string;
  usage?: {
    tokens: number;
    cost: number;
  };
}

export const aiService = {
  async generateContent(request: AIRequest): Promise<AIResponse> {
    try {
      const response = await axios.post('/api/ai/generate', request);
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to generate content'
      };
    }
  },

  async getAvailableModels(): Promise<AIResponse> {
    try {
      const response = await axios.get('/api/ai/models');
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to get available models'
      };
    }
  },

  async getUsage(): Promise<AIResponse> {
    try {
      const response = await axios.get('/api/ai/usage');
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        error: error.message || 'Failed to get usage data'
      };
    }
  }
};
