FROM pmdcosta/golang:1.11 AS builder

WORKDIR $GOPATH/src/github.com/pmdcosta/elklogs

# Update vendor dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Add code and compile it
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

# Final image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./
ENTRYPOINT ["./app"]
