if [ "${TRAVIS_PULL_REQUEST_BRANCH:-$TRAVIS_BRANCH}" != "master" ]; then
	go test -bench . ./... > current_bench && \
	git checkout master && \
	go test -bench . ./... > master_bench && \
	go get golang.org/x/tools/cmd/benchcmp && \
	benchcmp master_bench current_bench;
fi
