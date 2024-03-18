# Builder stage: Build the Go application
FROM golang:1.22 AS builder
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code and build it
COPY . .
# Ensure CGO is enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o myapp .

# Final stage: Use a Debian image that matches the builder stage
# FROM debian:bullseye-slim
# WORKDIR /app

# # Install necessary packages for Oracle Instant Client
# RUN apt-get update && apt-get install -y libaio1 unzip && \
#     apt-get clean && rm -rf /var/lib/apt/lists/*

# Final stage: Set up the environment with Oracle Instant Client on Ubuntu
FROM ubuntu:latest
WORKDIR /app

# Install necessary packages for Oracle Instant Client and other dependencies
# Ensure 'unzip' is included in the packages to install
RUN apt-get update && apt-get install -y libaio1 unzip && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Add the Oracle Instant Client zip file to the container and unzip it
COPY instantclient-basic-linux.x64-21.13.0.0.0dbru.zip /opt/oracle/
RUN unzip /opt/oracle/instantclient-basic-linux.x64-21.13.0.0.0dbru.zip -d /opt/oracle && \
    rm /opt/oracle/instantclient-basic-linux.x64-21.13.0.0.0dbru.zip && \
    ln -s /opt/oracle/instantclient_21_13 /opt/oracle/instantclient && \
    echo /opt/oracle/instantclient > /etc/ld.so.conf.d/oracle-instantclient.conf && \
    ldconfig

# Environment variables for Oracle Instant Client
ENV LD_LIBRARY_PATH=/opt/oracle/instantclient:$LD_LIBRARY_PATH
ENV ORACLE_BASE=/opt/oracle/instantclient
ENV TNS_ADMIN=/opt/oracle/instantclient
ENV ORACLE_HOME=/opt/oracle/instantclient

# Copy the built Go application from the builder stage
COPY --from=builder /app/myapp /app/

# Command to run the application
CMD ["./myapp"]
