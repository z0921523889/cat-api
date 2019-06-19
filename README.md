cat-api


#### Build Docker Image:
cd ./
docker build -t cat-api:0.0.0
docker run -i -t -p 8081:8080 cat-api:0.0.0
