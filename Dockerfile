# Frontend Dockerfile
FROM node:latest as node

WORKDIR /app

COPY . .

RUN npm install

# Backend Dockerfile
FROM golang:latest

WORKDIR /app

COPY . .

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y nodejs

RUN curl -sL https://deb.nodesource.com/setup_16.x | bash -
ENV PORT=8001
RUN apt-get install -y nodejs

RUN npm install  

RUN go build -o main .

CMD ["sh", "-c", "go run main.go & npm run lol"]


