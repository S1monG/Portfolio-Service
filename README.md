Serverless cloud functions handling requests made by my (portfolio) gh-pages website. Handlers are written in go and deployed to GCP.

Authorize to gcloud and run below script to deploy service to GCP

`gcloud functions deploy <path> \
  --region=europe-west1
  --runtime=go122 \
  --trigger-http \
  --allow-unauthenticated \
  --entry-point=<Function Name> \
  --source=.

Right now im deploying the source code. Thats fine but there are some advantages with deploying the binary instead. However, if I want to do that I would have to use Cloud Run in place of Cloud Functions.

Deploying Binary:
Build a docker image which contains the binary and a dockerfile that starts it, push the image to Google Container Registery and deploy the container to Google Cloud Run. Way overcomplicated for a small project like this tho.

Should add tests once the functions goes into production, nah