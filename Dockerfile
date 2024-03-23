FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod tidy
CMD [ "/app/movie_api" ]
