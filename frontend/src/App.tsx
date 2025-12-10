// frontend/src/App.tsx
import React, { useState, useEffect, Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Toaster } from 'react-hot-toast';
import { ErrorBoundary } from 'react-error-boundary';
import * as Sentry from '@sentry/react';

// Mantine UI Framework
import { MantineProvider, createTheme, LoadingOverlay } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { ModalsProvider } from '@mantine/modals';

// Context Providers
import { AuthProvider } from './contexts/AuthContext';
import { CartProvider } from './contexts/CartContext';
import { AIProvider } from './contexts/AIContext';
import { NotificationProvider } from './contexts/NotificationContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { ServicesProvider } from './contexts/ServicesContext';

// Core Components
import Layout from './components/Layout/Layout';
import Header from './components/Layout/Header';
import Sidebar from './components/Layout/Sidebar';
import Navbar from './components/Layout/Navbar';
import Footer from './components/Layout/Footer';
import LoadingScreen from './components/Common/LoadingScreen';
import ErrorFallback from './components/Common/ErrorFallback';
import MaintenanceMode from './components/Common/MaintenanceMode';

// Lazy-loaded Pages (Code Splitting)
const HomePage = lazy(() => import('./pages/HomePage'));
const StorePage = lazy(() => import('./pages/StorePage'));
const ServiceDetailPage = lazy(() => import('./pages/ServiceDetailPage'));
const CartPage = lazy(() => import('./pages/CartPage'));
const CheckoutPage = lazy(() => import('./pages/CheckoutPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const AIStudioPage = lazy(() => import('./pages/AIStudioPage'));
const ContentStudioPage = lazy(() => import('./pages/ContentStudioPage'));
const VideoStudioPage = lazy(() => import('./pages/VideoStudioPage'));
const AnalyticsPage = lazy(() => import('./pages/AnalyticsPage'));
const StrategyPage = lazy(() => import('./pages/StrategyPage'));
const OrdersPage = lazy(() => import('./pages/OrdersPage'));
const ProfilePage = lazy(() => import('./pages/ProfilePage'));
const SettingsPage = lazy(() => import('./pages/SettingsPage'));
const LoginPage = lazy(() => import('./pages/auth/LoginPage'));
const RegisterPage = lazy(() => import('./pages/auth/RegisterPage'));
const ForgotPasswordPage = lazy(() => import('./pages/auth/ForgotPasswordPage'));
const AdminDashboard = lazy(() => import('./pages/admin/AdminDashboard'));

// Sentry Configuration
Sentry.init({
  dsn: import.meta.env.VITE_SENTRY_DSN,
  environment: import.meta.env.VITE_ENVIRONMENT || 'development',
  integrations: [
    new Sentry.BrowserTracing({
      tracePropagationTargets: [window.location.origin],
    }),
    new Sentry.Replay(),
  ],
  tracesSampleRate: 0.1,
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1.0,
});

// Create Query Client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      cacheTime: 1000 * 60 * 30, // 30 minutes
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

// Mantine Theme Configuration (Sentry Dark Theme)
const theme = createTheme({
  colors: {
    dark: [
      '#0d1117', // Primary background
      '#161b22', // Secondary background
      '#21262d', // Tertiary background
      '#30363d', // Border color
      '#484f58',
      '#6e7681',
      '#8b949e',
      '#c9d1d9',
      '#f0f6fc', // Primary text
    ],
    purple: ['#bc8cff', '#7c3aed', '#6e2cc9'],
    blue: ['#58a6ff', '#1f6feb', '#0d419d'],
    green: ['#3fb950', '#238636', '#196127'],
    red: ['#f85149', '#da3633', '#b62324'],
    yellow: ['#e3b341', '#d29922', '#bb8009'],
  },
  primaryColor: 'purple',
  primaryShade: 6,
  fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", "Inter", "Roboto", sans-serif',
  headings: {
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    fontWeight: '700',
  },
  components: {
    Button: {
      defaultProps: {
        radius: 'md',
      },
      styles: {
        root: {
          fontWeight: 600,
        },
      },
    },
    Card: {
      styles: {
        root: {
          backgroundColor: '#161b22',
          border: '1px solid #30363d',
        },
      },
    },
    Input: {
      styles: {
        input: {
          backgroundColor: '#0d1117',
          borderColor: '#30363d',
          '&:focus': {
            borderColor: '#58a6ff',
          },
        },
      },
    },
  },
});

// Main App Component
const App: React.FC = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [isMaintenance, setIsMaintenance] = useState(false);
  const [servicesStatus, setServicesStatus] = useState<Record<string, any>>({});

  // Initialize application
  useEffect(() => {
    const initializeApp = async () => {
      try {
        // Check maintenance mode from Cloudflare KV
        const maintenanceRes = await fetch('/api/system/status');
        const systemStatus = await maintenanceRes.json();
        
        if (systemStatus.maintenanceMode) {
          setIsMaintenance(true);
          setIsLoading(false);
          return;
        }

        // Initialize services
        const servicesResponse = await fetch('/api/services/initialize', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        });

        if (servicesResponse.ok) {
          const servicesData = await servicesResponse.json();
          setServicesStatus(servicesData);
        }

        // Load user session if exists
        const token = localStorage.getItem('auth_token');
        if (token) {
          // Validate token with backend
          await fetch('/api/auth/validate', {
            headers: { Authorization: `Bearer ${token}` },
          });
        }

        setIsLoading(false);
      } catch (error) {
        console.error('Failed to initialize app:', error);
        Sentry.captureException(error);
        setIsLoading(false);
      }
    };

    initializeApp();

    // Setup service worker for PWA
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/sw.js').catch(console.error);
    }

    // Setup real-time notifications via WebSocket
    const setupWebSocket = () => {
      const ws = new WebSocket(import.meta.env.VITE_WS_URL || 'wss://ws.nawthtech.com');
      
      ws.onopen = () => {
        console.log('WebSocket connected');
        ws.send(JSON.stringify({ type: 'subscribe', channels: ['notifications', 'orders'] }));
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleWebSocketMessage(data);
      };

      return ws;
    };

    const ws = setupWebSocket();

    // Cleanup
    return () => {
      ws.close();
    };
  }, []);

  const handleWebSocketMessage = (data: any) => {
    switch (data.type) {
      case 'notification':
        // Handle real-time notification
        break;
      case 'order_update':
        // Handle order status update
        break;
      case 'ai_progress':
        // Handle AI processing progress
        break;
    }
  };

  if (isLoading) {
    return <LoadingScreen />;
  }

  if (isMaintenance) {
    return <MaintenanceMode />;
  }

  return (
    <Sentry.ErrorBoundary fallback={ErrorFallback}>
      <HelmetProvider>
        <QueryClientProvider client={queryClient}>
          <MantineProvider theme={theme} defaultColorScheme="dark">
            <ModalsProvider>
              <Notifications position="top-right" limit={5} />
              <Router>
                <AuthProvider>
                  <CartProvider>
                    <AIProvider>
                      <NotificationProvider>
                        <ThemeProvider>
                          <ServicesProvider>
                            <ErrorBoundary FallbackComponent={ErrorFallback}>
                              <div className="app-container">
                                <Header />
                                <div className="app-content">
                                  <Sidebar />
                                  <main className="main-content">
                                    <Navbar />
                                    <Suspense fallback={<LoadingOverlay visible />}>
                                      <Routes>
                                        {/* Public Routes */}
                                        <Route path="/" element={<HomePage />} />
                                        <Route path="/store" element={<StorePage />} />
                                        <Route path="/store/service/:id" element={<ServiceDetailPage />} />
                                        <Route path="/cart" element={<CartPage />} />
                                        <Route path="/login" element={<LoginPage />} />
                                        <Route path="/register" element={<RegisterPage />} />
                                        <Route path="/forgot-password" element={<ForgotPasswordPage />} />

                                        {/* Protected Routes */}
                                        <Route path="/checkout" element={
                                          <ProtectedRoute>
                                            <CheckoutPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/dashboard" element={
                                          <ProtectedRoute>
                                            <DashboardPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/ai-studio" element={
                                          <ProtectedRoute>
                                            <AIStudioPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/content-studio" element={
                                          <ProtectedRoute>
                                            <ContentStudioPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/video-studio" element={
                                          <ProtectedRoute>
                                            <VideoStudioPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/analytics" element={
                                          <ProtectedRoute>
                                            <AnalyticsPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/strategies" element={
                                          <ProtectedRoute>
                                            <StrategyPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/orders" element={
                                          <ProtectedRoute>
                                            <OrdersPage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/profile" element={
                                          <ProtectedRoute>
                                            <ProfilePage />
                                          </ProtectedRoute>
                                        } />
                                        <Route path="/settings" element={
                                          <ProtectedRoute>
                                            <SettingsPage />
                                          </ProtectedRoute>
                                        } />

                                        {/* Admin Routes */}
                                        <Route path="/admin/*" element={
                                          <AdminRoute>
                                            <AdminDashboard />
                                          </AdminRoute>
                                        } />

                                        {/* Catch-all route */}
                                        <Route path="*" element={<Navigate to="/" replace />} />
                                      </Routes>
                                    </Suspense>
                                  </main>
                                </div>
                                <Footer />
                              </div>
                            </ErrorBoundary>
                          </ServicesProvider>
                        </ThemeProvider>
                      </NotificationProvider>
                    </AIProvider>
                  </CartProvider>
                </AuthProvider>
              </Router>
              <Toaster 
                position="bottom-right"
                toastOptions={{
                  style: {
                    background: '#161b22',
                    color: '#f0f6fc',
                    border: '1px solid #30363d',
                  },
                }}
              />
            </ModalsProvider>
          </MantineProvider>
          {import.meta.env.DEV && <ReactQueryDevtools initialIsOpen={false} />}
        </QueryClientProvider>
      </HelmetProvider>
    </Sentry.ErrorBoundary>
  );
};

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('auth_token');
  
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

// Admin Route Component
const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const token = localStorage.getItem('auth_token');
  const userRole = localStorage.getItem('user_role');
  
  if (!token || userRole !== 'admin') {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};

export default App;