all: verify_tx

verify_tx:
	go build -o ./cmd/verify_tx/verify_tx ./cmd/verify_tx/...

clean:
	rm -f ./cmd/verify_tx/verify_tx
