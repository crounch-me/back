FROM node:lts-alpine as builder

WORKDIR /app

COPY package*.json ./
RUN npm i --production
COPY . .
RUN npm rebuild node-sass
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]