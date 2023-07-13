FROM python:alpine
WORKDIR /app
RUN apk add git
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
COPY . .
CMD ["python", "main.py"]
