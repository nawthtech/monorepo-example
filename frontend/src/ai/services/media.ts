import { aiService } from './api';

export interface MediaGenerationOptions {
  style: 'realistic' | 'anime' | 'digital_art' | '3d_render' | 'minimalist';
  aspectRatio: '1:1' | '16:9' | '9:16' | '4:5';
  quality: 'low' | 'medium' | 'high';
}

export class MediaService {
  // توليد صورة لمنشور اجتماعي
  async generateSocialMediaImage(
    platform: 'instagram' | 'twitter' | 'linkedin' | 'facebook',
    theme: string,
    text?: string
  ) {
    const dimensions = this.getPlatformDimensions(platform);
    const style = this.getPlatformStyle(platform);
    
    const prompt = this.buildImagePrompt(theme, text, style);
    
    return await aiService.generateImage(prompt, style);
  }
  
  // توليد صورة منتج
  async generateProductImage(productName: string, description: string) {
    const prompt = `Create a professional product image for: ${productName}
    
    Description: ${description}
    
    Style requirements:
    - Clean white background
    - Professional photography style
    - 3D render quality
    - Modern and minimalist
    - Good lighting and shadows
    - Product centered`;
    
    return await aiService.generateImage(prompt, 'realistic');
  }
  
  // توليد فيديو تعريفي
  async generateExplainerVideo(topic: string, duration: number = 30) {
    const prompt = `Create an animated explainer video about: ${topic}
    
    Requirements:
    - Duration: ${duration} seconds
    - Animated style
    - Professional quality
    - Clear visual storytelling
    - Engaging for business audience
    - Include text overlays if needed`;
    
    return await aiService.generateVideo(prompt, duration);
  }
  
  // توليد شعار
  async generateLogo(companyName: string, industry: string, style: string = 'modern') {
    const prompt = `Design a logo for ${companyName} in the ${industry} industry
    
    Style: ${style}
    
    Requirements:
    - Modern and professional
    - Scalable vector design
    - Works in color and black/white
    - Memorable and distinctive
    - Represents innovation and technology
    - Use brand colors if specified
    
    Describe the logo design in detail.`;
    
    return await aiService.generateContent({ prompt });
  }
  
  // أدوات مساعدة
  private getPlatformDimensions(platform: string): { width: number; height: number } {
    const dimensions: Record<string, { width: number; height: number }> = {
      instagram: { width: 1080, height: 1080 },
      twitter: { width: 1200, height: 675 },
      linkedin: { width: 1200, height: 627 },
      facebook: { width: 1200, height: 630 },
    };
    
    return dimensions[platform] || { width: 1024, height: 1024 };
  }
  
  private getPlatformStyle(platform: string): string {
    const styles: Record<string, string> = {
      instagram: 'modern, vibrant, eye-catching',
      twitter: 'clean, professional, brand-aligned',
      linkedin: 'corporate, professional, business-focused',
      facebook: 'social, engaging, community-oriented',
    };
    
    return styles[platform] || 'realistic';
  }
  
  private buildImagePrompt(theme: string, text?: string, style?: string): string {
    let prompt = `Create a social media image about: ${theme}`;
    
    if (text) {
      prompt += `\nText to include: "${text}"`;
    }
    
    if (style) {
      prompt += `\nStyle: ${style}`;
    }
    
    prompt += `
    
    Design requirements:
    - Use brand colors: purple (#7A3EF0) and neon cyan (#00F6FF)
    - Modern, futuristic AI/tech aesthetic
    - Clean, professional layout
    - High quality, sharp image
    - Optimized for digital display
    
    Make it visually striking and shareable.`;
    
    return prompt;
  }
}

export const mediaService = new MediaService();