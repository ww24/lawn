PROJECT_ID ?=

.PHONY: login
login:
	gcloud auth application-default login --project ${PROJECT_ID}

.PHONY: run
run:
	go run ./cmd/lawn
