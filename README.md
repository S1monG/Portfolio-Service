WIP Serverless cloud functions handling requests made by my (portfolio) gh-pages website. Handlers are written in go and deployed to GCP.

Authorize to gcloud and run below script to deploy service to GCP

`gcloud functions deploy <path> \
  --region=europe-west1
  --runtime=go122 \
  --trigger-http \
  --allow-unauthenticated \
  --entry-point=<Function Name> \
  --source=.


Alternativly for more control in the future I should build the functions in a docker container image, push the image to Google Container Registery and deploy the container to Google Cloud Run

Should add tests once the functions goes into production