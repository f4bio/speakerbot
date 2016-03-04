all:
	which glide || go get github.com/Masterminds/glide && glide install

build:
	go build -o speakerbot *.go

dev:
	which CompileDaemon || go get github.com/githubnemo/CompileDaemon && CompileDaemon -directory=. -exclude-dir=.git -exclude-dir=vendor -exclude=speakerbot -command=./speakerbot

docker-build:
	docker build -t dustinblackman/speakerbot .

test:
	./scripts/test.sh
