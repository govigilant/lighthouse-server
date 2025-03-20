FROM golang:1.24-alpine

RUN apk add --no-cache curl nodejs npm chromium

ENV CHROME_PATH=/usr/bin/chromium-browser

RUN npm install -g lighthouse

WORKDIR /app

COPY . .

RUN go build -o lighthouse-server

CMD ["./lighthouse-server"]

