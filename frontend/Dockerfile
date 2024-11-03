# syntax=docker/dockerfile:1

FROM node:23-alpine
WORKDIR /app

COPY . .

RUN npm install

CMD ["npm", "start"]