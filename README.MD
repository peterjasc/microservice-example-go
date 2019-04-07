# Micro service example Go app

## API calls to app
1. For specific Recipe IDs:
`/recipes?ids=1,2,3`
2. For all results, paginated at max by 10 by default
`recipes?skip=2&top=2`
skip - the number of recipes to skip
top - number of recipes to get (max 10)

## Running locally

1. `go build -o bin/main cmd/main.go && ./bin/main`
2. `curl "localhost:8080/recipes?ids=1,2,30"`

## Running on Docker 

1. `docker build -t recipeapp/peterjasc . && docker run -it recipeapp/peterjasc`
2. Find the container ID with `docker ps`
3. Find the IP of the container `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' [container ID]`
4. Send a request to the app `curl "172.17.0.2:8080/recipes?skip=2&top=2"` or `curl "172.17.0.2:8080/recipes?ids=1,2,3"`
