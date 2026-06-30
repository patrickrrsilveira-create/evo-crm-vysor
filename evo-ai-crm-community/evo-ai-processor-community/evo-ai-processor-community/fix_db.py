import psycopg2

def fix_db():
    conn = psycopg2.connect("postgres://postgres:evoai_dev_password@localhost:5432/evo_community")
    cur = conn.cursor()
    cur.execute("UPDATE evo_core_community_schema_migrations SET dirty=false WHERE version=17;")
    conn.commit()
    print("Database dirty flag cleared!")
    cur.close()
    conn.close()

if __name__ == '__main__':
    fix_db()
