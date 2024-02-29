FROM golang:1.21 as builder
ARG APPUSER=appuser
ENV USER=${APPUSER}
ENV UID=1001
WORKDIR /source
RUN mkdir app && \
    apt-get update && apt-get install -y build-essential && \
    adduser --disabled-password --gecos "" --no-create-home --uid "${UID}" "${USER}"
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -installsuffix cgo -o ./app/main ./cmd/

FROM scratch
ARG APPUSER=appuser
ENV UID=1001
EXPOSE 8080

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=${APPUSER}:${APPUSER} /source/app /app
WORKDIR /app
USER ${APPUSER}:${APPUSER}

ENTRYPOINT ["./main"]
