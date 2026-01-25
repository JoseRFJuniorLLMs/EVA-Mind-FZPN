#!/usr/bin/env python3
"""
Script para criar um API client de teste
"""

import psycopg2
import bcrypt
import uuid
from datetime import datetime

# Configuração do database remoto
DB_CONFIG = {
    "host": "104.248.219.200",
    "port": 5432,
    "database": "eva-db",
    "user": "postgres",
    "password": "Debian23@"
}

def create_test_client():
    """Cria um API client de teste"""
    print(f"Conectando ao PostgreSQL em {DB_CONFIG['host']}...")

    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cur = conn.cursor()

        print(f"✓ Conectado com sucesso!\n")

        # Gerar credenciais
        client_id = "eva_test_client"
        client_secret = "test_secret_123"
        client_secret_hash = bcrypt.hashpw(client_secret.encode(), bcrypt.gensalt()).decode()

        # Criar API client
        cur.execute("""
            INSERT INTO api_clients (
                id,
                client_name,
                client_type,
                client_id,
                client_secret_hash,
                scopes,
                rate_limit_per_minute,
                rate_limit_per_hour,
                rate_limit_per_day,
                webhook_url,
                webhook_secret,
                is_active,
                is_approved,
                created_at,
                updated_at
            ) VALUES (
                %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s
            )
            ON CONFLICT (client_id) DO UPDATE SET
                client_secret_hash = EXCLUDED.client_secret_hash,
                updated_at = EXCLUDED.updated_at
            RETURNING id, client_name, client_id;
        """, (
            str(uuid.uuid4()),
            "EVA Test Client",
            "third_party",  # Valid values: web_app, mobile_app, hospital_system, research_platform, ehr_system, third_party
            client_id,
            client_secret_hash,
            ['read:patients', 'write:patients', 'read:assessments', 'write:assessments',
             'read:medications', 'write:medications', 'read:trajectories', 'export:data'],
            60,    # 60 requests per minute
            3600,  # 3600 requests per hour
            50000, # 50000 requests per day
            None,  # webhook_url (opcional)
            None,  # webhook_secret
            True,  # is_active
            True,  # is_approved
            datetime.now(),
            datetime.now()
        ))

        result = cur.fetchone()
        conn.commit()

        print(f"✅ API Client criado com sucesso!\n")
        print(f"📋 Credenciais do Cliente:")
        print(f"  Client ID:     {client_id}")
        print(f"  Client Secret: {client_secret}")
        print(f"  \n  ⚠️  SALVE ESTAS CREDENCIAIS! O secret não pode ser recuperado depois.\n")

        print(f"🔑 Scopes Disponíveis:")
        print(f"  - read:patients       (Ler dados de pacientes)")
        print(f"  - write:patients      (Criar/atualizar pacientes)")
        print(f"  - read:assessments    (Ler assessments)")
        print(f"  - write:assessments   (Criar assessments)")
        print(f"  - read:medications    (Ler medicações)")
        print(f"  - write:medications   (Gerenciar medicações)")
        print(f"  - read:trajectories   (Ler predições de trajetória)")
        print(f"  - export:data         (Exportar dados LGPD/FHIR)")

        print(f"\n⚡ Rate Limits:")
        print(f"  - 60 requests/minuto")
        print(f"  - 3600 requests/hora")
        print(f"  - 50000 requests/dia")

        print(f"\n📚 Próximos Passos:")
        print(f"  1. Use estas credenciais para obter um token OAuth2")
        print(f"  2. Teste os endpoints da API")
        print(f"  3. Implemente integração no seu app Python/FastAPI")
        print(f"\n  Exemplo de uso:")
        print(f"  ```python")
        print(f"  # Obter token")
        print(f"  response = requests.post('http://localhost:8000/oauth/token', data={{")
        print(f"      'client_id': '{client_id}',")
        print(f"      'client_secret': '{client_secret}'")
        print(f"  }})")
        print(f"  token = response.json()['access_token']")
        print(f"  ")
        print(f"  # Usar token")
        print(f"  headers = {{'Authorization': f'Bearer {{token}}'}}")
        print(f"  response = requests.get('http://localhost:8000/api/v1/patients/1', headers=headers)")
        print(f"  ```")

        cur.close()
        conn.close()

    except psycopg2.Error as e:
        print(f"❌ Erro PostgreSQL: {e}")
    except Exception as e:
        print(f"❌ Erro: {e}")

if __name__ == "__main__":
    create_test_client()
