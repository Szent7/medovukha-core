FROM node:alpine AS frontend-builder
WORKDIR /frontend
ADD ./frontend/package.json .
ADD ./frontend/package-lock.json .
RUN npm install
COPY ./frontend .
RUN npm run build

FROM golang:alpine AS backend-builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go 

FROM alpine
WORKDIR /app/build
COPY --from=frontend-builder /frontend/build /app/build/
WORKDIR /app
COPY --from=backend-builder /build/main /app/main

EXPOSE 10015

CMD ["./main"]