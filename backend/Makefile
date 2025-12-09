# Makefile for NawthTech Backend

.PHONY: help build run test clean deps migrate

# المساعدة
help:
	@echo "أوامر NawthTech Backend المتاحة:"
	@echo ""
	@echo "  build     - بناء التطبيق"
	@echo "  run       - تشغيل التطبيق (التطوير)"
	@echo "  test      - تشغيل الاختبارات"
	@echo "  clean     - تنظيف الملفات المبنية"
	@echo "  deps      - تحديث الاعتمادات"
	@echo "  migrate   - تشغيل ترحيل قاعدة البيانات"
	@echo "  docker    - بناء صورة Docker"
	@echo ""

# بناء التطبيق
build:
	@echo "🏗️  بناء تطبيق نوذ تك..."
	go build -o bin/nawthtech main.go

# تشغيل التطبيق (التطوير)
run:
	@echo "🚀 تشغيل خادم التطوير..."
	go run main.go server --env development

# Email Management
setup-email:
	@echo "🚀 Setting up email routing..."
	@go run cmd/email-setup/main.go -action=setup

test-email:
	@if [ -z "$(email)" ]; then \
		echo "❌ Usage: make test-email email=test@example.com"; \
		exit 1; \
	fi
	@go run cmd/email-setup/main.go -action=test -email=$(email)

add-email:
	@if [ -z "$(email)" ]; then \
		echo "❌ Usage: make add-email email=user@example.com"; \
		exit 1; \
	fi
	@go run cmd/email-setup/main.go -action=add-email -email=$(email)

list-emails:
	@go run cmd/email-setup/main.go -action=list

# Run all setup
setup-all: setup-email
	@echo "✅ All services configured"

# تشغيل الاختبارات
test:
	@echo "🧪 تشغيل الاختبارات..."
	go test -v ./...

# تنظيف الملفات المبنية
clean:
	@echo "🧹 تنظيف الملفات المبنية..."
	rm -rf bin/
	rm -rf coverage.txt

# تحديث الاعتمادات
deps:
	@echo "📦 تحديث اعتمادات Go..."
	go mod tidy
	go mod download

# ترحيل قاعدة البيانات
migrate:
	@echo "🗄️  تشغيل ترحيل قاعدة البيانات..."
	go run main.go migrate --direction up

# بناء صورة Docker
docker:
	@echo "🐳 بناء صورة Docker..."
	docker build -t nawthtech/backend:latest .

# فحص الصحة
health:
	@echo "🔍 فحص صحة النظام..."
	go run main.go health

# عرض الإصدار
version:
	@echo "📋 عرض معلومات الإصدار..."
	go run main.go version

# التنسيق التلقائي
fmt:
	@echo "🎨 تنسيق الكود..."
	go fmt ./...

# التحقق من الكود
vet:
	@echo "🔍 التحقق من الكود..."
	go vet ./...