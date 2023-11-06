# base image for golang
FROM golang:latest
MAINTAINER Rohan Nagavardhan (rnagavar@umich.edu)

# define working directory
# equivalent of `RUN mkdir -p <desired path>`
WORKDIR /usr/src/backend

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/

# Load languages models.
RUN apt-get install -y -qq tesseract-ocr-eng

# pre-copy/cache golang dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy source code
COPY . .
RUN make build # build the application

EXPOSE 8080

ENTRYPOINT ["./bin/main.exe"]
