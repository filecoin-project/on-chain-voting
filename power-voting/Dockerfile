FROM node:22 AS build

WORKDIR /app

COPY package*.json ./

RUN npm config delete proxy
RUN npm config delete https-proxy

RUN npm install

COPY . .

RUN yarn build

FROM nginx:alpine

COPY --from=build /app/dist /voting/dev/dist

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
