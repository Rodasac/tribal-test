FROM node:alpine
ENV FLASK_APP main

WORKDIR /usr/src/app

COPY node .

RUN npm install

EXPOSE 3000

CMD [ "npx", "nodemon", "main.ts" ]