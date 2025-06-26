FROM node:24-alpine

WORKDIR /app

ENV NODE_ENV production

COPY package*.json ./

RUN npm install serve && npm ci --production

COPY . .

RUN npm run build

EXPOSE 8081

CMD ["./node_modules/.bin/serve", "-s", "build", "-l", "8081"]