import os
import sys

# Insert current dir so we can import src
sys.path.insert(0, os.path.abspath(os.path.dirname(__file__)))

os.environ["ARTIFACT_SERVICE_TYPE"] = "minio"
os.environ["ARTIFACT_ENDPOINT"] = "minio-api.vysortech.app.br"
os.environ["ARTIFACT_ACCESS_KEY"] = "admin"
os.environ["ARTIFACT_SECRET_KEY"] = "vysor@BR001"

from src.services.adk.artifacts.providers.minio_provider import MinIOProviderConfig
config = MinIOProviderConfig()
print("Enabled:", config.enabled)
print("Configured:", config.is_configured())
print("Valid:", config.validate())
