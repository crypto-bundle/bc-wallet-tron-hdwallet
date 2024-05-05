default: build_plugin

build_plugin:
	$(eval short_commit_id=$(shell git rev-parse --short HEAD))
	$(eval commit_id=$(shell git rev-parse HEAD))
	$(eval build_number=0)
	$(eval build_date=$(shell date +%s))
	$(eval release_tag=$(shell git describe --tags $(commit_id))-$(short_commit_id)-$(build_number))

	CGO_ENABLED=1 go build -trimpath -race -installsuffix cgo -gcflags all=-N \
		-ldflags "-linkmode external -extldflags -w \
			-X 'main.BuildDateTS=${BUILD_DATE_TS}' \
			-X 'main.BuildNumber=${BUILD_NUMBER}' \
			-X 'main.ReleaseTag=${RELEASE_TAG}' \
			-X 'main.CommitID=${COMMIT_ID}' \
			-X 'main.ShortCommitID=${SHORT_COMMIT_ID}'" \
		-buildmode=plugin \
		-o ./build/tron.so \
		./plugin

deploy:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval controller_container_registry=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-controller)
	$(eval parent_container_registry=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-api)
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-tron-hdwalleta-api)
	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	$(eval short_commit_id=$(shell git rev-parse --short HEAD))
	$(eval commit_id=$(shell git rev-parse HEAD))
	$(eval build_number=0)
	$(eval build_date=$(shell date +%s))
	$(eval release_tag=$(shell git describe --tags $(commit_id))-$(short_commit_id)-$(build_number))

	docker build \
		--ssh default=$(SSH_AUTH_SOCK) \
		--platform $(platform) \
		--build-arg PARENT_CONTAINER_IMAGE_NAME=$(parent_container_registry):latest \
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg COMMIT_ID=$(commit_id) \
		--build-arg SHORT_COMMIT_ID=$(short_commit_id) \
		--build-arg BUILD_NUMBER=$(build_number) \
		--build-arg BUILD_DATE_TS=$(build_date) \
		--tag $(container_registry):$(build_tag) .

	docker push $(container_registry):$(build_tag)

	helm --kube-context $(context) upgrade \
		--install bc-wallet-tron-hdwallet-api \
		--set "global.api_container_path=$(container_registry)" \
		--set "global.controller_container_path=$(controller_container_registry)" \
		--set "global.build_tag=$(build_tag)" \
		--set "global.env=$(env)" \
		--values=./deploy/helm/hdwallet/values.yaml \
		--values=./deploy/helm/hdwallet/values_$(env).yaml \
		./deploy/helm/api

.PHONY: hdwallet_proto deploy
