import { vi } from 'vitest'

// Mock global objects
vi.stubGlobal('global', {
  EventSource: vi.fn(),
  fetch: vi.fn(),
});

vi.stubGlobal('afterEach', vi.fn());

// Mock EventSource
vi.stubGlobal('EventSource', vi.fn(() => ({
  onopen: null,
  onmessage: null,
  onerror: null,
  close: vi.fn(),
})));

// Mock window
Object.defineProperty(window, 'scrollTo', {
  value: vi.fn(),
  writable: true,
});
