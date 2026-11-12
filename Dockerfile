FROM golang:1.24-bookworm

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v2.13.2/install.sh \
    | sh -s -- -b "$(go env GOPATH)/bin" v2.13.2

RUN git config --global --add safe.directory /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["bash"]
