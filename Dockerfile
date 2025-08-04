FROM golang:1.24
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download all
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /otusgruz
CMD [ "/otusgruz" ]