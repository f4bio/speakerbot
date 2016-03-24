all:
	if [ ! -d "./vendor" ]; then make bootstrap; fi
	go build -o speakerbot *.go

bootstrap:
	which glide || go get github.com/Masterminds/glide && glide install -s

dev:
	which CompileDaemon || go get github.com/githubnemo/CompileDaemon && CompileDaemon -directory=. -exclude-dir=.git -exclude-dir=vendor -exclude=speakerbot -command=./speakerbot

docker:
	docker build -t dustinblackman/speakerbot .

test:
	./scripts/test.sh
