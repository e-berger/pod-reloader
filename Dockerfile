FROM golang:1.21 as builder
ARG APPUSER=appuser
ENV USER=${APPUSER}
ENV UID=1001
WORKDIR /source
RUN mkdir app && \
    apt-get update && apt-get install -y build-essential
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -installsuffix cgo -o ./app/main ./cmd/

FROM golang:1.21

COPY --from=builder /source/app /app
WORKDIR /app

ENTRYPOINT ["./main"]
