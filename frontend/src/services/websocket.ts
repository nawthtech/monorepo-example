class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout = 1000;
  private listeners: Map<string, ((data: any) => void)[]> = new Map();

  connect() {
    const wsUrl = import.meta.env.VITE_WS_URL || 'wss://ws.nawthtech.com';
    
    try {
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        this.authenticate();
      };

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.handleMessage(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        this.reconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (error) {
      console.error('Failed to connect WebSocket:', error);
    }
  }

  private authenticate() {
    const token = localStorage.getItem('auth_token');
    if (token && this.ws?.readyState === WebSocket.OPEN) {
      this.send({
        type: 'auth',
        token,
      });
    }
  }

  private reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);
      
      setTimeout(() => {
        this.connect();
      }, this.reconnectTimeout * this.reconnectAttempts);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  send(data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.warn('WebSocket is not connected');
    }
  }

  subscribe(channel: string) {
    this.send({
      type: 'subscribe',
      channel,
    });
  }

  unsubscribe(channel: string) {
    this.send({
      type: 'unsubscribe',
      channel,
    });
  }

  on(channel: string, callback: (data: any) => void) {
    if (!this.listeners.has(channel)) {
      this.listeners.set(channel, []);
    }
    this.listeners.get(channel)?.push(callback);
  }

  off(channel: string, callback: (data: any) => void) {
    const callbacks = this.listeners.get(channel);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  private handleMessage(data: any) {
    const { type, channel, payload } = data;
    
    if (channel && this.listeners.has(channel)) {
      this.listeners.get(channel)?.forEach(callback => {
        callback(payload);
      });
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.listeners.clear();
  }
}

export const websocketService = new WebSocketService();