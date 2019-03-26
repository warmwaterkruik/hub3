PROJECT_ID=hub3-saas-nl-prod
SOURCE_BUCKET_NAME="websrc-${PROJECT_ID}"
DERIVATIVE_BUCKET_NAME="webdrv-${PROJECT_ID}"


deploy-update:
	#go mod vendor
	gcloud alpha functions deploy cloud-storage-update \
		--entry-point CloudUpdate \
		--memory 128MB \
		--region europe-west1 \
		--set-env-vars SOURCE_BUCKET_NAME=${SOURCE_BUCKET_NAME},DERIVATIVE_BUCKET_NAME=${DERIVATIVE_BUCKET_NAME},PROJECT_ID=${PROJECT_ID} \
		--runtime go111 \
		--trigger-event google.storage.object.finalize \
		--trigger-resource ${SOURCE_BUCKET_NAME}

deploy-delete:
	#go mod vendor
	gcloud alpha functions deploy cloud-storage-delete \
		--entry-point CloudDelete \
		--memory 128MB \
		--region europe-west1 \
		--set-env-vars SOURCE_BUCKET_NAME=${SOURCE_BUCKET_NAME},DERIVATIVE_BUCKET_NAME=${DERIVATIVE_BUCKET_NAME},PROJECT_ID=${PROJECT_ID} \
		--runtime go111 \
		--trigger-event google.storage.object.delete \
		--trigger-resource ${SOURCE_BUCKET_NAME}

deploy-resolve:
	#go mod vendor
	gcloud alpha functions deploy cloud-storage-resolve \
		--entry-point CloudResolve \
		--memory 128MB \
		--region europe-west1 \
		--set-env-vars SOURCE_BUCKET_NAME=${SOURCE_BUCKET_NAME},DERIVATIVE_BUCKET_NAME=${DERIVATIVE_BUCKET_NAME},PROJECT_ID=${PROJECT_ID} \
		--runtime go111 \
		--trigger-http

test:
	go test .
