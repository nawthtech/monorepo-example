import '@testing-library/jest-dom'
import { expect, afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import * as matchers from '@testing-library/jest-dom/matchers'
import React from 'react'

// توسيع expect بـ jest-dom matchers
expect.extend(matchers)

// 1. Mock لـ window.scrollTo
Object.defineProperty(window, 'scrollTo', {
  value: vi.fn(),
  writable: true,
})

// 2. Mock لـ window.matchMedia (مطلوب لـ MUI)
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// 3. Mock لـ window.localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
})

// 4. Mock لـ window.sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
}
Object.defineProperty(window, 'sessionStorage', {
  value: sessionStorageMock,
})

// 5. Mock لـ EventSource مع تحسينات
class MockEventSource {
  static readonly CONNECTING = 0
  static readonly OPEN = 1
  static readonly CLOSED = 2

  public onopen: (() => void) | null = null
  public onmessage: ((event: any) => void) | null = null
  public onerror: ((event: any) => void) | null = null
  public readyState = MockEventSource.CONNECTING
  public url: string
  public withCredentials = false

  constructor(url: string, eventSourceInitDict?: EventSourceInit) {
    this.url = url
    if (eventSourceInitDict?.withCredentials !== undefined) {
      this.withCredentials = eventSourceInitDict.withCredentials
    }
    
    // محاكاة فتح الاتصال
    setTimeout(() => {
      this.readyState = MockEventSource.OPEN
      if (this.onopen) this.onopen()
    }, 10)
  }

  close() {
    this.readyState = MockEventSource.CLOSED
  }

  // دالة مساعدة للمختبرات لمحاكاة الرسائل
  simulateMessage(data: string | object) {
    if (this.onmessage && this.readyState === MockEventSource.OPEN) {
      const eventData = typeof data === 'string' ? data : JSON.stringify(data)
      this.onmessage({ data: eventData })
    }
  }

  // دالة مساعدة لمحاكاة الأخطاء
  simulateError() {
    if (this.onerror) {
      this.onerror(new Event('error'))
      this.readyState = MockEventSource.CLOSED
    }
  }

  // دالة مساعدة لمحاكاة إغلاق الاتصال
  simulateClose() {
    this.readyState = MockEventSource.CLOSED
  }

  addEventListener(event: string, listener: any) {
    if (event === 'open') this.onopen = listener
    if (event === 'message') this.onmessage = listener
    if (event === 'error') this.onerror = listener
  }

  removeEventListener(event: string, listener: any) {
    if (event === 'open' && this.onopen === listener) this.onopen = null
    if (event === 'message' && this.onmessage === listener) this.onmessage = null
    if (event === 'error' && this.onerror === listener) this.onerror = null
  }

  dispatchEvent(event: Event): boolean {
    return true
  }
}

// Mock لـ EventSource العالمي
global.EventSource = MockEventSource as any

// 6. Mock للمتغيرات البيئية
vi.stubEnv('VITE_BACKEND_HOST', 'http://localhost:8080')
vi.stubEnv('VITE_API_URL', 'http://localhost:8080/api')
vi.stubEnv('VITE_APP_TITLE', 'NawthTech')
vi.stubEnv('MODE', 'test')

// 7. Mock للاستيرادات الناقصة في الملفات

// Mock لصفحات AI
vi.mock('../pages/AIDashboard/AIDashboard', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'ai-dashboard',
    'role': 'main'
  }, 'AI Dashboard Page')
}))

vi.mock('../pages/ContentGenerator/ContentGenerator', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'content-generator' 
  }, 'Content Generator Page')
}))

vi.mock('../pages/MediaStudio/MediaStudio', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'media-studio' 
  }, 'Media Studio Page')
}))

vi.mock('../pages/StrategyPlanner/StrategyPlanner', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'strategy-planner' 
  }, 'Strategy Planner Page')
}))

// Mock لمكونات AI
vi.mock('../ai/components/AIContentGenerator/AIContentGenerator', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'ai-content-generator' 
  }, 'AI Content Generator')
}))

vi.mock('../ai/components/AIMediaGenerator/AIMediaGenerator', () => ({
  default: () => React.createElement('div', { 
    'data-testid': 'ai-media-generator' 
  }, 'AI Media Generator')
}))

// Mock للـ store
vi.mock('../store', () => ({
  store: {
    getState: vi.fn(() => ({
      auth: {
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
      },
      ai: {
        loading: false,
        error: null,
        responses: [],
      },
    })),
    dispatch: vi.fn(),
    subscribe: vi.fn(() => vi.fn()),
    replaceReducer: vi.fn(),
  },
}))

// Mock لخدمات AI
vi.mock('../ai/hooks/useAI', () => ({
  useAI: () => ({
    loading: false,
    error: null,
    result: null,
    progress: 0,
    generateContent: vi.fn().mockResolvedValue({
      success: true,
      data: { content: 'محتوى نصي مُولد', model_used: 'gpt-3.5' }
    }),
    generateImage: vi.fn().mockResolvedValue({
      success: true,
      data: { url: 'https://example.com/image.png' }
    }),
    getUsage: vi.fn().mockResolvedValue({
      text_used: 100,
      text_limit: 1000,
      images_used: 5,
      images_limit: 50,
      videos_used: 0,
      videos_limit: 10,
    }),
    getAvailableModels: vi.fn().mockResolvedValue(['gpt-3.5', 'gpt-4']),
    cancel: vi.fn(),
    reset: vi.fn(),
  }),
}))

// Mock لـ React Router
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    BrowserRouter: ({ children }: { children: React.ReactNode }) => 
      React.createElement('div', {}, children),
    Routes: ({ children }: { children: React.ReactNode }) => 
      React.createElement('div', {}, children),
    Route: ({ element }: { element: React.ReactNode }) => element,
    Navigate: ({ to }: { to: string }) => 
      React.createElement('div', {}, `Navigate to: ${to}`),
    useNavigate: () => vi.fn(),
    useLocation: () => ({ pathname: '/test' }),
  }
})

// Mock لـ MUI
vi.mock('@mui/material/styles', () => ({
  ThemeProvider: ({ children, theme }: any) => 
    React.createElement('div', { 'data-testid': 'theme-provider' }, children),
  createTheme: vi.fn(() => ({})),
}))

vi.mock('@mui/material/CssBaseline', () => ({
  default: () => React.createElement('div', { 'data-testid': 'css-baseline' }),
}))

// Mock لـ React Redux
vi.mock('react-redux', () => ({
  Provider: ({ children, store }: any) => 
    React.createElement('div', { 'data-testid': 'redux-provider' }, children),
  useSelector: vi.fn(),
  useDispatch: () => vi.fn(),
}))

// Mock للأصول
vi.mock('../assets/mc', () => ({
  mc: (...args: string[]) => args.join(' '),
}))

// 8. Mock لـ ResizeObserver (مطلوب لبعض مكتبات الرسوم البيانية)
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// 9. Mock لـ IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
  takeRecords: vi.fn(() => []),
}))

// 10. Mock لـ fetch إذا لزم الأمر
global.fetch = vi.fn().mockResolvedValue({
  ok: true,
  json: vi.fn().mockResolvedValue({}),
  text: vi.fn().mockResolvedValue(''),
  status: 200,
  headers: new Map(),
})

// 11. Mock لـ requestAnimationFrame
global.requestAnimationFrame = vi.fn((callback) => {
  setTimeout(callback, 0)
  return 0
})

global.cancelAnimationFrame = vi.fn()

// 12. Mock لـ setTimeout و clearTimeout لتسريع الاختبارات
vi.useFakeTimers()

// تنظيف بعد كل اختبار
afterEach(() => {
  cleanup()
  vi.clearAllMocks()
  vi.useRealTimers()
  vi.useFakeTimers() // إعادة تعيين لاختبارات التالية
  localStorageMock.clear()
  sessionStorageMock.clear()
})

// تنظيف بعد كل مجموعة اختبارات
afterAll(() => {
  vi.clearAllTimers()
  vi.useRealTimers()
})