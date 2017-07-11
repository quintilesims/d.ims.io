test-small:
	go test -parallel 3 -count 1

test-large:
	go test -parallel 10 -count 3

test: test-small

.PHONY: test-small test-large test
