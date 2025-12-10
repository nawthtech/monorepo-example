import React, { createContext, useContext, useState, useEffect } from 'react';

interface ServiceStatus {
  name: string;
  status: 'healthy' | 'degraded' | 'unhealthy';
  message?: string;
}

interface ServicesContextType {
  services: ServiceStatus[];
  isLoading: boolean;
  refreshServices: () => Promise<void>;
  getService: (name: string) => ServiceStatus | undefined;
}

const ServicesContext = createContext<ServicesContextType | undefined>(undefined);

export const ServicesProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [services, setServices] = useState<ServiceStatus[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const loadServices = async () => {
    setIsLoading(true);
    try {
      const response = await fetch('/api/services/status');
      const data = await response.json();
      setServices(data);
    } catch (error) {
      console.error('Failed to load services:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const refreshServices = async () => {
    await loadServices();
  };

  const getService = (name: string) => {
    return services.find(s => s.name === name);
  };

  useEffect(() => {
    loadServices();
    
    // Refresh services every 30 seconds
    const interval = setInterval(loadServices, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <ServicesContext.Provider value={{ services, isLoading, refreshServices, getService }}>
      {children}
    </ServicesContext.Provider>
  );
};

export const useServices = () => {
  const context = useContext(ServicesContext);
  if (!context) {
    throw new Error('useServices must be used within ServicesProvider');
  }
  return context;
};