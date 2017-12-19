.PHONY: all

PROJECT_ID=$(shell gcloud config get-value project)
COMMIT_ID=$(shell git rev-parse --short HEAD)

all: build

build:
	@gcloud container builds submit --config=cloudbuild.yml --substitutions _COMMIT_ID=$(COMMIT_ID) ..
clean:
	@gcloud container images list-tags gcr.io/$(PROJECT_ID)/yamato \
		--format="value(digest)" \
		--filter="tags.len()=0" | \
		 xargs -n1 -IDIGEST \
		 gcloud container images delete gcr.io/$(PROJECT_ID)/yamato@sha256:DIGEST -q
