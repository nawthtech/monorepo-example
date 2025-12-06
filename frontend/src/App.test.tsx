import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import App from './App'

// Mock لـ EventSource
vi.stubGlobal('EventSource', vi.fn(() => ({
  onopen: null,
  onmessage: null,
  onerror: null,
  close: vi.fn(),
  readyState: 0,
})))

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />)
    // تحقق من وجود نص في الصفحة
    expect(screen.getByText('لوحة تحكم الذكاء الاصطناعي')).toBeInTheDocument()
  })
  
  it('displays AI Dashboard content', () => {
    render(<App />)
    // تحقق من وجود محتوى لوحة التحكم
    expect(screen.getByText('الصفحة الرئيسية لأدوات الذكاء الاصطناعي')).toBeInTheDocument()
  })
  
  it('shows all page titles', () => {
    render(<App />)
    // تحقق من وجود جميع عناوين الصفحات
    expect(screen.getByText('لوحة تحكم الذكاء الاصطناعي')).toBeInTheDocument()
    expect(screen.getByText('مولد المحتوى')).toBeInTheDocument()
    expect(screen.getByText('استوديو الوسائط')).toBeInTheDocument()
    expect(screen.getByText('مخطط الاستراتيجيات')).toBeInTheDocument()
  })
  
  it('has router working', () => {
    render(<App />)
    // تحقق من وجود عنصر التنقل
    expect(screen.getByText('Navigate to: /ai')).toBeInTheDocument()
  })
  
  it('has correct structure', () => {
    render(<App />)
    // تحقق من وجود الهيكل الأساسي
    expect(screen.getByText('Here\'s some unnecessary quotes for you to read...')).toBeInTheDocument()
    expect(screen.getByText('Start Quotes')).toBeInTheDocument()
  })
})