
#### build swagger doc:
    go get -u github.com/swaggo/swag/cmd/swag
    swag init -d ./src/app/router -g router.go -o ./src/app/docs

#### Build Docker Image:
    cd ./
    docker build -t cat-api .
    
#### Copy .env.example as .env file ####
    cp .env.example .env    
    
## run Docker Image

#### use Docker run
    docker run -i --env-file=.env.example -p 8085:8085 -t cat-api:0.0.0

#### Use Docker Compose
    change env file POSTGRES_HOST=db
    docker network create app_net
    docker-compose up -d