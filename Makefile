
IAMGE_REPO=d.caodong.xyz
PROJECT = scaffold
PKG=CDcoding2333/$(PROJECT)
VERSION=git-$(shell git describe --always --dirty)
IMAGE_TAG=$(VERSION)

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -a -tags netgo -installsuffix netgo -installsuffix cgo -ldflags '-w -s' -ldflags "-X main.Version=$(VERSION)" \
		-o ./build/linux/$(PROJECT) $(PKG)
		upx ./build/linux/$(PROJECT)

darwin:
	GOOS=darwin GOARCH=amd64 \
		go build -a -tags netgo -installsuffix netgo -ldflags "-X main.Version=$(VERSION)" \
		-o ./build/darwin/$(PROJECT) $(PKG)

ship:
	docker build -t ${IAMGE_REPO}/$(PROJECT):${IMAGE_TAG} .
	docker push ${IAMGE_REPO}/$(PROJECT):${IMAGE_TAG}

push:
	docker build -t ${IAMGE_REPO}/$(PROJECT):${IMAGE_TAG} .
	docker push ${IAMGE_REPO}/$(PROJECT):${IMAGE_TAG}
	docker rmi ${IAMGE_REPO}/$(PROJECT):${IMAGE_TAG}