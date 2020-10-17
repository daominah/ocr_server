# OCR server

Used for read simple CAPTCHA

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

I check some CAPTCHA and see Tesseract version [3.05.02](https://github.com/tesseract-ocr/tesseract/tree/3.05.02)
is better than [4.1.1](https://github.com/tesseract-ocr/tesseract/tree/4.1.1).
So I have to build old version from source instead of using apt.
See Dockerfile for detail.

## Config

Tesseract parameters can be changed to modify its behaviour
in [tesseract.cfg](./tesseract.cfg)

Doc: [Tesseract ControlParams](https://tesseract-ocr.github.io/tessdoc/ControlParams.html)

## Source

* This project is forked from [otiai10/ocrserver](https://github.com/otiai10/ocrserver)
* Go wrap library: [otiai10/gosseract](https://github.com/otiai10/gosseract)
* Origin Tesseract project in C++ [tesseract-ocr/tesseract](https://github.com/tesseract-ocr/tesseract)
