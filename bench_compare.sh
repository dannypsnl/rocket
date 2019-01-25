set -eu

if [ "${TRAVIS_PULL_REQUEST_BRANCH:-$TRAVIS_BRANCH}" != "master" ]; then
	go test -bench . ./... > current_bench.out
	git checkout master
	go test -bench . ./... > master_bench.out
	go get golang.org/x/tools/cmd/benchcmp
	benchcmp master_bench.out current_bench.out
fi
