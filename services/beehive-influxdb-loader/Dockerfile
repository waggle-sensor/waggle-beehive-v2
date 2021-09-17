FROM python:3-alpine
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt
COPY . .
ENTRYPOINT [ "python", "main.py" ]
