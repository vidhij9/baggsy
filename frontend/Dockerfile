# Using latest recommended LTS Alpine Node image (secure version)
FROM node:20.12.2-alpine3.19 AS builder

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev
COPY . .
RUN npm run build

# Using secure nginx-alpine stable image with fixed vulnerabilities
FROM nginxinc/nginx-unprivileged:1.25-alpine3.19-slim

# Copy built files to nginx
COPY --from=builder /app/build /usr/share/nginx/html

USER 101
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
