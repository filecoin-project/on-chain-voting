# 使用 Node.js 镜像作为构建环境
FROM node:22 as build-stage

# 设置工作目录
WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY power-voting/package*.json ./

# 安装依赖
RUN yarn install

# 复制项目文件
COPY power-voting .

# 构建项目
RUN yarn build

# 使用 Nginx 镜像作为生产环境
FROM nginx:alpine

# 将构建产物复制到 Nginx 的静态文件目录
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 复制自定义的 Nginx 配置文件
COPY power-voting/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露 80 端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]