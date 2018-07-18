FROM node:7

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install --production

COPY ./dist/src .

EXPOSE 3000

CMD ["npm", "run", "start"]