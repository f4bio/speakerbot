set -e
cd $PWD

go vet

if [[ $(golint *.go) ]]; then
  golint *.go
  echo "golint failed"
  exit 1
fi

if [[ $(gofmt -d ./*.go) ]]; then
  gofmt -d ./*.go
  echo "gofmt returned suggested changes, please run gofmt first. Exiting..."
  exit 1
fi

echo "Hooray! Tests passed."
