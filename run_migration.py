#!/usr/bin/env python3
"""
Script para executar migration SQL no PostgreSQL remoto
"""

import psycopg2
import sys

# Configuração do database remoto
DB_CONFIG = {
    "host": "104.248.219.200",
    "port": 5432,
    "database": "eva-db",
    "user": "postgres",
    "password": "Debian23@"
}

def run_migration(sql_file):
    """Executa arquivo SQL migration"""
    print(f"Conectando ao PostgreSQL em {DB_CONFIG['host']}...")

    try:
        # Conectar ao database
        conn = psycopg2.connect(**DB_CONFIG)
        conn.autocommit = True
        cur = conn.cursor()

        print(f"✓ Conectado com sucesso!")
        print(f"Executando migration: {sql_file}...")

        # Ler arquivo SQL
        with open(sql_file, 'r', encoding='utf-8') as f:
            sql_content = f.read()

        # Executar SQL
        cur.execute(sql_content)

        print(f"✓ Migration executada com sucesso!")

        # Verificar tabelas criadas
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

        print(f"\n✓ Tabelas criadas ({len(tables)}):")
        for table in tables:
            print(f"  - {table[0]}")

        # Verificar views criadas
        cur.execute("""
            SELECT viewname
            FROM pg_views
            WHERE schemaname = 'public'
            AND viewname LIKE 'v_api%'
            ORDER BY viewname;
        """)

        views = cur.fetchall()

        print(f"\n✓ Views criadas ({len(views)}):")
        for view in views:
            print(f"  - {view[0]}")

        # Fechar conexão
        cur.close()
        conn.close()

        print("\n🎉 SPRINT 7 Integration Layer deployment concluído!")

    except psycopg2.Error as e:
        print(f"❌ Erro PostgreSQL: {e}")
        sys.exit(1)
    except FileNotFoundError:
        print(f"❌ Arquivo não encontrado: {sql_file}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Erro: {e}")
        sys.exit(1)

if __name__ == "__main__":
    migration_file = "migrations/010_integration_layer.sql"
    run_migration(migration_file)
