import urllib.request
import json

url = 'https://openrouter.ai/api/v1/audio/speech'
data = json.dumps({'model':'hexgrad/kokoro-82m','input':'test','voice':'af_heart'}).encode('utf-8')
req = urllib.request.Request(url, data=data, headers={'Content-Type': 'application/json'})

try:
    with urllib.request.urlopen(req) as response:
        print("Status:", response.status)
        print("Content-Type:", response.headers.get('Content-Type'))
        print("Body:", response.read()[:200])
except urllib.error.HTTPError as e:
    print("Status:", e.code)
    print("Content-Type:", e.headers.get('Content-Type'))
    print("Body:", e.read()[:200])
