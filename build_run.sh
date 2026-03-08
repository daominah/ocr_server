set -e

docker build --tag=daominah/ocr_server:v5.5.2 --tag=daominah/ocr_server:latest .
# if you want to push to docker hub, run this:
# docker push daominah/ocr_server:v5.5.2
# docker push daominah/ocr_server:latest

docker rm -f ocr_server
docker run -dit --restart always --name=ocr_server -p=35735:35735 daominah/ocr_server:v5.5.2

# To run tests interactively in the build environment:
docker build --target builder --tag=ocr_server_builder .

docker rm -f ocr_run_go_test
docker run -dit --workdir //go/src/app --name=ocr_run_go_test ocr_server_builder bash -c 'echo "try command:
ERODE_RADIUS=1 go test -v ./pkg/driver/ocr/ -run=TestRecognizeOverlapChars" && sleep infinity'

# Note: building Tesseract from source is very time consuming,
# do not prune Docker if you are still working on this repo (prune clears cache).

set +e
