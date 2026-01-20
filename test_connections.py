#!/usr/bin/env python3
"""
Script de teste de conexão Redis e Neo4j
Usa as credenciais do .env
"""

import redis
from neo4j import GraphDatabase
import sys

# Configurações (ajuste conforme seu .env)
REDIS_HOST = "127.0.0.1"
REDIS_PORT = 6379
REDIS_PASSWORD = None  # Se tiver senha, coloque aqui

NEO4J_URI = "bolt://localhost:7687"
NEO4J_USERNAME = "neo4j"
NEO4J_PASSWORD = "Debian23@"  # Ajuste conforme necessário

def test_redis():
    """Testa conexão com Redis"""
    print("🔍 Testando Redis...")
    try:
        r = redis.Redis(
            host=REDIS_HOST,
            port=REDIS_PORT,
            password=REDIS_PASSWORD,
            decode_responses=True
        )
        
        # Teste de ping
        r.ping()
        print("✅ Redis: Conectado com sucesso!")
        
        # Teste de escrita/leitura
        r.set("test_key", "test_value")
        value = r.get("test_key")
        print(f"   Teste de escrita/leitura: {value}")
        r.delete("test_key")
        
        return True
    except redis.ConnectionError as e:
        print(f"❌ Redis: Erro de conexão - {e}")
        return False
    except redis.AuthenticationError as e:
        print(f"❌ Redis: Erro de autenticação - {e}")
        print("   Dica: Verifique se o Redis requer senha")
        return False
    except Exception as e:
        print(f"❌ Redis: Erro desconhecido - {e}")
        return False

def test_neo4j():
    """Testa conexão com Neo4j"""
    print("\n🔍 Testando Neo4j...")
    try:
        driver = GraphDatabase.driver(
            NEO4J_URI,
            auth=(NEO4J_USERNAME, NEO4J_PASSWORD)
        )
        
        # Verificar conectividade
        driver.verify_connectivity()
        print("✅ Neo4j: Conectado com sucesso!")
        
        # Teste de query simples
        with driver.session() as session:
            result = session.run("RETURN 1 AS num")
            record = result.single()
            print(f"   Teste de query: {record['num']}")
        
        driver.close()
        return True
    except Exception as e:
        print(f"❌ Neo4j: Erro - {e}")
        print(f"   URI: {NEO4J_URI}")
        print(f"   User: {NEO4J_USERNAME}")
        return False

def main():
    print("=" * 50)
    print("🧪 Teste de Conexões - EVA-Mind")
    print("=" * 50)
    
    redis_ok = test_redis()
    neo4j_ok = test_neo4j()
    
    print("\n" + "=" * 50)
    print("📊 Resumo:")
    print("=" * 50)
    print(f"Redis:  {'✅ OK' if redis_ok else '❌ FALHOU'}")
    print(f"Neo4j:  {'✅ OK' if neo4j_ok else '❌ FALHOU'}")
    
    if redis_ok and neo4j_ok:
        print("\n🎉 Todos os serviços estão funcionando!")
        sys.exit(0)
    else:
        print("\n⚠️ Alguns serviços falharam. Verifique as configurações.")
        sys.exit(1)

if __name__ == "__main__":
    main()
