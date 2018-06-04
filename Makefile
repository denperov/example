git_branch ?= $(shell git rev-parse --abbrev-ref HEAD)

docker_repo  ?= denperov
docker_tag   ?= $(git_branch)

run: images
	docker_tag=${docker_tag} docker-compose up

images: image-front-api image-auth-api image-items-api image-offers-api

image-front-api:
	docker build \
		-t $(docker_repo)/owm-task-front-api:$(docker_tag) \
		-f ./services/front-api/Dockerfile .

image-auth-api:
	docker build \
		-t $(docker_repo)/owm-task-auth-api:$(docker_tag) \
		-f ./services/auth-api/Dockerfile .

image-items-api:
	docker build \
		-t $(docker_repo)/owm-task-items-api:$(docker_tag) \
		-f ./services/items-api/Dockerfile .

image-offers-api:
	docker build \
		-t $(docker_repo)/owm-task-offers-api:$(docker_tag) \
		-f ./services/offers-api/Dockerfile .
