FROM node:18-alpine3.15 as build

RUN apk add --no-cache protoc

WORKDIR /app

ENV protoc_gen_grpc_web=./node_modules/.bin/protoc-gen-grpc-web

COPY ./package.json .
COPY ./yarn.lock .

RUN yarn install

COPY ./ .

EXPOSE 8081

RUN yarn global add serve
RUN yarn build

CMD ["serve" "-s" "build" "-l" "8081"]