all: under

under:
	go build

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

clean:
	rm -rf under_pressure
