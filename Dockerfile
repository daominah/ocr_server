FROM golang:1.14.9

RUN cat /etc/os-release
RUN apt-get -qy update

# install tesseract

RUN apt-get install -qy libleptonica-dev libtesseract-dev
RUN apt-get install -qy libtool m4 automake cmake pkg-config
RUN apt-get install -qy libicu-dev libpango1.0-dev libcairo-dev

RUN cd /opt && git clone https://github.com/tesseract-ocr/tesseract
WORKDIR /opt/tesseract
RUN git reset --hard 4.1.1
RUN ./autogen.sh &&\
    ./configure --enable-debug LDFLAGS="-L/usr/local/lib" CFLAGS="-I/usr/local/include"
RUN make -j 8
RUN make install && ldconfig
RUN tesseract --version

ENV TESSDATA_PREFIX=/usr/local/share/tessdata
ENV TESSDATA_REPO=https://github.com/tesseract-ocr/tessdata_best
WORKDIR ${TESSDATA_PREFIX}
RUN wget -q ${TESSDATA_REPO}/raw/4.1.0/eng.traineddata
RUN wget -q ${TESSDATA_REPO}/raw/4.1.0/vie.traineddata
RUN wget -q ${TESSDATA_REPO}/raw/4.1.0/chi_sim.traineddata

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
