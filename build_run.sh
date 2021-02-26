docker build --tag=daominah/ocr_server . &&\
    docker rm -f ocr_server ;\
    docker run -dit --restart always --name=ocr_server -p=35735:35735 daominah/ocr_server
