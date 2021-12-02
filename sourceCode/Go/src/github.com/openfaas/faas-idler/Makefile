TAG?=latest-dev
.PHONY: build
build:
	docker build -t openfaas/faas-idler:${TAG} .

.PHONY: push
push:
	docker push openfaas/faas-idler:${TAG}

.PHONY: ci-armhf-build
ci-armhf-build:
	docker build -t openfaas/faas-idler:${TAG}-armhf . -f Dockerfile.armhf

.PHONY: ci-armhf-push
ci-armhf-push:
	docker push openfaas/faas-idler:${TAG}-armhf

.PHONY: ci-arm64-build
ci-arm64-build:
	docker build -t openfaas/faas-idler:${TAG}-arm64 . -f Dockerfile.arm64

.PHONY: ci-arm64-push
ci-arm64-push:
	docker push openfaas/faas-idler:${TAG}-arm64

.PHONY: ci-ppc64le-build
ci-ppc64le-build:
	docker build -t openfaas/faas-idler:${TAG}-ppc64le . -f Dockerfile.ppc64le

.PHONY: ci-ppc64le-push
ci-ppc64le-push:
	docker push openfaas/faas-idler:${TAG}-ppc64le
