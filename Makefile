# Install plugins:
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

plugin_proto:
	protoc -I ./plugins/tron/ -I . \
    		--go_out=./plugins/tron/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./plugins/tron/ \
    		--go-grpc_opt=paths=source_relative \
    		--openapiv2_out=logtostderr=true:./plugins/tron/ \
    		--grpc-gateway_out=./plugins/tron/ \
    		--grpc-gateway_opt=logtostderr=true \
    		--grpc-gateway_opt=paths=source_relative \
    		--doc_out=./plugins/tron/ \
    		--doc_opt=markdown,$@.md \
    		./plugins/tron/*.proto

default: hdwallet

plugin: plugin_proto
	$(eval short_commit_id=$(shell git rev-parse --short HEAD))
	$(eval commit_id=$(shell git rev-parse HEAD))
	$(eval build_number=0)
	$(eval build_date=$(shell date +%s))
	$(eval release_tag=$(shell git describe --tags $(commit_id))-$(short_commit_id)-$(build_number))

	CGO_ENABLED=1 go build -race -installsuffix cgo -gcflags all=-N \
		-ldflags "-linkmode external -extldflags -w \
			-X 'main.BuildDateTS=${BUILD_DATE_TS}' \
			-X 'main.BuildNumber=${BUILD_NUMBER}' \
			-X 'main.ReleaseTag=${RELEASE_TAG}' \
			-X 'main.CommitID=${COMMIT_ID}' \
			-X 'main.ShortCommitID=${SHORT_COMMIT_ID}'" \
		-buildmode=plugin \
		-o ../../build/tron.so \
		./plugins/tron

deploy:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-tron-hdwallet)
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
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg COMMIT_ID=$(commit_id) \
		--build-arg SHORT_COMMIT_ID=$(short_commit_id) \
		--build-arg BUILD_NUMBER=$(build_number) \
		--build-arg BUILD_DATE_TS=$(build_date) \
		--tag $(container_registry):$(build_tag) .

	docker push $(container_registry):$(build_tag)

	helm --kube-context $(context) upgrade \
		--install bc-wallet-tron-hdwallet-api \
		--set "global.container_registry=$(container_registry)" \
		--set "global.build_tag=$(build_tag)" \
		--set "global.env=$(env)" \
		--values=./deploy/helm/hdwallet/values.yaml \
		--values=./deploy/helm/hdwallet/values_$(env).yaml \
		./deploy/helm/api

.PHONY: hdwallet_proto deploy
