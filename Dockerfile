FROM golang:1.15-alpine AS build
WORKDIR /src
COPY * /src
RUN CGO_ENABLED=0 go build -o proxy

FROM scratch
COPY --from=build /src/proxy /src/proxy
ENTRYPOINT ["/src/proxy"]