import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import App from './App'

// Mock لـ EventSource
const mockEventSource = {
  onopen: null,
  onmessage: null,
  onerror: null,
  close: vi.fn(),
  readyState: 0,
}

describe('App Component', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset EventSource mock
    global.EventSource = vi.fn(() => mockEventSource) as any
  })

  it('renders without crashing', () => {
    render(<App />)
    
    // تحقق من وجود العناصر الأساسية
    expect(screen.getByTestId('redux-provider')).toBeInTheDocument()
    expect(screen.getByTestId('theme-provider')).toBeInTheDocument()
    expect(screen.getByTestId('css-baseline')).toBeInTheDocument()
  })

  it('contains router provider', () => {
    render(<App />)
    expect(screen.getByTestId('ai-dashboard')).toBeInTheDocument()
  })

  it('has correct routing structure', () => {
    render(<App />)
    
    // تحقق من وجود العناصر المتوقعة
    const mainElement = screen.getByRole('main')
    expect(mainElement).toBeInTheDocument()
    
    // تحقق من أن AI Dashboard معروض (الصفحة الافتراضية)
    expect(screen.getByTestId('ai-dashboard')).toBeInTheDocument()
  })

  it('initializes EventSource when connection opens', () => {
    render(<App />)
    
    // تحقق من أن EventSource تم استدعاؤه
    expect(global.EventSource).toHaveBeenCalled()
  })

  it('handles EventSource errors gracefully', () => {
    // محاكاة خطأ في EventSource
    const errorEventSource = {
      onopen: null,
      onmessage: null,
      onerror: vi.fn(),
      close: vi.fn(),
      readyState: 2, // CLOSED
    }
    
    global.EventSource = vi.fn(() => errorEventSource) as any
    
    render(<App />)
    
    // يجب أن يعالج التطبيق الخطأ دون أن يتحطم
    expect(screen.getByTestId('ai-dashboard')).toBeInTheDocument()
  })
})

describe('App Routing', () => {
  it('renders AI Dashboard at /ai route', () => {
    window.history.pushState({}, 'AI Dashboard', '/ai')
    render(<App />)
    expect(screen.getByTestId('ai-dashboard')).toBeInTheDocument()
  })
})