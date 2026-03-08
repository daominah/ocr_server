# ---- build stage ----
FROM golang:1.26 AS builder
RUN apt-get -qy update

# build Tesseract from source

RUN apt-get install -qy libleptonica-dev
RUN apt-get install -qy libtool m4 automake cmake pkg-config
RUN apt-get install -qy libicu-dev libpango1.0-dev libcairo-dev
RUN ln -s /usr/lib/x86_64-linux-gnu/libleptonica.so /usr/lib/x86_64-linux-gnu/liblept.so && ldconfig

RUN cd /opt && git clone https://github.com/tesseract-ocr/tesseract
WORKDIR /opt/tesseract
RUN git checkout 5.5.2
RUN ./autogen.sh &&\
    ./configure --enable-debug LDFLAGS="-L/usr/local/lib" CFLAGS="-I/usr/local/include"
RUN make -j 8
RUN make install && ldconfig
RUN tesseract --version

ENV TESSDATA_PREFIX=/usr/local/share/tessdata
ENV TESSDATA_REPO=https://github.com/tesseract-ocr/tessdata_best
WORKDIR ${TESSDATA_PREFIX}
RUN wget -q ${TESSDATA_REPO}/raw/main/eng.traineddata
RUN wget -q ${TESSDATA_REPO}/raw/main/vie.traineddata
RUN wget -q ${TESSDATA_REPO}/raw/main/chi_sim.traineddata

# build this app

COPY go.mod /go.mod
COPY go.sum /go.sum
RUN cd / && go mod download

WORKDIR $GOPATH/src/app
COPY . $GOPATH/src/app
# Verify OCR accuracy on standard captcha/language samples; must pass.
RUN go test -v ./pkg/driver/ocr/ -run=TestRecognize

RUN cp tesseract.cfg /tesseract.cfg
RUN go build -o ocr_server ./cmd/ocr-server/

# ---- final stage ----
# Still use the Golang image to run tests inside the container if needed,
# since installing Tesseract on the dev computer is complicated
FROM golang:1.26

RUN apt-get -qy update
RUN apt-get install -qy libicu76
RUN apt-get install -qy libpango-1.0-0 libpangocairo-1.0-0 libcairo2
RUN apt-get install -qy libleptonica-dev
RUN rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/lib/libtesseract.so* /usr/local/lib/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libleptonica.so.6* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/local/share/tessdata /usr/local/share/tessdata
COPY --from=builder /go/src/app/ocr_server /ocr_server
COPY --from=builder /go/src/app/tesseract.cfg /tesseract.cfg

RUN ldconfig

ENV TESSDATA_PREFIX=/usr/local/share/tessdata
ENV PORT=35735
WORKDIR /
CMD ["/ocr_server"]
