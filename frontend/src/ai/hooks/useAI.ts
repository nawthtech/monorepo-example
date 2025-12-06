import { useState, useCallback, useRef } from 'react';
import { aiService, type AIRequest, type AIResponse } from '../services/api'; 
import { contentService } from '../services/content';
import { analysisService } from '../services/analysis';
import { mediaService } from '../services/media';

// تعريف أنواع اللغة
type Language = 'ar' | 'en' | 'fr' | 'es';

interface UseAIOptions {
  onSuccess?: (data: any) => void;
  onError?: (error: Error) => void;
  showNotifications?: boolean;
}

// تعريف أنواع الدوال المساعدة
interface BlogPostOptions {
  language?: Language;
  length?: number;
  tone?: string;
}

interface SocialMediaOptions {
  language?: Language;
  hashtags?: boolean;
  emojis?: boolean;
}

export const useAI = (options: UseAIOptions = {}) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<any>(null);
  const [progress, setProgress] = useState<number>(0);
  
  const abortControllerRef = useRef<AbortController | null>(null);
  
  // توليد محتوى نصي عام
  const generateContent = useCallback(async (request: AIRequest): Promise<AIResponse> => {
    setLoading(true);
    setError(null);
    setProgress(0);
    
    abortControllerRef.current = new AbortController();
    
    try {
      // محاكاة التقدم
      const progressInterval = setInterval(() => {
        setProgress(prev => Math.min(prev + 10, 90));
      }, 500);
      
      const response = await aiService.generateContent(request);
      
      clearInterval(progressInterval);
      setProgress(100);
      setResult(response);
      
      if (options.onSuccess) {
        options.onSuccess(response);
      }
      
      return response;
    } catch (err: any) {
      setError(err.message || 'An error occurred');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
      setTimeout(() => setProgress(0), 1000);
    }
  }, [options]);
  
  // توليد مقال
  const generateBlogPost = useCallback(async (
    topic: string, 
    options: BlogPostOptions = {}
  ): Promise<any> => {
    const { language = 'ar', length = 500, tone = 'professional' } = options;
    
    setLoading(true);
    setError(null);
    
    try {
      // تحويل language إلى النوع الصحيح
      const lang = language as Language;
      const result = await contentService.generateBlogPost(topic, { 
        language: lang, 
        length, 
        tone 
      });
      
      setResult(result);
      if (options.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to generate blog post');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [options.onSuccess, options.onError]);
  
  // توليد منشور وسائط اجتماعية
  const generateSocialMediaPost = useCallback(async (
    platform: 'twitter' | 'linkedin' | 'instagram' | 'facebook',
    topic: string,
    options: SocialMediaOptions = {}
  ): Promise<any> => {
    const { language = 'ar', hashtags = true, emojis = true } = options;
    
    setLoading(true);
    setError(null);
    
    try {
      // تحويل language إلى النوع الصحيح
      const lang = language as Language;
      const result = await contentService.generateSocialMediaPost(platform, topic, { 
        language: lang, 
        hashtags, 
        emojis 
      });
      
      setResult(result);
      if (options.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || `Failed to generate ${platform} post`);
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [options.onSuccess, options.onError]);
  
  // تحليل اتجاهات السوق
  const analyzeMarketTrends = useCallback(async (
    industry: string, 
    timeframe: string
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await analysisService.analyzeMarketTrends(industry, timeframe);
      
      setResult(result);
      if (options.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to analyze market trends');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [options.onSuccess, options.onError]);
  
  // توليد صورة
  const generateImage = useCallback(async (
    prompt: string, 
    style: string = 'realistic',
    options?: { aspectRatio?: string; size?: string }
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await mediaService.generateSocialMediaImage('instagram', prompt, style, options);
      setResult(result);
      
      if (options?.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to generate image');
      if (options?.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  // إلغاء العملية الجارية
  const cancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      setLoading(false);
      setError('Operation cancelled');
    }
  }, []);
  
  // إعادة التعيين
  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
    setResult(null);
    setProgress(0);
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }, []);
  
  // الحصول على النماذج المتاحة
  const getAvailableModels = useCallback(async (): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await aiService.getAvailableModels();
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to get available models');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  // الحصول على الاستخدام
  const getUsage = useCallback(async (): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await aiService.getUsage();
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to get usage data');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  // توليد نص مع نموذج ولغة محددة (لإصلاح الخطأ السابق)
  const generateText = useCallback(async (
    prompt: string, 
    model?: string, 
    language?: string
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      // تحويل language إلى النوع الصحيح إذا كانت معطاة
      const lang = language ? language as Language : undefined;
      
      // هنا يمكنك استدعاء خدمة النص المناسبة
      const request: AIRequest = {
        prompt,
        type: 'text',
        options: {
          model,
          language: lang,
        }
      };
      
      const result = await aiService.generateContent(request);
      setResult(result);
      
      if (options.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to generate text');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [options.onSuccess, options.onError]);
  
  // ترجمة نص
  const translateText = useCallback(async (
    text: string,
    targetLanguage: Language,
    sourceLanguage?: Language
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      // تحويل اللغات إلى النوع الصحيح
      const targetLang = targetLanguage as Language;
      const sourceLang = sourceLanguage as Language | undefined;
      
      const result = await contentService.translateText(text, targetLang, sourceLang);
      setResult(result);
      
      if (options.onSuccess) {
        options.onSuccess(result);
      }
      
      return result;
    } catch (err: any) {
      setError(err.message || 'Failed to translate text');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [options.onSuccess, options.onError]);
  
  return {
    // State
    loading,
    error,
    result,
    progress,
    
    // Actions
    generateContent,
    generateBlogPost,
    generateSocialMediaPost,
    analyzeMarketTrends,
    generateImage,
    generateText,
    translateText,
    getAvailableModels,
    getUsage,
    
    // Control
    cancel,
    reset,
    
    // Services
    contentService,
    analysisService,
    mediaService,
  };
};