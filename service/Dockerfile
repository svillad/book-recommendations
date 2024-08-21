FROM golang:latest

ENV SDIR=$GOPATH/src/github.com/book-recommendations/service

# Specify that we now need to execute any commands in this directory.
WORKDIR $SDIR
ADD . $SDIR

# RUN go get
COPY go.mod go.sum ./
RUN go mod download

# Compile the binary exe for our app.
RUN go build -o books_service .

# Start the application.
ENTRYPOINT ["./books_service"]