FROM node:20 AS builder
WORKDIR /app
ARG BACKEND_CORE_URL
ENV BACKEND_CORE_URL=$BACKEND_CORE_URL
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/templates/default.conf.template

CMD ["nginx", "-g", "daemon off;"]