ARG PARENT_CONTAINER_IMAGE_NAME="/crypto-bundle/bc-wallet-common-hdwallet-api:latest"

FROM golang:1.22.2-bookworm AS gobuild

ENV GO111MODULE on
ENV GOSUMDB off
# add go-base repo to exceptions as a private repository.
ENV GOPRIVATE $GOPRIVATE,github.com/crypto-bundle

# add private github token
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		git \
        openssh-client \
        build-essential \
	; \
	rm -rf /var/lib/apt/lists/* && \
    mkdir -p -m 0700 ~/.ssh && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    git config --global url."git@github.com".insteadOf "https://github.com/"

WORKDIR /src

# Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever).
COPY go.* ./

COPY . .

# Compile! Should only compile our sources since everything else is precompiled.
ARG RACE=-race
ARG CGO=1
ARG NETWORK_CHAIN_ID=195
ARG RELEASE_TAG="v0.0.0-00000000-100500"
ARG COMMIT_ID="0000000000000000000000000000000000000000"
ARG SHORT_COMMIT_ID="0000000"
ARG BUILD_NUMBER="100500"
ARG BUILD_DATE_TS="1713280105"
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /src/bin && \
    GOOS=linux CGO_ENABLED=${CGO} go build ${RACE} \
        -gcflags all=-N \
        -ldflags "-linkmode external -extldflags -w \
            -X 'main.BuildDateTS=${BUILD_DATE_TS}' \
            -X 'main.BuildNumber=${BUILD_NUMBER}' \
            -X 'main.ReleaseTag=${RELEASE_TAG}' \
            -X 'main.CommitID=${COMMIT_ID}' \
            -X 'main.ShortCommitID=${SHORT_COMMIT_ID}'" \
        -buildmode=plugin \
        -o /src/bin/hdwallet_plugin_tron.so \
        ./plugin

FROM $PARENT_CONTAINER_IMAGE_NAME

ARG PLUGIN_ROOT="/usr/local/lib/"
ENV HDWALLET_PLUGIN_PATH="$PLUGIN_ROOT/hdwallet_plugin_tron.so"
ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$PLUGIN_ROOT"

COPY --from=gobuild /src/bin/hdwallet_plugin_tron.so $PLUGIN_ROOT