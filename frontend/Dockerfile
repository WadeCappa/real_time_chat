# syntax=docker/dockerfile:1

FROM node:23-alpine
WORKDIR /app

COPY . .

ARG MODE

RUN npm install
RUN npm install -g serve

RUN ./setup_env.sh $MODE
RUN npm run build 

EXPOSE 8000

CMD ["serve", "-s", "-p", "3000", "build"]