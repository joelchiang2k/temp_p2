# Frontend Dockerfile
FROM node:latest as node

WORKDIR /app

COPY . .

RUN npm install

CMD ["/usr/local/bin/npm", "run", "dev"]



