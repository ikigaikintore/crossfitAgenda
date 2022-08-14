FROM debian:11 as base

RUN apt-get update && \
    apt-get install -y apt-transport-https wget gnupg lsb-release && \
    echo "deb https://notesalexp.org/tesseract-ocr5/$(lsb_release -cs)/ $(lsb_release -cs) main"\
    | tee /etc/apt/sources.list.d/notesalexp.list > /dev/null && \
    wget -O - https://notesalexp.org/debian/alexp_key.asc | apt-key add - && \
    apt-get update && \
    apt-get install -y tesseract-ocr libtesseract-dev libleptonica-dev

FROM debian:11 as ocr
COPY --from=base /usr/include/leptonica /usr/include/leptonica
COPY --from=base /usr/bin/tesseract /usr/bin/tesseract
COPY --from=base /usr/include/tesseract /usr/include/tesseract
COPY --from=base /usr/share/tesseract-ocr/5/tessdata /usr/share/tesseract-ocr/5/tessdata

ENTRYPOINT ["/usr/bin/tesseract"]
