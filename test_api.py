#!/usr/bin/env python3
"""
EVA-Mind-FZPN - API Test Script
Testa todos os endpoints da API
"""

import requests
import json
from datetime import datetime

# Config
API_BASE_URL = "http://localhost:8000"
GO_SERVICE_URL = "http://localhost:8081"

# Colors for terminal
GREEN = '\033[92m'
RED = '\033[91m'
YELLOW = '\033[93m'
BLUE = '\033[94m'
RESET = '\033[0m'

def print_success(msg):
    print(f"{GREEN}✓{RESET} {msg}")

def print_error(msg):
    print(f"{RED}✗{RESET} {msg}")

def print_info(msg):
    print(f"{BLUE}ℹ{RESET} {msg}")

def print_warning(msg):
    print(f"{YELLOW}⚠{RESET} {msg}")

def test_service_health():
    """Test if services are running"""
    print("\n" + "=" * 60)
    print("🏥 HEALTH CHECKS")
    print("=" * 60)

    # Test Go service
    try:
        response = requests.get(f"{GO_SERVICE_URL}/health", timeout=2)
        if response.status_code == 200:
            print_success("Go Integration Service: OK")
        else:
            print_error(f"Go Integration Service: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Go Integration Service: {str(e)}")
        return False

    # Test Python API
    try:
        response = requests.get(f"{API_BASE_URL}/health", timeout=2)
        if response.status_code == 200:
            print_success("Python API Server: OK")
            data = response.json()
            print(f"  Status: {data['status']}")
        else:
            print_error(f"Python API Server: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Python API Server: {str(e)}")
        return False

    return True

def test_oauth_flow():
    """Test OAuth2 authentication"""
    print("\n" + "=" * 60)
    print("🔐 OAUTH2 AUTHENTICATION")
    print("=" * 60)

    # Use test client credentials (created by create_test_client.py)
    client_id = "eva_test_client"
    client_secret = "test_secret_123"

    print_info(f"Authenticating with client_id: {client_id}")

    try:
        response = requests.post(
            f"{API_BASE_URL}/oauth/token",
            data={
                "username": client_id,  # OAuth2PasswordRequestForm usa username
                "password": client_secret,
                "grant_type": "password"
            }
        )

        if response.status_code == 200:
            token_data = response.json()
            access_token = token_data['access_token']
            print_success("OAuth2 token obtido com sucesso")
            print(f"  Token: {access_token[:50]}...")
            print(f"  Expires in: {token_data['expires_in']} seconds")
            return access_token
        else:
            print_error(f"Falha na autenticação: HTTP {response.status_code}")
            print(f"  Response: {response.text}")
            print_warning("Execute create_test_client.py primeiro para criar credenciais")
            return None

    except Exception as e:
        print_error(f"Erro na autenticação: {str(e)}")
        return None

def test_patient_endpoints(token):
    """Test patient endpoints"""
    print("\n" + "=" * 60)
    print("👤 PATIENT ENDPOINTS")
    print("=" * 60)

    headers = {"Authorization": f"Bearer {token}"}

    # Test list patients
    print_info("Testing GET /api/v1/patients")
    try:
        response = requests.get(f"{API_BASE_URL}/api/v1/patients", headers=headers)
        if response.status_code == 200:
            data = response.json()
            print_success(f"Lista de pacientes obtida: {data['total_count']} total")
            if data['data']:
                print(f"  Primeiro paciente: {data['data'][0]['name']}")
        else:
            print_error(f"Falha: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Erro: {str(e)}")

    # Test get patient by ID
    print_info("Testing GET /api/v1/patients/1")
    try:
        response = requests.get(f"{API_BASE_URL}/api/v1/patients/1", headers=headers)
        if response.status_code == 200:
            patient = response.json()
            print_success(f"Paciente obtido: {patient['name']}")
            print(f"  Age: {patient['age']}")
            print(f"  Gender: {patient['gender']}")
        elif response.status_code == 404:
            print_warning("Paciente #1 não encontrado")
        else:
            print_error(f"Falha: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Erro: {str(e)}")

def test_fhir_endpoints(token):
    """Test FHIR endpoints"""
    print("\n" + "=" * 60)
    print("🏥 FHIR ENDPOINTS")
    print("=" * 60)

    headers = {"Authorization": f"Bearer {token}"}

    # Test FHIR patient
    print_info("Testing GET /api/v1/fhir/patients/1")
    try:
        response = requests.get(f"{API_BASE_URL}/api/v1/fhir/patients/1", headers=headers)
        if response.status_code == 200:
            fhir_patient = response.json()
            print_success(f"FHIR Patient obtido")
            print(f"  Resource Type: {fhir_patient.get('resourceType')}")
            print(f"  ID: {fhir_patient.get('id')}")
        else:
            print_error(f"Falha: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Erro: {str(e)}")

    # Test FHIR bundle
    print_info("Testing GET /api/v1/fhir/bundle/1")
    try:
        response = requests.get(f"{API_BASE_URL}/api/v1/fhir/bundle/1", headers=headers)
        if response.status_code == 200:
            bundle = response.json()
            print_success(f"FHIR Bundle obtido")
            print(f"  Type: {bundle.get('type')}")
            print(f"  Entries: {len(bundle.get('entry', []))}")
        else:
            print_error(f"Falha: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Erro: {str(e)}")

def test_export_endpoints(token):
    """Test export endpoints"""
    print("\n" + "=" * 60)
    print("📦 EXPORT ENDPOINTS")
    print("=" * 60)

    headers = {"Authorization": f"Bearer {token}"}

    # Test LGPD export
    print_info("Testing GET /api/v1/export/lgpd/1")
    try:
        response = requests.get(f"{API_BASE_URL}/api/v1/export/lgpd/1", headers=headers)
        if response.status_code == 200:
            export = response.json()
            print_success(f"LGPD Export obtido")
            print(f"  Export date: {export.get('export_date')}")
            print(f"  Patient ID: {export.get('patient_id')}")
        else:
            print_error(f"Falha: HTTP {response.status_code}")
    except Exception as e:
        print_error(f"Erro: {str(e)}")

def test_rate_limiting(token):
    """Test rate limiting"""
    print("\n" + "=" * 60)
    print("⏱️  RATE LIMITING")
    print("=" * 60)

    headers = {"Authorization": f"Bearer {token}"}

    print_info("Fazendo 65 requests rápidos (limite é 60/min)...")

    success_count = 0
    rate_limited_count = 0

    for i in range(1, 66):
        try:
            response = requests.get(f"{API_BASE_URL}/api/v1/patients", headers=headers, timeout=1)
            if response.status_code == 200:
                success_count += 1
            elif response.status_code == 429:
                rate_limited_count += 1
                if rate_limited_count == 1:
                    print_warning(f"Rate limit atingido no request #{i}")
        except Exception:
            pass

    print_success(f"Requests bem-sucedidos: {success_count}")
    if rate_limited_count > 0:
        print_success(f"Rate limiting funcionando: {rate_limited_count} requests bloqueados")
    else:
        print_warning("Rate limiting não foi acionado")

def main():
    """Run all tests"""
    print("\n" + "=" * 70)
    print(" " * 15 + "🧪 EVA-MIND API TEST SUITE")
    print("=" * 70)
    print(f"Timestamp: {datetime.now()}")

    # Test 1: Health checks
    if not test_service_health():
        print_error("\n❌ Serviços não estão rodando!")
        print_info("Execute start_services.bat primeiro")
        return

    # Test 2: OAuth2
    token = test_oauth_flow()
    if not token:
        print_error("\n❌ Falha na autenticação!")
        print_info("Execute create_test_client.py para criar credenciais de teste")
        return

    # Test 3: Patient endpoints
    test_patient_endpoints(token)

    # Test 4: FHIR endpoints
    test_fhir_endpoints(token)

    # Test 5: Export endpoints
    test_export_endpoints(token)

    # Test 6: Rate limiting
    test_rate_limiting(token)

    # Summary
    print("\n" + "=" * 70)
    print(" " * 20 + "✅ TESTES CONCLUÍDOS")
    print("=" * 70)
    print("\n📚 Para mais informações, acesse: http://localhost:8000/docs")
    print()

if __name__ == "__main__":
    main()
