import sys
import os
sys.path.append(os.getcwd())
from src.services.adk.tools.text_to_speech import _clean_text_for_tts

text = 'Patrick agora eu fiquei preocupado! A gente tá tendo algum problema de conexão aí. Vou tentar mandar mais uma vez pra ver se vai. Se não chegar dessa vez eu vou pedir pra minha equipe técnica dar uma olhada agora mesmo pra gente resolver isso pra você!\n[VIDEO_LINK: https://minio-api.vysortech.app.br/midias/Ganader_Brasil.mp4]'
print('CLEANED:')
print(_clean_text_for_tts(text))
