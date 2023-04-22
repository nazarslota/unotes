FROM node:19-alpine as builder

WORKDIR /app
COPY . .

RUN npm install
RUN npm run build

FROM nginx:1-alpine

COPY --from=builder /app/build /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]
