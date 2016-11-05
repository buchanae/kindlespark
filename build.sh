mkdir -p builds
export GOPATH=`pwd`/gopackages/

function build() {
  go build -o builds/kindlespark-$GOOS-$GOARCH kindlespark.go
}

export GOOS=linux
export GOARCH=amd64
build

export GOOS=windows
export GOARCH=amd64
build

export GOOS=windows
export GOARCH=386
build

export GOOS=darwin
export GOARCH=amd64
build
