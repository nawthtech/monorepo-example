import React, { useState } from 'react';
import { ShoppingCart, User, Search, Menu } from 'lucide-react';

interface HeaderProps {
  cartCount: number;
  onCartClick: () => void;
  onProfileClick: () => void;
  onSearch: (query: string) => void;
}

export const Header: React.FC<HeaderProps> = ({
  cartCount,
  onCartClick,
  onProfileClick,
  onSearch
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchQuery);
  };

  return (
    <header className="store-header">
      {/* شعار NawthTech */}
      <div className="header-logo">
        <svg width="140" height="30" viewBox="0 0 160 40" className="logo-svg">
          <defs>
            <linearGradient id="headerGradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" style={{ stopColor: '#bc8cff', stopOpacity: 1 }} />
              <stop offset="50%" style={{ stopColor: '#7c3aed', stopOpacity: 1 }} />
              <stop offset="100%" style={{ stopColor: '#58a6ff', stopOpacity: 1 }} />
            </linearGradient>
          </defs>
          <rect x="5" y="5" width="30" height="30" rx="8" fill="url(#headerGradient)" opacity="0.9" />
          <text x="20" y="25" fontFamily="'Segoe UI', 'Inter', sans-serif" fontWeight="900" fontSize="16" textAnchor="middle" fill="#ffffff">NT</text>
          <text x="50" y="28" fontFamily="'Segoe UI', 'Inter', sans-serif" fontWeight="800" fontSize="24" fill="url(#headerGradient)">NawthTech</text>
          <circle cx="145" cy="15" r="6" fill="#3fb950" />
          <text x="145" y="18" fontFamily="'Segoe UI', 'Inter', sans-serif" fontWeight="900" fontSize="8" textAnchor="middle" fill="#ffffff">AI</text>
        </svg>
      </div>

      {/* زر القائمة للموبايل */}
      <button 
        className="mobile-menu-btn"
        onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
        aria-label="Toggle menu"
      >
        <Menu size={24} />
      </button>

      {/* البحث */}
      <form onSubmit={handleSearch} className="header-search">
        <input
          type="text"
          className="search-input"
          placeholder="ابحث في الخدمات..."
          value={searchQuery}
          onChange={(e) => {
            setSearchQuery(e.target.value);
            if (e.target.value === '') onSearch('');
          }}
          dir="rtl"
        />
        <span className="search-icon">
          <Search size={18} />
        </span>
      </form>

      {/* القائمة المتنقلة */}
      <div className={`mobile-menu ${mobileMenuOpen ? 'open' : ''}`}>
        <button className="action-btn mobile-cart-btn" onClick={onCartClick}>
          <ShoppingCart size={18} />
          <span>السلة</span>
          {cartCount > 0 && (
            <span className="cart-badge">{cartCount}</span>
          )}
        </button>
        
        <button className="action-btn" onClick={onProfileClick}>
          <User size={18} />
          <span>حسابي</span>
        </button>
      </div>

      {/* الأزرار - للديسكتوب */}
      <div className="header-actions">
        <button className="action-btn" onClick={onCartClick}>
          <ShoppingCart size={18} />
          <span>السلة</span>
          {cartCount > 0 && (
            <span className="cart-badge">{cartCount}</span>
          )}
        </button>
        
        <button className="action-btn" onClick={onProfileClick}>
          <User size={18} />
          <span>حسابي</span>
        </button>
      </div>
    </header>
  );
};