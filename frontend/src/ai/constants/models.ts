export const AI_MODELS = [
  {
    id: 'gemini-2.0-flash',
    name: 'Gemini 2.0 Flash',
    provider: 'Google',
    capabilities: ['text', 'translation', 'summarization'],
    isLocal: false,
    isFree: true,
    maxTokens: 8192,
    languages: ['ar', 'en', 'fr', 'es', 'de'],
  },
  {
    id: 'llama3.2-3b',
    name: 'Llama 3.2 3B',
    provider: 'Meta',
    capabilities: ['text', 'code', 'reasoning'],
    isLocal: true,
    isFree: true,
    maxTokens: 4096,
    languages: ['en', 'ar', 'es'],
  },
  {
    id: 'mistral-7b',
    name: 'Mistral 7B',
    provider: 'Mistral AI',
    capabilities: ['text', 'translation', 'summarization'],
    isLocal: true,
    isFree: true,
    maxTokens: 32768,
    languages: ['en', 'fr', 'es', 'de', 'it'],
  },
  {
    id: 'qwen2.5-7b',
    name: 'Qwen 2.5 7B',
    provider: 'Alibaba',
    capabilities: ['text', 'translation', 'code', 'reasoning'],
    isLocal: true,
    isFree: true,
    maxTokens: 32768,
    languages: ['en', 'zh', 'ar', 'fr', 'es'],
  },
];

export const CONTENT_TYPES = [
  { id: 'blog_post', name: 'مقال مدونة', icon: '📝' },
  { id: 'social_media', name: 'منشور وسائط اجتماعية', icon: '📱' },
  { id: 'email', name: 'بريد إلكتروني', icon: '📧' },
  { id: 'ad_copy', name: 'نص إعلامي', icon: '📢' },
  { id: 'product_description', name: 'وصف منتج', icon: '📦' },
];

export const TONES = [
  { id: 'professional', name: 'مهني', description: 'مناسب للأعمال والشركات' },
  { id: 'casual', name: 'غير رسمي', description: 'مناسب للوسائط الاجتماعية' },
  { id: 'persuasive', name: 'إقناعي', description: 'مناسب للإعلانات والمبيعات' },
  { id: 'informative', name: 'إعلامي', description: 'مناسب للمحتوى التعليمي' },
];

export const LANGUAGES = [
  { id: 'ar', name: 'العربية', nativeName: 'العربية', flag: '🇸🇦' },
  { id: 'en', name: 'English', nativeName: 'English', flag: '🇺🇸' },
  { id: 'fr', name: 'Français', nativeName: 'Français', flag: '🇫🇷' },
  { id: 'es', name: 'Español', nativeName: 'Español', flag: '🇪🇸' },
];

export const MEDIA_STYLES = [
  { id: 'realistic', name: 'واقعي', description: 'صور فوتوغرافية واقعية' },
  { id: 'anime', name: 'أنمي', description: 'أسلوب رسوم متحركة ياباني' },
  { id: 'digital_art', name: 'فن رقمي', description: 'فن رقمي وإبداعي' },
  { id: '3d_render', name: 'ثلاثي الأبعاد', description: 'تصميم ثلاثي الأبعاد' },
  { id: 'minimalist', name: 'بسيط', description: 'تصميم بسيط وأنيق' },
];