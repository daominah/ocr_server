# OCR server

Used for reading simple CAPTCHA. Powered by Tesseract 5.5.2

## API

### [/base64](http://127.0.0.1:35735/base64) POST

Example body:

````
{
    "base64": "iVBORw0KGgoAAAANSUhEUgAAABYAAAAkCAMAAAC62DqvAAAAP1BMVEUAAAAkJSgjKCgoKCglJSgjJSckJSgjJSkkJCYnJycnJycqKiokJSgjJSckJSgkJCclJSglJSklJSwaGhokJSjbbGjNAAAAFHRSTlMA8joT0HDMil0hGgbVtaWVUkQpCmqOj4cAAABQSURBVCjP5cg3DoAwEATAdcSBjP//VtCVvkWixlMOPqopFL0utmayWt8ek15puibTtT/cGtki0UWQnl3Xxcj7/u2gf/EPoL+B/gHlXJcdb24oNg3pSN9UAQAAAABJRU5ErkJggg==",
    "whitelist": "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}
````

base64 field can have prefix `data:image/png;base64,`

## Install

This project Dockerfile builds Tesseract from source, so you can choose a suitable version.

Convenient commands for dev:

````bash
docker build --tag=daominah/ocr_server .

docker rm -f ocr_server

docker run -dit --restart always --name=ocr_server -p=35735:35735 daominah/ocr_server
````

Published images on Docker Hub: `daominah/ocr_server:v5.5.2`, `daominah/ocr_server:v4.1.1`

## Tesseract trained data

The Dockerfile downloads trained data from
[tessdata_best](https://github.com/tesseract-ocr/tessdata_best) (highest accuracy, LSTM only).

There are three official trained data repositories — all essentially frozen as of early 2024:

| Repo                                                            | Engine                      | Last updated |
|-----------------------------------------------------------------|-----------------------------|--------------|
| [tessdata](https://github.com/tesseract-ocr/tessdata)           | Legacy + LSTM               | 2024-03-09   |
| [tessdata_best](https://github.com/tesseract-ocr/tessdata_best) | LSTM only                   | 2024-03-09   |
| [tessdata_fast](https://github.com/tesseract-ocr/tessdata_fast) | LSTM only (integer, faster) | 2024-08-01   |

Languages included: `eng`, `vie`, `chi_sim`.

## Config

Tesseract parameters can be changed to modify its behaviour
in [tesseract.cfg](./tesseract.cfg)

Doc: [Tesseract improve quality](https://github.com/tesseract-ocr/tessdoc/blob/master/ImproveQuality.md)

## Source

* This project is forked from [otiai10/ocrserver](https://github.com/otiai10/ocrserver)
* Go wrap library: [otiai10/gosseract](https://github.com/otiai10/gosseract)
* Origin Tesseract project in C++ [tesseract-ocr/tesseract](https://github.com/tesseract-ocr/tesseract)
