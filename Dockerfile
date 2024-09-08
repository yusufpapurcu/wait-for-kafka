FROM public.ecr.aws/docker/library/golang:1.22 as build

ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod .
COPY go.sum .

# Copy the code into the container
COPY main.go .

# Build the application and copy somewhere convienient
RUN go build -o wait-for-kafka .

# create our new image with just the stuff we need
FROM public.ecr.aws/docker/library/alpine:3.20.3
WORKDIR /app
COPY --from=build /app/ ./

CMD ["/app/wait-for-kafka"]