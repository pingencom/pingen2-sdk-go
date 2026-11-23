FROM golang:1.27

WORKDIR /app

COPY . .

RUN go mod download

CMD ["bash"]
