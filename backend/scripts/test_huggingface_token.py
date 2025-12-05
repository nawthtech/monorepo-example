#!/usr/bin/env python3
"""
اختبار Hugging Face Token لـ NawthTech
"""

import os
import requests
import json
from dotenv import load_dotenv

# تحميل environment variables
load_dotenv()

def test_token():
    """اختبار صلاحية Token"""
    token = os.getenv("HUGGINGFACE_TOKEN")
    
    if not token:
        print("❌ لم يتم العثور على HUGGINGFACE_TOKEN في .env")
        return False
    
    print(f"🔑 Token: {token[:10]}...")
    
    # اختبار API
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    # اختبار 1: التحقق من صلاحية Token
    print("\n🔍 اختبار صلاحية Token...")
    try:
        response = requests.get(
            "https://huggingface.co/api/whoami",
            headers=headers,
            timeout=10
        )
        
        if response.status_code == 200:
            user_info = response.json()
            print(f"✅ Token صالح")
            print(f"👤 المستخدم: {user_info.get('name', 'غير معروف')}")
            print(f"📧 البريد: {user_info.get('email', 'غير معروف')}")
            print(f"🏢 المنظمة: {user_info.get('orgs', [])}")
        else:
            print(f"❌ Token غير صالح: {response.status_code}")
            return False
            
    except Exception as e:
        print(f"❌ خطأ في الاتصال: {e}")
        return False
    
    # اختبار 2: التحقق من أذونات النماذج
    print("\n🔍 اختبار أذونات النماذج...")
    models_to_test = [
        "google/flan-t5-xl",
        "stabilityai/stable-diffusion-xl-base-1.0",
        "openai/whisper-large-v3",
    ]
    
    for model in models_to_test:
        try:
            response = requests.get(
                f"https://huggingface.co/api/models/{model}",
                headers=headers,
                timeout=10
            )
            
            if response.status_code == 200:
                print(f"✅ يمكن الوصول إلى: {model}")
            else:
                print(f"❌ لا يمكن الوصول إلى: {model} ({response.status_code})")
                
        except Exception as e:
            print(f"⚠️ خطأ في النموذج {model}: {e}")
    
    # اختبار 3: اختبار Inference API
    print("\n🔍 اختبار Inference API...")
    test_payload = {
        "inputs": "مرحباً، هذا اختبار من NawthTech",
        "parameters": {
            "max_new_tokens": 50
        }
    }
    
    try:
        response = requests.post(
            "https://api-inference.huggingface.co/models/google/flan-t5-xl",
            headers=headers,
            json=test_payload,
            timeout=30
        )
        
        if response.status_code == 200:
            result = response.json()
            print(f"✅ Inference يعمل: {result}")
        elif response.status_code == 503:
            print("⚠️ النموذج قيد التحميل، جرب لاحقاً")
        else:
            print(f"❌ Inference فشل: {response.status_code}")
            print(f"الرد: {response.text}")
            
    except Exception as e:
        print(f"⚠️ خطأ في Inference: {e}")
    
    return True

def check_rate_limits():
    """التحقق من حدود الاستخدام"""
    print("\n📊 التحقق من حدود الاستخدام...")
    
    token = os.getenv("HUGGINGFACE_TOKEN")
    headers = {"Authorization": f"Bearer {token}"}
    
    try:
        response = requests.get(
            "https://huggingface.co/api/billing/usage",
            headers=headers,
            timeout=10
        )
        
        if response.status_code == 200:
            usage = response.json()
            print("✅ معلومات الاستخدام:")
            print(json.dumps(usage, indent=2, ensure_ascii=False))
        else:
            print("❌ لا يمكن الحصول على معلومات الاستخدام")
            
    except Exception as e:
        print(f"⚠️ خطأ: {e}")

def main():
    print("=" * 60)
    print("🤖 NawthTech Hugging Face Token Tester")
    print("=" * 60)
    
    if test_token():
        check_rate_limits()
        print("\n🎉 جميع الاختبارات اكتملت بنجاح!")
        print("\n📝 ملاحظات:")
        print("1. تأكد من حفظ Token في .env")
        print("2. الحدود: 30 طلب/دقيقة مجاناً")
        print("3. بعض النماذج تحتاج تحميل عند أول طلب")
    else:
        print("\n❌ هناك مشاكل في Token، راجع الإعدادات")

if __name__ == "__main__":
    main()