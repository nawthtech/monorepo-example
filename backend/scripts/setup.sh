#!/bin/bash

echo "🚀 بدء إعداد مشروع NawthTech Backend مع Cloudflare D1..."

# التحقق من تثبيت Go
if ! command -v go &> /dev/null; then
    echo "❌ Go غير مثبت. يرجى تثبيت Go 1.25 أو أحدث"
    echo "   قم بزيارة: https://go.dev/dl/"
    exit 1
fi

echo "✅ Go مثبت: $(go version)"

# التحقق من إصدار Go
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+')
if [[ "$GO_VERSION" < "go1.25" ]]; then
    echo "⚠️  إصدار Go $GO_VERSION - يوصى باستخدام Go 1.21 أو أحدث"
fi

# الانتقال إلى مجلد المشروع
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

echo "📁 المجلد الحالي: $(pwd)"

# تنظيف الاعتمادات القديمة (بحذر)
if [ -f "go.mod" ]; then
    echo "📋 go.mod موجود بالفعل، حفظ نسخة احتياطية..."
    cp go.mod go.mod.backup
fi

if [ -f "go.sum" ]; then
    echo "📋 go.sum موجود بالفعل، حفظ نسخة احتياطية..."
    cp go.sum go.sum.backup
fi

# تهيئة المشروع
echo "📦 تهيئة مشروع Go..."
go mod init github.com/nawthtech/nawthtech/backend

# إضافة الاعتمادات الأساسية للـ Backend مع Cloudflare D1
echo "📥 إضافة الاعتمادات الأساسية..."

# 1. Framework الأساسي (Gin)
echo "   🕸️  إضافة Gin Web Framework..."
go get github.com/gin-gonic/gin@v1.9.1

# 2. قاعدة البيانات (Cloudflare D1 - SQLite)
echo "   🗄️  إضافة SQLite driver لـ Cloudflare D1..."
go get github.com/mattn/go-sqlite3@v1.14.19

# 3. Cloudinary للرفع
echo "   ☁️  إضافة Cloudinary SDK..."
go get github.com/cloudinary/cloudinary-go/v2@v2.5.1

# 4. CLI Tool
echo "   🛠️  إضافة CLI Tool..."
go get github.com/urfave/cli/v2@v2.25.7

# 5. البريد الإلكتروني
echo "   📧 إضافة مكتبة البريد الإلكتروني..."
go get gopkg.in/gomail.v2@v2.0.0
go get github.com/resend/resend-go/v2@latest

# 6. Environment Variables
echo "   🔧 إضافة مكتبة متغيرات البيئة..."
go get github.com/joho/godotenv@v1.5.1

# 7. JWT Authentication
echo "   🔐 إضافة JWT للمصادقة..."
go get github.com/golang-jwt/jwt/v5@v5.0.0

# 8. Testing & Assertions
echo "   🧪 إضافة مكتبات الاختبار..."
go get github.com/stretchr/testify@v1.8.4

# 9. Logging
echo "   📝 إضافة مكتبة التسجيل (Logging)..."
go get go.uber.org/zap@v1.26.0

# 10. Configuration
echo "   ⚙️  إضافة مكتبة التكوين..."
go get github.com/spf13/viper@v1.17.0

# 11. Validation
echo "   ✓ إضافة مكتبة التحقق..."
go get github.com/go-playground/validator/v10@v10.15.5

# 12. HTTP Client
echo "   🌐 إضافة HTTP Client..."
go get github.com/go-resty/resty/v2@v2.10.0

# 13. Slack Integration
echo "   💬 إضافة تكامل Slack..."
go get github.com/slack-go/slack@v0.12.3

# 14. Stripe Payments
echo "   💳 إضافة Stripe للدفع..."
go get github.com/stripe/stripe-go/v76@v76.0.0

# 15. UUID Generation
echo "   🆔 إضافة مكتبة UUID..."
go get github.com/google/uuid@v1.4.0

# 16. CORS Middleware
echo "   🌍 إضافة CORS Middleware..."
go get github.com/gin-contrib/cors@v1.5.0

# 17. Compression
echo "   🗜️  إضافة ضغط GZIP..."
go get github.com/gin-contrib/gzip@v1.0.0

# 18. Rate Limiting
echo "   ⏱️  إضافة Rate Limiting..."
go get golang.org/x/time/rate@latest

# 19. Cryptography
echo "   🔒 إضافة مكتبات التشفير..."
go get golang.org/x/crypto@latest

# 20. Cloudflare Workers (اختياري)
echo "   ⚡ إضافة Cloudflare Workers API..."
go get github.com/cloudflare/cloudflare-go@v0.86.0

# تنظيف وتحديث الاعتمادات
echo "🔧 تنظيف وتحديث الاعتمادات..."
go mod tidy

# تحميل الاعتمادات
echo "📥 تحميل الاعتمادات..."
go mod download

# التحقق من الصحة
echo "🔍 التحقق من صحة الاعتمادات..."
go mod verify

# اختبار البناء
echo "🏗️ اختبار بناء المشروع..."
if go build -o /tmp/nawthtech-backend ./cmd/server; then
    echo "✅ البناء نجح!"
    rm -f /tmp/nawthtech-backend
else
    echo "❌ فشل البناء!"
    echo "⚠️  تحقق من الأخطاء أعلاه"
    exit 1
fi

# إنشاء الهيكل الأساسي للمجلدات
echo "📁 إنشاء هيكل المجلدات..."
mkdir -p cmd/server
mkdir -p internal/{config,db,handlers,services,middleware,utils,models,logger,routes}
mkdir -p api/{v1,v2}
mkdir -p scripts
mkdir -p data
mkdir -p uploads
mkdir -p logs
mkdir -p tests

# إنشاء ملف .env مثال
echo "📄 إنشاء ملف .env.example..."
cat > .env.example << 'EOF'
# ==================== الأساسية ====================
APP_NAME=nawthtech
APP_VERSION=1.0.0
ENVIRONMENT=development
PORT=8080
DEBUG=true

# ==================== URLs ====================
API_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
WORKER_API_URL=https://api.nawthtech.com
WORKER_API_KEY=""

# ==================== الأمان ====================
JWT_SECRET=""
REFRESH_SECRET=""
ENCRYPTION_KEY=""
API_KEY=""

# ==================== قاعدة البيانات (Cloudflare D1) ====================
DB_DRIVER=sqlite3
DATABASE_URL=database.nawthtech.com

# ==================== CORS ====================
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080,http://localhost:3001,https://nawthtech.com,www.nawthtech.com,https://www.nawthtech.com,https://api.nawthtech.com,https://database.nawthtech.com
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With,X-API-Key
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=300

# ==================== البريد ====================
EMAIL_ENABLED=false
EMAIL_PROVIDER=resend
RESEND_API_KEY=""
EMAIL_FROM=support@nawthtech.com
EMAIL_FROM_NAME=NawthTech

# ==================== الرفع ====================
UPLOAD_MAX_SIZE=10485760  # 10MB
UPLOAD_PATH=./uploads
UPLOAD_ALLOWED_TYPES=image/jpeg,image/png,image/gif,image/webp,application/pdf
CLOUDINARY_URL=""

# ==================== التخزين المؤقت ====================
CACHE_ENABLED=true
CACHE_TYPE=memory
REDIS_URL=redis://localhost:6379
CACHE_TTL=5m

# ==================== الأمان ====================
RATE_LIMIT=100
RATE_WINDOW=1m

# ==================== الخدمات ====================
SLACK_TOKEN=""
SLACK_CHANNEL=general
SLACK_APP_NAME=nawthtech

STRIPE_SECRET_KEY=""
STRIPE_WEBHOOK_SECRET=""
STRIPE_PUBLISHABLE_KEY=""

CLOUDINARY_CLOUD_NAME=""
CLOUDINARY_API_KEY=""
CLOUDINARY_API_SECRET=""

# ==================== الذكاء الاصطناعي ====================
OPENAI_API_KEY=""
OPENAI_MODEL=gpt-4-turbo-preview
GEMINI_API_KEY=""
EOF

echo "📄 إنشاء ملف .gitignore..."
cat > .gitignore << 'EOF'
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
nawthtech-backend

# IDE
.vscode/
.idea/
*.swp
*.swo

# Environment
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Logs
*.log
logs/

# Database
*.db
*.db-journal
data/

# Uploads
uploads/

# Coverage
coverage.out
coverage.html

# OS
.DS_Store
Thumbs.db

# Backup files
*.backup
*.bak

# Temp files
tmp/
temp/

# Go
vendor/
EOF

echo "📄 إنشاء ملف Makefile للمهام الشائعة..."
cat > Makefile << 'EOF'
.PHONY: build run test clean deps migrate dev

# بناء المشروع
build:
	@echo "🏗️  بناء المشروع..."
	go build -o nawthtech-backend ./cmd/server

# تشغيل المشروع
run:
	@echo "🚀 تشغيل الخادم..."
	go run ./cmd/server

# تشغيل في وضع التطوير
dev:
	@echo "🔧 تشغيل في وضع التطوير..."
	ENVIRONMENT=development go run ./cmd/server

# تشغيل جميع الاختبارات
test:
	@echo "🧪 تشغيل الاختبارات..."
	go test ./... -v

# تشغيل اختبارات سريعة
test-short:
	@echo "⚡ اختبارات سريعة..."
	go test ./... -short

# تنظيف الملفات المبنية
clean:
	@echo "🧹 تنظيف الملفات..."
	rm -f nawthtech-backend
	rm -f coverage.out
	rm -rf dist/

# تحديث الاعتمادات
deps:
	@echo "📦 تحديث الاعتمادات..."
	go mod tidy
	go mod download

# تشغيل عمليات الترحيل
migrate:
	@echo "🔄 تشغيل عمليات الترحيل..."
	go run ./scripts/migrate.go

# فحص الجودة
lint:
	@echo "🔍 فحص الكود..."
	gofmt -d .
	golangci-lint run

# نسخة احتياطية لقاعدة البيانات
backup:
	@echo "💾 نسخ احتياطي لقاعدة البيانات..."
	cp data/nawthtech.db data/nawthtech.db.backup.$(shell date +%Y%m%d_%H%M%S)

# استعادة نسخة احتياطية
restore:
	@echo "🔄 استعادة قاعدة البيانات..."
	@if [ -f "data/nawthtech.db.backup" ]; then \
		cp data/nawthtech.db.backup data/nawthtech.db; \
		echo "✅ تمت الاستعادة"; \
	else \
		echo "❌ لا توجد نسخة احتياطية"; \
	fi

# عرض المساعدة
help:
	@echo "أوامر Makefile المتاحة:"
	@echo "  build     - بناء المشروع"
	@echo "  run       - تشغيل المشروع"
	@echo "  dev       - تشغيل في وضع التطوير"
	@echo "  test      - تشغيل جميع الاختبارات"
	@echo "  test-short - اختبارات سريعة"
	@echo "  clean     - تنظيف الملفات المبنية"
	@echo "  deps      - تحديث الاعتمادات"
	@echo "  migrate   - تشغيل عمليات الترحيل"
	@echo "  lint      - فحص جودة الكود"
	@echo "  backup    - نسخ احتياطي للقاعدة"
	@echo "  restore   - استعادة القاعدة"
	@echo "  help      - عرض هذه الرسالة"
EOF

echo "📄 إنشاء ملف README للـ Backend..."
cat > backend/README.md << 'EOF'
# NawthTech Backend

Backend API لمشروع NawthTech مبني بـ Go و Cloudflare D1.

## 🏗️ البنية التقنية

### المكونات الأساسية
- **Go 1.25+** - لغة البرمجة
- **Gin** - Web Framework
- **Cloudflare D1** - قاعدة البيانات (SQLite)
- **Cloudinary** - تخزين الملفات
- **JWT** - المصادقة

### الهيكل
