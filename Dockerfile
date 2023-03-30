# Frontend Dockerfile
FROM node:latest as node

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

# Backend Dockerfile
FROM golang:latest

WORKDIR /app

COPY . .

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y nodejs

RUN curl -sL https://deb.nodesource.com/setup_16.x | bash -

RUN apt-get install -y nodejs

RUN npm install  

RUN go build -o main .

CMD ["sh", "-c", "npm run dev & go run main.go"]


