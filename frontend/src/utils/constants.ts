export const VIEW_TABS = [
  { id: 'all', label: 'جميع الخدمات' },
  { id: 'featured', label: 'المميزة' },
  { id: 'ai', label: 'الذكاء الاصطناعي' },
  { id: 'trending', label: 'الرائجة' },
  { id: 'new', label: 'الجديدة' },
] as const;

export const CATEGORIES = [
  { id: 'all', label: 'جميع الخدمات' },
  { id: 'instagram', label: 'إنستغرام' },
  { id: 'tiktok', label: 'تيك توك' },
  { id: 'twitter', label: 'تويتر' },
  { id: 'youtube', label: 'يوتيوب' },
  { id: 'facebook', label: 'فيسبوك' },
  { id: 'followers', label: 'المتابعين' },
  { id: 'likes', label: 'الإعجابات' },
  { id: 'comments', label: 'التعليقات' },
  { id: 'analytics', label: 'التحليلات' },
] as const;

export const PLATFORM_COLORS = {
  instagram: '#E4405F',
  tiktok: '#000000',
  twitter: '#1DA1F2',
  youtube: '#FF0000',
  facebook: '#1877F2',
  analytics: '#3fb950',
} as const;