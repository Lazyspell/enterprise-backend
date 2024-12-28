# Build the Go Binary
FROM golang:1.23 AS build_sales
ENV CGO_ENABLED=0
ARG BUILD_REF

# Copy the source code into the container
COPY . /enterprise-backend

# Build the service binary
WORKDIR /enterprise-backend/apis/services/sales
RUN go build -o sales-service -ldflags "-X main.build=${BUILD_REF}" 

# Run the Go Binary in Alpine
FROM alpine:3.20
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /service -G sales -S sales
COPY --from=build_sales --chown=sales:sales /enterprise-backend/apis/services/sales/sales-service /service/sales-service
WORKDIR /service
USER sales
CMD ["./sales-service"]  # Corrected path to the binary

# Metadata labels
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="Jeremy Elam <jelam2975@gmail.com>" \
      org.opencontainers.image.source="github.com/lazyspell/enterprise-backend/sales" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Lazyspell"

