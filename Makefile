# Install plugins:
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
#  go get -d github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
#  go get -d github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

hdwallet_proto:
	protoc -I ./pkg/proto/ -I . -I ./pkg/proto/ \
    		--go_out=./pkg/grpc/hdwallet_api/proto/ \
    		--go_opt=paths=source_relative \
    		--go-grpc_out=./pkg/grpc/hdwallet_api/proto/ \
    		--go-grpc_opt=paths=source_relative \
    		--openapiv2_out=logtostderr=true:./docs/hdwallet_api/ \
    		--grpc-gateway_out=./pkg/grpc/hdwallet_api/proto/ \
    		--grpc-gateway_opt=logtostderr=true \
    		--grpc-gateway_opt=paths=source_relative \
    		--doc_out=./docs/hdwallet_api/ \
    		--doc_opt=markdown,$@.md \
    		./pkg/proto/*.proto

default: hdwallet

plugin:
	CGO_ENABLED=1 go build -race -v -installsuffix cgo -o ./build/api -ldflags "-linkmode external -extldflags -w" ./cmd/hdwallet_api
	CGO_ENABLED=1 go build -race -v -installsuffix cgo -o ./build/tron.so -ldflags "-linkmode external -extldflags -w"  -buildmode=plugin ./plugins/tron

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
