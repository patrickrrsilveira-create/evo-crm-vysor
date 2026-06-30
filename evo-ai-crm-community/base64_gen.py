import base64

s = '''import re, os
with open('src/api/a2a_routes.py', 'r', encoding='utf-8') as f: c = f.read()

hack = """if url:
    import httpx, asyncio
    phone = params.get('contact', {}).get('phone_number', '')
    if phone:
        phone = phone.replace('+', '').replace('-', '').replace(' ', '').strip()
        try:
            api_url = 'http://evolution_go:4000/message/sendWhatsAppAudio/evolution_go'
            headers = {'apikey': 'X8G9W2M4V5N7B3L1K6J0H9P2Y3T5C8F1', 'Content-Type': 'application/json'}
            asyncio.create_task(httpx.AsyncClient().post(api_url, headers=headers, json={'number': phone, 'audio': url}, timeout=10))
            artifacts.append({'artifactId': str(uuid.uuid4()), 'parts': [{'type': 'text', 'text': '🎙️ Áudio enviado pelo Evolution'}]})
        except:
            artifacts.append({'artifactId': str(uuid.uuid4()), 'parts': [{'type': 'file', 'url': url, 'mimeType': mime_type}]})
    else:
        artifacts.append({'artifactId': str(uuid.uuid4()), 'parts': [{'type': 'file', 'url': url, 'mimeType': mime_type}]})
"""

p = r'(\s*)if url:\s*artifacts\.append\(\{\s*"artifactId": str\(uuid\.uuid4\(\)\),\s*"parts": \[\{\s*"type": "audio",\s*"url": url,\s*"mimeType": mime_type\s*\}\]\s*\}\)'

c = re.sub(p, lambda m: m.group(1) + hack.replace('\\n', '\\n' + m.group(1)), c)
c = c.replace('"type": "audio"', '"type": "file"')
with open('src/api/a2a_routes.py', 'w', encoding='utf-8') as f: f.write(c)
'''

with open('hack.b64', 'w') as f:
    f.write(base64.b64encode(s.encode('utf-8')).decode('utf-8'))
