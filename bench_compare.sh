if [ "${TRAVIS_PULL_REQUEST_BRANCH:-$TRAVIS_BRANCH}" != "master" ]; then
	REMOTE_URL="$(git config --get remote.origin.url)";
	cd ${TRAVIS_BUILD_DIR}/.. && \
	git clone ${REMOTE_URL} "${TRAVIS_REPO_SLUG}-bench" && \
	cd "${TRAVIS_REPO_SLUG}-bench" && \
	git checkout ${TRAVIS_COMMIT} && \
	go test -bench . ./... > current_bench && \
	git checkout master && \
	go test -bench . ./... > master_bench && \
	go get golang.org/x/tools/cmd/benchcmp && \
	benchcmp master_bench current_bench;
fi
