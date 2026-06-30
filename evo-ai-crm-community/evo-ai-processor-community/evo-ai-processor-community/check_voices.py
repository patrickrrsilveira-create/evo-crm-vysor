import urllib.request
import json
url = 'https://texttospeech.googleapis.com/v1/voices?key=AIzaSyC3c8gRxB_YccNyEoOmQvXCpfgJYZe6PxY'
req = urllib.request.Request(url)
data = json.loads(urllib.request.urlopen(req).read().decode())
for v in data['voices']:
    if v['name'] in ['pt-BR-Chirp3-HD-Schedar', 'pt-BR-Chirp3-HD-Achird', 'pt-BR-Chirp3-HD-Achernar', 'pt-BR-Neural2-A', 'pt-BR-Neural2-B', 'pt-BR-Neural2-C']:
        print(f"{v['name']}: {v['ssmlGender']}")
