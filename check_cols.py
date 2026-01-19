import psycopg2

DB_URL = "postgres://postgres:Debian23%40@127.0.0.1:5432/eva-db?sslmode=disable"

def check():
    try:
        conn = psycopg2.connect(DB_URL)
        cur = conn.cursor()
        cur.execute("SELECT column_name FROM information_schema.columns WHERE table_name = 'idosos';")
        cols = [c[0] for c in cur.fetchall()]
        print(f"COLUMNS: {cols}")
        conn.close()
    except Exception as e:
        print(f"ERROR: {e}")

if __name__ == "__main__":
    check()
