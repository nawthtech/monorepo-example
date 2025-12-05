#!/bin/bash

# تعيين API key
export GEMINI_API_KEY="your-api-key-here"

# تشغيل test
cd backend
go run cmd/test-gemini/main.go