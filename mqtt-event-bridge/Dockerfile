FROM node:alpine

# don't install dev dependencies
ENV NODE_ENV=production

# Create app directory
WORKDIR /usr/src/app

# Bundle app source
COPY ./package*.json ./
RUN npm install

COPY . .
EXPOSE 8080

CMD [ "npm", "start" ]