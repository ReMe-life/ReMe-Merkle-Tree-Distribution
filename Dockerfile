# Start from the latest golang base image
FROM amd64/golang:1.14 AS dev-build
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
######## Start a new stage from scratch #######
FROM amd64/golang:1.14
#Service variables
ARG DB_HOST=
ARG DB_PORT=
ARG DB_USER=
ARG DB_PASS=
ARG DB_NAME=
ARG NODE_URL=
ARG SECRET_KEY=
ARG CONTRACT_ADDRESS=
ARG SAVE_PERIOD=
ARG API_PORT=8080
#ENV set
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT
ENV DB_USER=$DB_USER
ENV DB_PASS=$DB_PASS
ENV DB_NAME=$DB_NAME
ENV NODE_URL=$NODE_URL
ENV SECRET_KEY=$SECRET_KEY
ENV CONTRACT_ADDRESS=$CONTRACT_ADDRESS
ENV SAVE_PERIOD=$SAVE_PERIOD
ENV API_PORT=$API_PORT
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=dev-build /app/main .
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./main"]