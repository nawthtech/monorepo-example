import '@testing-library/jest-dom'
import { vi } from 'vitest'

// Mock لـ EventSource
global.EventSource = vi.fn(() => ({
  onopen: null,
  onmessage: null,
  onerror: null,
  close: vi.fn(),
  readyState: 0,
}))

// Mock لـ window.scrollTo
Object.defineProperty(window, 'scrollTo', {
  value: vi.fn(),
  writable: true,
})

// Mock لـ window.matchMedia
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

// Mock للمتغيرات البيئية
vi.stubEnv('VITE_BACKEND_HOST', 'http://localhost:8080')

// تنظيف بعد كل اختبار
afterEach(() => {
  vi.clearAllMocks()
})