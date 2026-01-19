import os
import psycopg2
from urllib.parse import urlparse

# DATABASE_URL from .env
DB_URL = "postgres://postgres:Debian23%40@127.0.0.1:5432/eva-db?sslmode=disable"

def migrate():
    print("🚀 Iniciando migração de banco de dados (Voz)...")
    
    try:
        conn = psycopg2.connect(DB_URL)
        conn.autocommit = True
        cur = conn.cursor()
        
        # 1. Verifica se a coluna já existe
        cur.execute("""
            SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name='idosos' AND column_name='voice_name';
        """)
        exists = cur.fetchone()
        
        if exists:
            print("⚠️ Coluna 'voice_name' já existe. Nada a fazer.")
        else:
            # 2. Cria a coluna com valor default 'Aoede'
            print("Creating column 'voice_name'...")
            cur.execute("""
                ALTER TABLE idosos 
                ADD COLUMN voice_name VARCHAR(50) DEFAULT 'Aoede';
            """)
            print("✅ Coluna 'voice_name' criada com sucesso!")
            
            # 3. Atualiza registros existentes (redundante com DEFAULT mas bom garantir)
            cur.execute("UPDATE idosos SET voice_name = 'Aoede' WHERE voice_name IS NULL;")
            print("✅ Registros existentes atualizados.")

        conn.close()
        print("✅ Migração concluída.")

    except Exception as e:
        print(f"❌ Erro na migração: {e}")

if __name__ == "__main__":
    migrate()
