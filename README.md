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

## How to run

- Build the binary using cached dependencies:

```bash
go build -o=bin/sportsnews -mod=vendor cmd/sportsnews/main.go
```

- Create [config/custom.yaml](config/custom.yaml) file.
- Add settings based on [config/default.yaml](config/default.yaml) where all fields are described.
- Options at config/custom.yaml will override everything in config/default.yaml
- Run using:

```bash
if test -r config/default.yaml && test -r config/custom.yaml;then ./bin/sportsnews;fi
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

- You can easily run example mongo database with docker using:

```bash
docker run -d -p 27017:27017 --name some-mongo \
	-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
	-e MONGO_INITDB_ROOT_PASSWORD=secret \
	mongo
```

## Run tests

```bash
go test ./...
```

## What does not work

- If hash of 2 different news is the same the poller will be confused and will treat them as same articles
- There is no pruning data built into the system
- Getting more than 100 latest news form hullcity was not possible
- There is no sorting of articles list when returning it
- There is no swagger json for the API
- When article has no details it shows in field hasDetails, but not in response status code