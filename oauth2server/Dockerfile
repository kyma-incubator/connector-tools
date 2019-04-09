FROM node:alpine

# don't install dev dependencies
ENV NODE_ENV=production

# Create app directory
WORKDIR /usr/src/app

# Bundle app source
COPY ./package*.json ./
COPY ./app.js ./
COPY ./model.js ./
RUN npm install

EXPOSE 8080

CMD [ "npm", "start" ]