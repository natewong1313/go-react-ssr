FROM node:16-alpine as build-frontend
ADD ./frontend-mui ./build
WORKDIR /build
RUN npm install

FROM golang:1.21-alpine
RUN apk add git
ADD . /app
COPY --from=build-frontend /build ./app/frontend-mui
WORKDIR /app

RUN go mod download
RUN go get -u github.com/natewong1313/go-react-ssr
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -X main.APP_ENV=production" -a -o main

RUN chmod +x ./main
EXPOSE 8080
CMD ["./main"]

