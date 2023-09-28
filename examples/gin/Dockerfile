FROM golang:1.21-alpine as build-backend
RUN apk add git
ADD . /build
WORKDIR /build

RUN go mod download
RUN go get -u github.com/natewong1313/go-react-ssr
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -X main.APP_ENV=production" -a -o main


FROM node:16-alpine as build-frontend

ADD ./frontend-tailwind /frontend-tailwind
WORKDIR /frontend-tailwind

RUN npm install

FROM node:16-alpine
COPY --from=build-backend /build/main ./app/main
COPY --from=build-frontend /frontend-tailwind ./app/frontend-tailwind

WORKDIR /app
RUN chmod +x ./main
EXPOSE 8080
CMD ["./main"]

