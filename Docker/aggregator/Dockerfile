# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Conor Hackett <chackett1@gmail.com>"

# Set the Current Working Directory inside the container
#RUN mkdir /app
WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN cd cmd/aggregator/ && go build -o /app/app .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]