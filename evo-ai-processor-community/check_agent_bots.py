import psycopg2

def check():
    conn = psycopg2.connect("postgres://postgres:evoai_dev_password@localhost:5433/evo_community")
    cur = conn.cursor()
    cur.execute("SELECT id, name, outgoing_url FROM agent_bots;")
    rows = cur.fetchall()
    for row in rows:
        print(row)
    cur.close()
    conn.close()

if __name__ == '__main__':
    check()
