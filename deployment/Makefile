PROJECT_ID=$(shell gcloud config get-value project)
COMMIT_ID=$(shell git rev-parse --short HEAD)
PROJECT_ROOT=$(shell git rev-parse --show-toplevel)

all: build

build:
	@gcloud container builds submit --config=cloudbuild.yml --substitutions _COMMIT_ID=$(COMMIT_ID) ..
clean:
	@gcloud container images list-tags gcr.io/$(PROJECT_ID)/yamato \
		--format="value(digest)" \
		--filter="tags.len()=0" | \
		 xargs -n1 -IDIGEST \
		 gcloud container images delete gcr.io/$(PROJECT_ID)/yamato@sha256:DIGEST -q
deploy:
	@gcloud beta compute instances create-with-container bot-server-$(COMMIT_ID) \
		--zone us-west1-c \
		--boot-disk-size 10GB \
		--machine-type=f1-micro \
		--container-image=gcr.io/$(PROJECT_ID)/yamato:$(COMMIT_ID) \
		--container-env-file=$(PROJECT_ROOT)/compose/.env.container
.PHONY: all
