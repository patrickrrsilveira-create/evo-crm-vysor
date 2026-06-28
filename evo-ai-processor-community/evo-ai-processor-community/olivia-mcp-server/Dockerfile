FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY main.py .
# The fastmcp default stdio server runs when executed directly.
CMD ["python", "main.py"]
