FROM node:19-alpine as builder

WORKDIR /usr/app

COPY package*.json /usr/app/
RUN npm install

COPY . /usr/app/
COPY .env .env.local? /usr/app/

RUN npm run build

FROM nginx:1.23-alpine

EXPOSE 80
EXPOSE 443

RUN apk update && apk add --no-cache openssl
RUN openssl req -x509 -nodes -days 365 -subj "/O=Unotes/CN=unotes" -newkey rsa:2048 \
    -keyout /etc/ssl/private/selfsigned.key -out /etc/ssl/certs/selfsigned.crt;

COPY ./docker/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/app/build/ /usr/share/nginx/html/

CMD ["nginx", "-g", "daemon off;"]
