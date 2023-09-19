FROM golang:1.21-alpine as build-backend
RUN apk add git
ADD . /build
WORKDIR /build

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o main


FROM node:16-alpine as build-frontend

ADD ./frontend /frontend
WORKDIR /frontend

RUN npm install

FROM alpine:latest
COPY --from=build-backend /build/main ./app/main
COPY --from=build-frontend /frontend ./app/frontend

WORKDIR /app
RUN chmod +x ./main
EXPOSE 8080
CMD ["./main"]
