FROM golang:1.24

RUN apt update && apt install -y curl && \
	curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
	apt install -y nodejs npm chromium && \
	npm install -g lighthouse

WORKDIR /app

COPY . .

RUN go build -o lighthouse-server

CMD ["./lighthouse-server"]

