# syntax=docker/dockerfile:1

FROM node:23-alpine
WORKDIR /app

COPY public/ /react-docker-example/public
COPY src/ /react-docker-example/src
COPY package.json /react-docker-example/

RUN npm install

CMD ["npm", "start"]