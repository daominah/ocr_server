FROM golang:1.14.9

RUN cat /etc/os-release

RUN apt-get -qy update
RUN apt-get install -qy libleptonica-dev libtesseract-dev tesseract-ocr

# Load languages
RUN apt-get install -y tesseract-ocr-eng tesseract-ocr-vie

COPY go.mod /go.mod
COPY go.sum /go.sum
RUN cd / && go mod download

COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN go build -o ocr_server

ENV PORT=35735
CMD $GOPATH/src/app/ocr_server
