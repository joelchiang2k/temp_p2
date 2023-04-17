# Frontend Dockerfile
FROM node:latest as node

WORKDIR /app

COPY . .

RUN npm install

CMD ["npm", "run", "dev"]



