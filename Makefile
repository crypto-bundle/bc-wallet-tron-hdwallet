default: build_plugin

build_plugin:
	$(eval SHORT_COMMIT_ID=$(shell git rev-parse --short HEAD))
	$(eval COMMIT_ID=$(shell git rev-parse HEAD))
	$(eval BUILD_NUMBER=0)
	$(eval BUILD_DATE_TS=$(shell date +%s))
	$(eval RELEASE_TAG=$(shell git describe --tags $(COMMIT_ID))-$(SHORT_COMMIT_ID)-$(BUILD_NUMBER))

	CGO_ENABLED=1 go build -trimpath -race -installsuffix cgo -gcflags all=-N \
		-ldflags "-linkmode external -extldflags -w -s \
			-X 'main.NetworkChainID=195' \
			-X 'main.BuildDateTS=${BUILD_DATE_TS}' \
			-X 'main.BuildNumber=${BUILD_NUMBER}' \
			-X 'main.ReleaseTag=${RELEASE_TAG}' \
			-X 'main.CommitID=${COMMIT_ID}' \
			-X 'main.ShortCommitID=${SHORT_COMMIT_ID}'" \
		-buildmode=plugin \
		-o ./build/tron.so \
		./plugin

test_plugin:
	CGO_ENABLED=1 go build -trimpath -race -trimpath -installsuffix cgo \
		-gcflags all=-N \
		-o ./build/loader_test \
		-ldflags "-linkmode external -extldflags -w -s" \
		./cmd/loader_test

	./build/loader_test

deploy:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval migrator_container_path=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-migrator)
	$(eval controller_container_path=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-controller)
	$(eval parent_api_container_path=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-api)
	$(eval target_container_path=$(repository)/crypto-bundle/bc-wallet-tron-hdwallet-api)
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
		--build-arg RACE= \
		--build-arg PARENT_CONTAINER_IMAGE_NAME=$(parent_api_container_path):latest \
		--build-arg NETWORK_CHAIN_ID=195 \
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg COMMIT_ID=$(commit_id) \
		--build-arg SHORT_COMMIT_ID=$(short_commit_id) \
		--build-arg BUILD_NUMBER=$(build_number) \
		--build-arg BUILD_DATE_TS=$(build_date) \
		--tag $(target_container_path):$(build_tag) \
		--tag $(target_container_path):latest .

	docker push $(target_container_path):$(build_tag)
	docker push $(target_container_path):latest

	helm --kube-context $(context) upgrade \
		--install bc-wallet-tron-hdwallet \
		--set "global.migrator.image.path=$(migrator_container_path)" \
		--set "global.migrator.image.tag=latest" \
		--set "global.api.image.path=$(target_container_path)" \
		--set "global.api.image.tag=$(build_tag)" \
		--set "global.controller.image.path=$(controller_container_path)" \
		--set "global.controller.image.tag=latest" \
		--set "global.env=$(env)" \
		--values=./deploy/helm/hdwallet/values.yaml \
		--values=./deploy/helm/hdwallet/values_$(env).yaml \
		./deploy/helm/hdwallet

.PHONY: hdwallet_proto deploy