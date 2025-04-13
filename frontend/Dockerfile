# Stage 1: Build frontend
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev
COPY . .
RUN npm run build

# Stage 2: Serve frontend with Nginx on Port 8080 explicitly
FROM nginx:1.25-alpine
COPY --from=builder /app/build /usr/share/nginx/html

# Update Nginx config to listen on 8080
RUN sed -i 's/listen  .*/listen 8080;/g' /etc/nginx/conf.d/default.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
