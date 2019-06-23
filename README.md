cat-api

#### Build Docker Image:
cd ./
docker build -t cat-api:0.0.0

#### run Docker Image
## windows escape character is ` linux,mac is \
docker run -i `
-e POSTGRES_HOST='192.168.99.100' `
-e POSTGRES_PORT='5432' `
-e POSTGRES_USER='postgres' `
-e POSTGRES_PASSWORD='password' `
-p 8081:8080 `
-t cat-api:0.0.0