import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import App from './App'

describe('App', () => {
  beforeEach(() => {
    // Reset mocks before each test
    vi.clearAllMocks()
  })

  it('renders main heading and button', () => {
    render(<App />)
    
    expect(screen.getByRole('main')).toBeInTheDocument()
    expect(screen.getByText("Here's some unnecessary quotes for you to read...")).toBeInTheDocument()
    expect(screen.getByText('Start Quotes')).toBeInTheDocument()
  })

  it('toggles button text when clicked', async () => {
    render(<App />)
    
    const button = screen.getByText('Start Quotes')
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(screen.getByText('Stop Quotes')).toBeInTheDocument()
    })
  })

  it('initially shows no messages', () => {
    render(<App />)
    
    const messages = screen.queryAllByRole('paragraph')
    expect(messages.length).toBe(0)
  })

  it('starts and stops connection when button is clicked', async () => {
    render(<App />)
    
    // Start connection
    const button = screen.getByText('Start Quotes')
    fireEvent.click(button)
    
    await waitFor(() => {
      expect(screen.getByText('Stop Quotes')).toBeInTheDocument()
    })
    
    // Stop connection
    fireEvent.click(screen.getByText('Stop Quotes'))
    
    await waitFor(() => {
      expect(screen.getByText('Start Quotes')).toBeInTheDocument()
    })
  })
})