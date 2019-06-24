
#### Build Docker Image:
    cd ./
    docker build -t cat-api:0.0.0

#### run Docker Image
## windows escape character is ` linux,mac is \
    docker run -i --env-file=.env.example -p 8080:8080 -t cat-api:0.0.0

#### Copy .env.example as .env file ####
    cp .env.example .env