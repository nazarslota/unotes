FROM node:19-alpine as builder

WORKDIR /usr/app

COPY package*.json /usr/app/
RUN npm install

COPY . /usr/app/
COPY .env /usr/app/
COPY .env.local /usr/app/

RUN npm run build

FROM nginx:1.23-alpine

EXPOSE 80

COPY ./docker/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/app/build/ /usr/share/nginx/html/

CMD ["nginx", "-g", "daemon off;"]
