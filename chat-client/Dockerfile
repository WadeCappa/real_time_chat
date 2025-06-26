# Build stage
FROM node:24-alpine AS build

WORKDIR /usr/src/app

ENV NODE_ENV production

COPY package*.json ./

RUN npm ci 

COPY . .

RUN npm run build

# Deploy stage
FROM nginx:1.29-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=build /usr/src/app/build/ /usr/share/nginx/html

EXPOSE 8081
