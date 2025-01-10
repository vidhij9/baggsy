# Stage 1: Build the React app
FROM node:18-alpine AS build

WORKDIR /frontend

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve with a lightweight web server (e.g., nginx)
FROM nginx:alpine

# Copy the built React files
COPY --from=build /frontend/build /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
