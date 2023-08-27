FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

CMD ["/main"]

# This dockerfile will make a container which will be stored on a cloud image repository and then used to deploy a
# K8s job that runs once per day to collect the data needed
# The data will be stored in a database and then used to create a dashboard
# Since most processes are parallilsed the job should run in almost constant time and run in well under 1 minute (~2 seconds)

