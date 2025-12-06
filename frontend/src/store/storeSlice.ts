import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

interface Service {
  id: string;
  name: string;
  platform: string;
  description: string;
  price: number;
  currency: string;
  quantity: number;
  features: string[];
  stats: {
    orders: number;
    success_rate: number;
    delivery_time: string;
  };
}

interface CartItem {
  serviceId: string;
  quantity: number;
  price: number;
}

interface StoreState {
  services: Service[];
  cart: CartItem[];
  selectedCategory: string;
  searchQuery: string;
  loading: boolean;
  error: string | null;
  featuredServices: string[];
}

const initialState: StoreState = {
  services: [],
  cart: JSON.parse(localStorage.getItem('cart') || '[]'),
  selectedCategory: 'all',
  searchQuery: '',
  loading: false,
  error: null,
  featuredServices: [],
};

const storeSlice = createSlice({
  name: 'store',
  initialState,
  reducers: {
    // تحميل الخدمات
    loadServicesStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    
    // نجاح تحميل الخدمات
    loadServicesSuccess: (state, action: PayloadAction<Service[]>) => {
      state.loading = false;
      state.services = action.payload;
      state.featuredServices = action.payload
        .filter(service => service.stats.orders > 1000)
        .map(service => service.id);
    },
    
    // فشل تحميل الخدمات
    loadServicesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    
    // إضافة إلى السلة
    addToCart: (state, action: PayloadAction<CartItem>) => {
      const existingItem = state.cart.find(item => item.serviceId === action.payload.serviceId);
      
      if (existingItem) {
        existingItem.quantity += action.payload.quantity;
      } else {
        state.cart.push(action.payload);
      }
      
      localStorage.setItem('cart', JSON.stringify(state.cart));
    },
    
    // إزالة من السلة
    removeFromCart: (state, action: PayloadAction<string>) => {
      state.cart = state.cart.filter(item => item.serviceId !== action.payload);
      localStorage.setItem('cart', JSON.stringify(state.cart));
    },
    
    // تحديث الكمية
    updateCartQuantity: (state, action: PayloadAction<{ serviceId: string; quantity: number }>) => {
      const item = state.cart.find(item => item.serviceId === action.payload.serviceId);
      if (item) {
        item.quantity = action.payload.quantity;
      }
      localStorage.setItem('cart', JSON.stringify(state.cart));
    },
    
    // تفريغ السلة
    clearCart: (state) => {
      state.cart = [];
      localStorage.removeItem('cart');
    },
    
    // تغيير التصنيف
    changeCategory: (state, action: PayloadAction<string>) => {
      state.selectedCategory = action.payload;
    },
    
    // تغيير بحث
    changeSearch: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
    
    // مسح الأخطاء
    clearStoreError: (state) => {
      state.error = null;
    },
  },
});

export const {
  loadServicesStart,
  loadServicesSuccess,
  loadServicesFailure,
  addToCart,
  removeFromCart,
  updateCartQuantity,
  clearCart,
  changeCategory,
  changeSearch,
  clearStoreError,
} = storeSlice.actions;

export default storeSlice.reducer;