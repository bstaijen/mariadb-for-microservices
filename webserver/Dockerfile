FROM node:argon

# Create app directory
RUN mkdir -p /usr/src/app

# Copy files to app directory
COPY . /usr/src/app

# Install bower and bower_components for the webapp
WORKDIR /usr/src/app/webapp
RUN npm install -g bower
RUN bower install --allow-root

# Install components for proxy
WORKDIR /usr/src/app/proxy
RUN npm install

# Start
EXPOSE 4999
CMD [ "npm", "start" ]