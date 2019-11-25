FROM python:3.7-alpine

WORKDIR /app
COPY . .
RUN pip install -r requirements.txt

CMD tavern-ci test_easyclai_api.tavern.yaml -v