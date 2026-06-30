import sys
import os
sys.path.append(os.getcwd())
from src.services.adk.tools.text_to_speech import _clean_text_for_tts

text = 'Patrick agora eu fiquei preocupado! \n[VIDEO_LINK: https://minio-api.vysortech.app.br/midias/Ganader_Brasil.mp4]'
print('CLEANED:')
print(_clean_text_for_tts(text))
