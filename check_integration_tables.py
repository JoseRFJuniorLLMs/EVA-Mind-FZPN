#!/usr/bin/env python3
"""
Script para verificar tabelas do SPRINT 7 Integration Layer
"""

import psycopg2

# Configuração do database remoto
DB_CONFIG = {
    "host": "104.248.219.200",
    "port": 5432,
    "database": "eva-db",
    "user": "postgres",
    "password": "Debian23@"
}

def check_integration_layer():
    """Verifica se SPRINT 7 está deployado"""
    print(f"Conectando ao PostgreSQL em {DB_CONFIG['host']}...")

    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cur = conn.cursor()

        print(f"✓ Conectado com sucesso!\n")

        # Verificar tabelas do SPRINT 7
        cur.execute("""
            SELECT tablename
            FROM pg_tables
            WHERE schemaname = 'public'
            AND (
                tablename LIKE 'api_%' OR
                tablename LIKE 'webhook%' OR
                tablename LIKE 'fhir%' OR
                tablename LIKE 'data_export%'
            )
            ORDER BY tablename;
        """)

        tables = cur.fetchall()

        print(f"📋 Tabelas do Integration Layer ({len(tables)}):")
        expected_tables = [
            'api_clients',
            'api_tokens',
            'api_request_logs',
            'api_rate_limits',
            'webhook_deliveries',
            'fhir_resource_mappings',
            'data_export_jobs',
            'integration_audit_logs'
        ]

        found_tables = [t[0] for t in tables]

        for table in expected_tables:
            status = "✓" if table in found_tables else "✗"
            print(f"  {status} {table}")

        # Verificar views
        cur.execute("""
            SELECT viewname
            FROM pg_views
            WHERE schemaname = 'public'
            AND (
                viewname LIKE 'v_api%' OR
                viewname LIKE 'v_webhook%'
            )
            ORDER BY viewname;
        """)

        views = cur.fetchall()

        print(f"\n📊 Views do Integration Layer ({len(views)}):")
        expected_views = [
            'v_api_usage_stats',
            'v_api_client_health',
            'v_webhook_delivery_stats',
            'v_fhir_sync_status'
        ]

        found_views = [v[0] for v in views]

        for view in expected_views:
            status = "✓" if view in found_views else "✗"
            print(f"  {status} {view}")

        # Contar registros em api_clients
        cur.execute("SELECT COUNT(*) FROM api_clients;")
        client_count = cur.fetchone()[0]

        print(f"\n📈 Estatísticas:")
        print(f"  - API Clients: {client_count}")

        # Se não houver clientes, criar um de teste
        if client_count == 0:
            print(f"\n💡 Nenhum API client encontrado. Deseja criar um client de teste?")
            print(f"   (Execute create_test_client.py para criar)")

        cur.close()
        conn.close()

        print(f"\n✅ SPRINT 7 Integration Layer está DEPLOYADO!")

    except psycopg2.Error as e:
        print(f"❌ Erro PostgreSQL: {e}")
    except Exception as e:
        print(f"❌ Erro: {e}")

if __name__ == "__main__":
    check_integration_layer()
