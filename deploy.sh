gcloud functions deploy postalcode-search \
--entry-point PostalCodeSearch \
--runtime go116 \
--trigger-http \
--allow-unauthenticated \
--project postalcode-firebase \
--region asia-northeast1