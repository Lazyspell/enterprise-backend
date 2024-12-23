#Build the Go Binary.
FROM golang:1.23.4 as build_sales
ENV CGO_ENABLE 0
ARG BUILD_REF

# Copy the source code into the container
COPY . /enterprise-backend

#BUild the service binary
WORKDIR /enterprise-backend/apis/services/sales
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.20
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /enterprise-backend -G sales -S sales
COPY --from=build_sales --chown=sales:sales /enterprise-backend/apis/services/sales/sales /enterprise-backend/sales
WORKDIR /enterprise-backend
USER sales
CMD ["./sales"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="Jeremy Elam <jelam2975@gmail.com" \
      org.opencontainers.image.source="github.com/lazyspell/enterprise-backend/sales" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Lazyspell"
