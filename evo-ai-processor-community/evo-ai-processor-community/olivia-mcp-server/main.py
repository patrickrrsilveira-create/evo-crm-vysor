import sqlite3
import datetime
import json
from mcp.server.fastmcp import FastMCP

# Create a FastMCP server
mcp = FastMCP("Olivia Executiva MCP")

# Configurar banco de dados local
DB_FILE = "olivia_mock.db"

def init_db():
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    c.execute('''
        CREATE TABLE IF NOT EXISTS meetings (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT,
            time TEXT,
            participants TEXT
        )
    ''')
    c.execute('''
        CREATE TABLE IF NOT EXISTS emails (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            recipient TEXT,
            subject TEXT,
            body TEXT,
            sent_at TEXT
        )
    ''')
    c.execute('''
        CREATE TABLE IF NOT EXISTS notes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            topic TEXT,
            content TEXT,
            created_at TEXT
        )
    ''')
    conn.commit()
    conn.close()

init_db()

@mcp.tool()
def check_calendar(date: str = None) -> str:
    """Verifica a agenda da diretoria para um dia específico.
    Args:
        date: A data no formato YYYY-MM-DD. Se vazio, pega o dia de hoje.
    """
    if not date:
        date = datetime.date.today().strftime("%Y-%m-%d")
        
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    c.execute("SELECT title, time, participants FROM meetings WHERE time LIKE ?", (f"{date}%",))
    meetings = c.fetchall()
    conn.close()
    
    if not meetings:
        return f"A agenda para {date} está livre."
    
    result = f"Agenda para {date}:\n"
    for m in meetings:
        result += f"- {m[1]}: {m[0]} (Com: {m[2]})\n"
    return result

@mcp.tool()
def schedule_meeting(title: str, time: str, participants: str) -> str:
    """Agenda um novo compromisso na agenda da diretoria.
    Args:
        title: O título da reunião.
        time: Data e hora no formato YYYY-MM-DD HH:MM.
        participants: Nomes ou e-mails dos participantes.
    """
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    c.execute("INSERT INTO meetings (title, time, participants) VALUES (?, ?, ?)", (title, time, participants))
    conn.commit()
    conn.close()
    return f"Reunião '{title}' agendada com sucesso para {time} com {participants}."

@mcp.tool()
def send_email(recipient: str, subject: str, body: str) -> str:
    """Envia um e--mail em nome da diretoria.
    Args:
        recipient: E-mail do destinatário.
        subject: Assunto do e-mail.
        body: Corpo do e-mail.
    """
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    sent_at = datetime.datetime.now().isoformat()
    c.execute("INSERT INTO emails (recipient, subject, body, sent_at) VALUES (?, ?, ?, ?)", (recipient, subject, body, sent_at))
    conn.commit()
    conn.close()
    return f"E-mail enviado para {recipient} com assunto '{subject}'."

@mcp.tool()
def save_note(topic: str, content: str) -> str:
    """Salva uma anotação, ata de reunião ou relatório no arquivo da secretaria.
    Args:
        topic: O título ou tópico da anotação.
        content: O conteúdo completo do relatório ou anotação.
    """
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    created_at = datetime.datetime.now().isoformat()
    c.execute("INSERT INTO notes (topic, content, created_at) VALUES (?, ?, ?)", (topic, content, created_at))
    conn.commit()
    conn.close()
    return f"Anotação sobre '{topic}' salva com sucesso."

@mcp.tool()
def read_reports(topic_query: str = "") -> str:
    """Busca relatórios e anotações salvos pela secretaria.
    Args:
        topic_query: Palavra-chave para buscar no tópico. Deixe vazio para listar as últimas 5.
    """
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    if topic_query:
        c.execute("SELECT topic, content, created_at FROM notes WHERE topic LIKE ? ORDER BY created_at DESC LIMIT 5", (f"%{topic_query}%",))
    else:
        c.execute("SELECT topic, content, created_at FROM notes ORDER BY created_at DESC LIMIT 5")
    notes = c.fetchall()
    conn.close()
    
    if not notes:
        return "Nenhuma anotação encontrada."
        
    result = "Relatórios/Anotações encontrados:\n"
    for n in notes:
        result += f"\n--- {n[0]} ({n[2]}) ---\n{n[1]}\n"
    return result

if __name__ == "__main__":
    mcp.run()
