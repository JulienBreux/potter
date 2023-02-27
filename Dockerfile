ARG GO_VERSION=1.20.0

FROM golang:${GO_VERSION}-alpine AS build

ARG VERSION=dev
ARG DATE=n/a
ARG COMMIT=n/a

WORKDIR /src

COPY go.mod go.sum Makefile ./
COPY internal internal
COPY pkg pkg
COPY cmd cmd
COPY views views


RUN apk --no-cache add --update make libx11-dev git gcc libc-dev curl && make build

FROM gcr.io/distroless/static AS final

LABEL maintainer="Julien BREUX <julien.breux@gmail.com>"
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /src/bin/app /app
COPY --from=build --chown=nonroot:nonroot /src/views /views

ENTRYPOINT ["/app"]
