# Base image for building the app
FROM node:18-alpine AS builder
ARG REACT_APP_IPFS_URL
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN REACT_APP_IPFS_URL=${REACT_APP_IPFS_URL} npm run build

# Base image for running the app
FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
