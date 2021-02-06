FROM golang:latest
RUN mkdir /app
COPY . /app
WORKDIR /app
ENV GOPROXY="https://goproxy.io,direct"
RUN go mod download
RUN go build -o go-xops .
EXPOSE 9000
CMD ["/app/go-xops"]