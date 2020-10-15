FROM golang:1.14.9

RUN cat /etc/os-release
RUN apt-get -qy update

# install tesseract

RUN apt-get install -qy libleptonica-dev libtesseract-dev
RUN cd /opt && git clone https://github.com/tesseract-ocr/tesseract
WORKDIR /opt/tesseract
RUN git reset --hard 3.05.02
RUN apt-get install -qy libtool m4 automake cmake pkg-config
RUN apt-get install -qy libicu-dev libpango1.0-dev libcairo-dev
RUN ./autogen.sh &&\
    ./configure --enable-debug LDFLAGS="-L/usr/local/lib" CFLAGS="-I/usr/local/include"
RUN make -j 8
RUN make install && ldconfig
RUN tesseract --version

ENV TESSDATA_PREFIX=/usr/local/share/tessdata
RUN cd ${TESSDATA_PREFIX} &&\
    wget -q https://github.com/tesseract-ocr/tessdata/raw/3.04.00/eng.traineddata

# build this app

COPY go.mod /go.mod
COPY go.sum /go.sum
RUN cd / && go mod download

WORKDIR $GOPATH/src/app
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN go test -v -run=TestReadCaptcha
RUN go build -o ocr_server

ENV PORT=35735
CMD $GOPATH/src/app/ocr_server
