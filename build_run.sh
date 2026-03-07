set -e

docker build --tag=daominah/ocr_server:v5.5.2 --tag=daominah/ocr_server:latest .
# if you want to push to docker hub, run this:
# docker push daominah/ocr_server:v5.5.2
# docker push daominah/ocr_server:latest

docker rm -f ocr_server
docker run -dit --restart always --name=ocr_server -p=35735:35735 daominah/ocr_server:v5.5.2

set +e

# To run tests interactively in the build environment:
#   docker build --target builder --tag=ocr_server_builder .
#   docker run --rm -dit --workdir //go/src/app --name=ocr_server_builder ocr_server_builder bash
