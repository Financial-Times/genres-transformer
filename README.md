# genres-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/genres-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/genres-transformer/tree/master)

Retrieves Genres taxonomy from TME vie the structure service and transforms the genres to the internal UP json model.
The service exposes endpoints for getting all the genres and for getting genre by uuid.

# Usage
`go get github.com/Financial-Times/genres-transformer`

`$GOPATH/bin/genres-transformer --port=8080 -base-url="http://localhost:8080/transformers/genres/" -structure-service-base-url="http://metadata.internal.ft.com:83" -structure-service-username="user" -structure-service-password="pass" -structure-service-principal-header="app-preditor"`
```
export|set PORT=8080
export|set BASE_URL="http://localhost:8080/transformers/genres/"
export|set STRUCTURE_SERVICE_BASE_URL="http://metadata.internal.ft.com:83"
export|set STRUCTURE_SERVICE_USERNAME="user"
export|set STRUCTURE_SERVICE_PASSWORD="pass"
export|set PRINCIPAL_HEADER="app-preditor"
$GOPATH/bin/genres-transformer
```

With Docker:

`docker build -t coco/genres-transformer .`

`docker run -ti --env BASE_URL=<base url> --env STRUCTURE_SERVICE_BASE_URL=<structure service url> --env STRUCTURE_SERVICE_USERNAME=<user> --env STRUCTURE_SERVICE_PASSWORD=<pass> --env PRINCIPAL_HEADER=<header> coco/genres-transformer`

# Deployment status
Because there were just a few genres, this was not deployer in any environment and the genre data import was done from a local machine.
That's why this app doesn't have service files.