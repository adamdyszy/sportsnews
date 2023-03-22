[![GoDoc](http://godoc.org/github.com/adamdyszy/sportsnews?status.png)](http://godoc.org/github.com/adamdyszy/sportsnews)

# Sports News

Sports News tool for collecting data from news feeds with list and details,
and then exposing that information through API.

## How it works

- It runs cron scheduled news poller that will
  - Poll list of N newest newses from specified news list URL
  - Save them as articles into storage and mark new ones as articles without details
  - Poll details of articles from specified news details URL
- Serve http router that will handle rest requests:
  - GET at "/articles" path, return all articles in json
  - GET at "/articles/{id}" path, where {id} is article id available when querying articles list, return article in json
  - You can see returned structures at [types/article.go](types/article.go)
- Application will continue to serve API and run cron jobs indefinitely even if database will fail at some point,
but if it will start working again at some point then it should work again if same connection details.

## Quickstart

Quick examples to run using docker.
The container names has hash in their names so they are unique.
It is generated using git commit or if not possible then md5sum of go files excluding vendor.

### Quickstart MongoDB

Example using docker and mongo.
It will run 2 containers: mongodb, newsServer.
NewsServer will run with this config: [config/quickstart/mongo.yaml](config/quickstart/mongo.yaml)

- Start with:

```bash
make quickstart
```

- Then you should be able to access the app within a minute with curl or accessing the web:

```bash
curl localhost:8080/articles
```

- You can choose some id from the result and get single article using:

```bash
curl localhost:8080/articles/{ID}
```

- Since the detailed and list poller started at the same time detailed poller didn't know what details to get.
- It will get the details on the next scheduled time (3 minutes in this example), but if you want to force it you can just restart the server:

```bash
make quickstart-restart-server
```

- To kill and delete these containers run:

```bash
make quickstart-kill
```

### Quickstart Memory

It will run container: newsServer.
NewsServer will run with this config: [config/quickstart/memory.yaml](config/quickstart/memory.yaml)

- Start with:

```bash
make quickstart-mem
```

- Then you should be able to access the app within a minute with curl or accessing the web:

```bash
curl localhost:8080/articles
```

- You can choose some id from the result and get single article using:

```bash
curl localhost:8080/articles/{ID}
```

- To kill and delete these containers run:

```bash
make quickstart-mem-kill
```

## How to run

- Build the binary using cached dependencies:

```bash
go build -o=bin/sportsnews -mod=vendor cmd/sportsnews/main.go
```

- Create [config/custom.yaml](config/custom.yaml) file.
- Add settings based on [config/default.yaml](config/default.yaml) where all fields are described.
- Options at config/custom.yaml will override everything in config/default.yaml
- Run with config checks using:

```bash
if test -r config/default.yaml && test -r config/custom.yaml;then ./bin/sportsnews;fi
```

- You can run binary using different config with:

```bash
./bin/sportsnews --customConfigFile <PATH to your config file>
```

- You can just run it using makefile:

```bash
make run
```

or with custom config file:

```bash
make run CONFIG_FILE=<PATH to your config file>
```

- You can also build and run the docker image for the server, but before
that make sure you have your config ready in config/custom.yaml or set CONFIG_FILE.

```bash
make docker-build
make docker-run #CONFIG_FILE=<PATH to your config file>
# or
make docker-run-background #CONFIG_FILE=<PATH to your config file>
```

- You can get config of your running server in docker using:

```bash
docker logs <your_server_docker_container_name> 2>&1 | head -n 1
```

```
2023-03-22T14:55:21Z    INFO    poller/poller.go:34     Starting poller with this config.       {"workerKind": "NewsPoller", "config": {"TeamId":"t94","RunOnceAtBoot":true,"List":{"URL":"https://www.wearehullcity.co.uk/api/incrowd/getnewlistinformation","Count":100,"Schedule":"@every 1h"},"Details":{"URL":"https://www.wearehullcity.co.uk/api/incrowd/getnewsarticleinformation","Schedule":"@every 5m"}}}
```

## MongoDB configuration

- Check mongoDB documentation here for how to get your database running https://docs.mongodb.com
- Ensure you configure all fields, including user and password in your config file.

```yaml
mongoStorage: # options for storageKind mongo
  uri: "mongodb://localhost:27017" # connection uri
  name: "newsDB" # database name
  articlesColl: "articles" # articles collection name
  user: "mongoadmin" # username when connecting to db
  password: "secret" # password when connecting to db
  timeoutSeconds: 60 # how many seconds should db wait for execution of queries before cancellation
```

- If you want to use auth inside the uri without providing username and password do it like that:

```yaml
mongoStorage:
  uri: "mongodb://mongoadmin:secret@localhost:27017" # connection uri with auth
  user: "" # empty username means the auth is inside uri
```

- With these setting when poller pulls data into the storage you can see them using:

```bash
mongosh mongodb://mongoadmin:secret@localhost:27017 # connection uri with auth
```

```rust
use newsDB // database name
db.articles.find() // shows all articles
db.articles.find({hasDetails: false}) // shows all articles without details
db.articles.find({hasDetails: true}) // shows all articles with details
db.articles.find({newsId: "444478"}) // shows article that was pulled from news with id 444478
db.articles.find({id: "a6a25a18-3bb2-5444-b2af-4c4ffa872110"}) // shows article filtering with its id
db.articles.countDocuments({hasDetails: false}) // how many articles without details
// db.articles.deleteMany({}) // delete all articles
// db.articles.deleteMany({id: "a6a25a18-3bb2-5444-b2af-4c4ffa872110"}) // delete articles with given id
```

## Run tests

```bash
make test
```

## What does not work

- There is no pruning data built into the system
- Getting more than 100 latest news form hullcity was not possible
- There is no sorting of articles list when returning it
- When article has no details it shows in field hasDetails, but not in response status code

## TODO

- Add more tests
- Add more descriptions for godoc
- Add swagger json
- Add travis for CI/CD