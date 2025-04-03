FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npx prisma migrate deploy

RUN go build -tags netgo -ldflags '-s -w' -o app

EXPOSE 4000

CMD ./app