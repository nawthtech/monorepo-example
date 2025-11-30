#!/bin/bash

echo "🚀 بدء إعداد مشروع NawthTech Backend..."

# التحقق من تثبيت Go
if ! command -v go &> /dev/null; then
    echo "❌ Go غير مثبت. يرجى تثبيت Go 1.21 أو أحدث"
    exit 1
fi

echo "✅ Go مثبت: $(go version)"

# الانتقال إلى مجلد المشروع
cd "$(dirname "$0")/.."

# تنظيف الاعتمادات القديمة
echo "🧹 تنظيف الاعتمادات القديمة..."
rm -f go.mod go.sum

# تهيئة المشروع
echo "📦 تهيئة مشروع Go..."
go mod init github.com/nawthtech/nawthtech/backend

# إضافة الاعتمادات الأساسية
echo "📥 إضافة الاعتمادات الأساسية..."
go get github.com/gin-gonic/gin@v1.9.1
go get go.mongodb.org/mongo-driver@v1.12.1
go get github.com/cloudinary/cloudinary-go/v2@v2.5.1
go get github.com/urfave/cli/v2@v2.25.7
go get gopkg.in/gomail.v2@v2.0.0
go get github.com/joho/godotenv@v1.5.1
go get github.com/golang-jwt/jwt/v5@v5.0.0

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
if go build -o /tmp/nawthtech-backend ./...; then
    echo "✅ البناء نجح!"
    rm -f /tmp/nawthtech-backend
else
    echo "❌ فشل البناء!"
    exit 1
fi

echo "🎉 تم إعداد مشروع NawthTech Backend بنجاح!"
echo ""
echo "📋 الاعتمادات المثبتة:"
go list -m all | head -20