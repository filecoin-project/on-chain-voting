FROM node:22 as build-stage

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

COPY frontend/package*.json ./

RUN pnpm install

COPY frontend .

RUN pnpm run build

FROM nginx:alpine

COPY --from=build-stage /app/dist /usr/share/nginx/html

COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]