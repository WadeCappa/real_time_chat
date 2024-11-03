# syntax=docker/dockerfile:1

FROM node:23-alpine
WORKDIR /app

COPY . .

RUN npm install
RUN npm run build 
RUN npm install -g serve

EXPOSE 8000

CMD ["serve", "-s", "-p", "8000", "build"]