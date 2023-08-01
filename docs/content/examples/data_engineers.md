---
title: "Ensign for Data Engineers"
weight: 20
date: 2023-05-17T17:03:41-04:00
---

We love data engineers &mdash; it's how a lot of us got our starts in tech. One of the main reasons we made Ensign is to make it easier for you to put your data in motion. We know that a clumsy ETL routine can quickly turn a data lake into a data landfill.

In this example we'll see how to move data around with Ensign. We'll be using [Watermill-Ensign](https://github.com/rotationalio/watermill-ensign) and [Watermill](https://watermill.io) to call a Weather API and insert weather data into a PostgreSQL database. If you haven't used Watermill yet, you're in for a treat! Check out this [introductory post](https://rotational.io/blog/prototyping-eda-with-watermill/) that covers the basics.

Just want the code? Check out this [repo](https://github.com/rotationalio/ensign-examples/tree/main/go/weather_data) for the full example.

## ETL Design

The architecture for this weather ingestor is composed of three components:
- An Ensign publisher that calls the Weather API and publishes the weather data to a topic.
- An Ensign subscriber that listens on this topic and runs a check against the PostgreSQL database to see if this is a new record. The weather data doesn't change that often, so it is possible to receive a duplicate record. If the record is new, we'll put the record into a second topic.
- A sql publisher that inserts the records from the second topic into the database.

The Ensign subscriber and the sql publisher are chained together using the `router` and `handler` functionality described in [this post](https://rotational.io/blog/prototyping-eda-with-watermill/).

## Prerequisites

This tutorial assumes that the following steps have been completed:
- You have installed **watermill**, **ensign**, **watermill-ensign**, and **watermill-sql**.
- You have received an Ensign Client ID and Client Secret.  Refer to the [getting started guide]({{< ref "/getting-started/ensign" >}}) on how to obtain the key.
- You have received an API key from the [Weather API website (it's free!)](https://www.weatherapi.com).
- You have [Docker](https://www.docker.com/) installed and running on your machine.

## Project Setup

First, you will need to set the environment variables for `ENSIGN_CLIENT_ID` and `ENSIGN_CLIENT_SECRET` from [your API Key]({{< ref "/getting-started/ensign#ensign-keys" >}}). ([Need a new key?](https://rotational.app/))

Next, let's create a root directory called `weather_data` for the application.

```bash
mkdir weather_data
```

We will then create two subfolders, one for the component that calls the Weather API to get the latest weather data and the other for the component that receives the data and inserts it into the database.

```bash
cd weather_data
mkdir producer
mkdir consumer
```
## Create the Ensign Publisher

Creating a publisher is very straightforward as you can see below.

```golang
publisher, err := ensign.NewPublisher(
		ensign.PublisherConfig{
			EnsureCreateTopic: true,
			Marshaler: ensign.EventMarshaler{},
		},
		logger,
	)
```

### Call the Weather API

Before we call the Weather API, we need to create the following structs:


First, a high level struct to represent the updates that come back from the Weather API.

```golang
type Response struct {
	Current Current `json:"current,omitempty"`
}
```

Next, a more detailed struct to help us parse all of the components of the Weather API's response. Some of this will depend on how much detail you need to ingest for your downstream data users (will the data scientists on your team complain if you forget to ingest the full text description provided by the response?)

```golang
type Current struct {
	LastUpdated string            `json:"last_updated,omitempty"`
	TempF       float64           `json:"temp_f,omitempty"`
	Condition   *CurrentCondition `json:"condition,omitempty"`
	WindMph     float64           `json:"wind_mph,omitempty"`
	WindDir     string            `json:"wind_dir,omitempty"`
	PrecipIn    float64           `json:"precip_in,omitempty"`
	Humidity    int32             `json:"humidity,omitempty"`
	FeelslikeF  float64           `json:"feelslike_f,omitempty"`
	VisMiles    float64           `json:"vis_miles,omitempty"`
}

type CurrentCondition struct {
	Text string `json:"text,omitempty"`
}
```

Finally, a struct to represent whatever structure makes the most sense for the weather data in *your* organization (e.g. with your company's database schemas or use cases in mind):

```golang
type ApiWeatherInfo struct {
	LastUpdated   string
	Temperature   float64
	FeelsLike     float64
	Humidity      int32
	Condition     string
	WindMph       float64
	WindDirection string
	Visibility    float64
	Precipitation float64
}
```

Here is the code to create the `request` object:

```golang
req, err := http.NewRequest("GET", "http://api.weatherapi.com/v1/current.json?", nil)
```

Next, we will define the query parameters and add it to `req`.  Note that you will need to create an environment variable called `WAPIKEY` that will be set to the API key you received from the Weather API.

```golang
q := req.URL.Query()
q.Add("key", os.Getenv("WAPIKEY"))
q.Add("q", "Washington DC")
req.URL.RawQuery = q.Encode()
```

Let's create the http client to call the Weather API, parse the `response` object, and create an `ApiWeatherInfo` object.

```golang
// create the client object
client := &http.Client{}

// retrieve the response
resp, err := client.Do(req)

// read the body of the response
body, _ := ioutil.ReadAll(resp.Body)

// unmarshal the body into a Response object
err = json.Unmarshal(body, &response)

// Convert the Response object into a ApiWeatherInfo object
current := response.Current
currentWeatherInfo := ApiWeatherInfo{
    LastUpdated:   current.LastUpdated,
    Temperature:   current.TempF,
    FeelsLike:     current.FeelslikeF,
    Humidity:      current.Humidity,
    Condition:     current.Condition.Text,
    WindMph:       current.WindMph,
    WindDirection: current.WindDir,
    Visibility:    current.VisMiles,
    Precipitation: current.PrecipIn,
}
```
Here is the complete function with some additional error handling.

<script src="https://gist.github.com/rebeccabilbro/1c8393527da818171a9b0e7f3f5ce871.js"></script>

### Publish the Data to a Topic

We'll create a helper method `publishWeatherData` that takes in a publisher and a channel that is used as a signal to stop publishing. First, create a ticker to call the Weather API every 5 seconds. Next, we will call the `GetCurrentWeather` function that we constructed previously to retrieve weather data, serialize it, construct a Watermill message, and publish the message to the `current_weather` topic.

```golang
func publishWeatherData(publisher message.Publisher, closeCh chan struct{}) {
	//weather doesn't change that often - call the Weather API every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		//if a signal has been sent through closeCh, publisher will stop publishing
		case <-closeCh:
			ticker.Stop()
			return

		case <-ticker.C:
		}

		//call the API to get the weather data
		weatherData, err := GetCurrentWeather()
		if err != nil {
			fmt.Println("Issue retrieving weather data: ", err)
			continue
		}

		//serialize the weather data
		payload, err := json.Marshal(weatherData)
		if err != nil {
			fmt.Println("Could not marshall weatherData: ", err)
			continue
		}

		//construct a watermill message
		msg := message.NewMessage(watermill.NewUUID(), payload)

		// Use a middleware to set the correlation ID, it's useful for debugging
		middleware.SetCorrelationID(watermill.NewShortUUID(), msg)

		//publish the message to the "current weather" topic
		err = publisher.Publish("current_weather", msg)
		if err != nil {
			fmt.Println("cannot publish message: ", err)
			continue
		}
	}
}
```
### Start the First Stream

Next we want to create a long-running process that will continue pinging the Weather API, parsing the response, and publishing it for downstream consumption.

In our `main()` function, we will first create a logger using `watermill.NewStdLogger`. Then we create the publisher and create the `closeCh` channel that will be used to send a signal to the publisher to stop publishing. We then pass the `publisher` and the `closeCh` to the `publishWeatherData` function and run it in a goroutine.

We then create another channel that listens for a `os.Interrupt` signal and will close when it receives the signal. If it receives the signal (e.g. because we want to stop the process, or if something goes wrong), it closes and the code moves on to close the `closeCh` channel and that notifies the publisher to stop publishing.

```golang
func main() {
	// add a logger
	logger := watermill.NewStdLogger(false, false)
	logger.Info("Starting the producer", watermill.LogFields{})

	//create the publisher
	publisher, err := ensign.NewPublisher(
		ensign.PublisherConfig{
			Marshaler: ensign.EventMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	//used to signal the publisher to stop publishing
	closeCh := make(chan struct{})

	go publishWeatherData(publisher, closeCh)

	// wait for SIGINT - this will end processing
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// signal for the publisher to stop publishing
	close(closeCh)

	logger.Info("All messages published", nil)
}
```

## Create the Ensign Subscriber

Next we'll create a subscriber in much the same way as we created our publisher.
Our subscriber will be in charge of listing for incoming weather messages from the publisher, and checking to see if the incoming data is actually new.

```golang
subscriber, err := ensign.NewSubscriber(
		ensign.SubscriberConfig{
			EnsureCreateTopic: true,
			Unmarshaler: ensign.EventMarshaler{},
		},
		logger,
	)
```
### Connect to the Database

Let's write a quick function `createPostgresConnection` that will allow us to connect to our database. You will need to create the following environment variables to connect to your local PostgreSQL database: `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB`. You can use any values of your choosing for these variables. For convenience, in this example, we'll use a docker PostgreSQL container. The function is below.

```golang
func createPostgresConnection() *stdSQL.DB {
	host := "weather_db"
	port := 5432
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname
	)
	db, err := stdSQL.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createQuery := `CREATE TABLE IF NOT EXISTS weather_info (
		id SERIAL NOT NULL PRIMARY KEY,
		last_updated VARCHAR(50) NOT NULL,
		temperature DECIMAL,
		feels_like DECIMAL,
		humidity INTEGER,
		condition VARCHAR(36),
		wind_mph DECIMAL,
		wind_direction VARCHAR(36),
		visibility DECIMAL,
		precipitation DECIMAL,
		created_at VARCHAR(100) NOT NULL
	);`
	_, err = db.ExecContext(context.Background(), createQuery)
	if err != nil {
		panic(err)
	}
	log.Println("created table weather_info")

	return db
}
```
Note that above, the host is defined as `weather_db` and will match the name of the container. More on that later in the tutorial!

We will also define a `WeatherInfo` struct that will contain the fields of the `weather_info` table. It is similar to the `ApiWeatherInfo` struct with the exception of the `CreatedAt`, which is an additional field in the table.

```golang
type WeatherInfo struct {
	LastUpdated   string
	Temperature   float64
	FeelsLike     float64
	Humidity      int32
	Condition     string
	WindMph       float64
	WindDirection string
	Visibility    float64
	Precipitation float64
	CreatedAt     string
}
```

### Does the Record Exist?

Next, we'll create a function that will check the database to see if the record already exists there. Before we create the function, we need to create a `dbHandler` struct which implements the `stdSQL.DB` interface that will be used to query the database.

```golang
type dbHandler struct {
	db *stdSQL.DB
}
```

We then write our `checkRecordExists` function, which will get executed when a message arrives on the `current_info` topic.

The first step is to unmarshal the message into a `ApiWeatherInfo` object. Next we will execute a query to see if a record with the `LastUpdated` value exists and if it does not, we will create a new `WeatherInfo` object and add the current timestamp for the `CreatedAt` field. We will then marshal this object, create, and return a new watermill message which will get published to the `weather_info` topic. If `LastUpdated` already exists, we simply note that the record exists and return `nil`.

```golang
func (d dbHandler) checkRecordExists(msg *message.Message) ([]*message.Message, error) {
	weatherInfo := ApiWeatherInfo{}

	err := json.Unmarshal(msg.Payload, &weatherInfo)
	if err != nil {
		return nil, err
	}

	log.Printf("received weather info: %+v", weatherInfo)

	var count int
	query := "SELECT count(*) FROM weather_info WHERE last_updated = $1"

	err = d.db.QueryRow(query, weatherInfo.LastUpdated).Scan(&count)
	switch {
	case err != nil:
		return nil, err
	default:
		if count > 0 {
			log.Println("Found existing record in the database")
			// not throwing an error here because this is not an issue
			return nil, nil
		}
		newWeatherInfo := WeatherInfo{
			LastUpdated:   weatherInfo.LastUpdated,
			Temperature:   weatherInfo.Temperature,
			FeelsLike:     weatherInfo.FeelsLike,
			Humidity:      weatherInfo.Humidity,
			Condition:     weatherInfo.Condition,
			WindMph:       weatherInfo.WindMph,
			WindDirection: weatherInfo.WindDirection,
			Visibility:    weatherInfo.Visibility,
			Precipitation: weatherInfo.Precipitation,
			CreatedAt:     time.Now().String(),
		}
		log.Println(newWeatherInfo)
		log.Println(len(newWeatherInfo.CreatedAt))
		var payload bytes.Buffer
		encoder := gob.NewEncoder(&payload)
		err := encoder.Encode(newWeatherInfo)
		if err != nil {
			panic(err)
		}

		newMessage := message.NewMessage(watermill.NewULID(), payload.Bytes())
		return []*message.Message{newMessage}, nil
	}
}
```

## Prepping the Database

Now we need to create a SQL publisher. A SQL publisher is a Watermill implementation of a SQL based pub/sub mechanism whereby you can use publishers to insert or upsert records and you can use subscribers to retrieve records. For more details, refer to this [post](https://watermill.io/pubsubs/sql/) on the Watermill site.

The SQL publisher is created as follows. Note that we are going to set `AutoInitializeSchema` to `false` because we've already created the table. The `postgresSchemaAdapter` is an extension of Watermill's `SchemaAdapter` which provides the schema-dependent queries and arguments.

```golang
pub, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        postgresSchemaAdapter{},
			AutoInitializeSchema: false,
		},
		logger,
	)
```

Then we'll create a `SchemaInitializingQueries` function to make a table based on the topic name, but since we've already created the table, we'll set this parameter to `false` and simply return an empty list.

```golang
func (p postgresSchemaAdapter) SchemaInitializingQueries(topic string) []string {
	return []string{}
}
```


### Inserting New Data

Next we will create an `InsertQuery` function that unmarshals the list of messages and creates an `insertQuery` sql statement that will be executed. The `topic` name is the same as the table name and it is used in the `insertQuery` sql statement. It then extracts the fields of the `WeatherInfo` object and puts them in the list of `args`. It then returns the sql statement and the arguments.

```golang
func (p postgresSchemaAdapter) InsertQuery(topic string, msgs message.Messages) (string, []interface{}, error) {
	insertQuery := fmt.Sprintf(
		`INSERT INTO %s (last_updated, temperature, feels_like, humidity, condition, wind_mph, wind_direction, visibility, precipitation, created_at) VALUES %s`,
		topic,
		strings.TrimRight(strings.Repeat(`($1,$2,$3,$4,$5,$6,$7,$8,$9,$10),`, len(msgs)), ","),
	)

	var args []interface{}
	for _, msg := range msgs {
		weatherInfo := WeatherInfo{}

		decoder := gob.NewDecoder(bytes.NewBuffer(msg.Payload))
		err := decoder.Decode(&weatherInfo)
		if err != nil {
			return "", nil, err
		}

		args = append(
			args,
			weatherInfo.LastUpdated,
			weatherInfo.Temperature,
			weatherInfo.FeelsLike,
			weatherInfo.Humidity,
			weatherInfo.Condition,
			weatherInfo.WindMph,
			weatherInfo.WindDirection,
			weatherInfo.Visibility,
			weatherInfo.Precipitation,
			weatherInfo.CreatedAt,
		)
	}

	return insertQuery, args, nil
}
```

Here's the part where you'd probably set up a subscriber stream for those data scientists who are teaching ChatGPT to be more conversant about the weather (or something like that). You'll want to create custom `SelectQuery` and `UnmarshalMessage` functions, but since we are not using a SQL subscriber for this tutorial, we'll skip that for now.

```golang
func (p postgresSchemaAdapter) SelectQuery(topic string, consumerGroup string, offsetsAdapter sql.OffsetsAdapter) (string, []interface{}) {
	// No need to implement this method, as PostgreSQL subscriber is not used in this example.
	return "", nil
}

func (p postgresSchemaAdapter) UnmarshalMessage(row *stdSQL.Row) (offset int, msg *message.Message, err error) {
	return 0, nil, errors.New("not implemented")
}
```

### Start the Second Stream

Now we need to put all those last pieces together so that we're storing data back to PostgreSQL.

Here we will use the router functionality described [here](https://rotational.io/blog/prototyping-eda-with-watermill/). We could have had the Ensign subscriber do the entire work of checking the database and inserting new records and not create a publisher at all. However, by decoupling the checking and the inserting functions, we can enable Ensign to scale up and down independently, which could save some serious $$$ depending on how much throughput you're dealing with.

Here, we instantiate a new router, add a `SignalsHandler` plugin that will shut down the router if it receives a `SIGTERM` message. We also add a `Recoverer` middleware which handles any panics sent by the handler.

Next, we will create the PostgreSQL connection and create the `weather_info` table.  We then create the ensign subscriber and the sql publisher.

We will then add a handler to the router called `weather_info_inserter`, pass in the subscriber topic, subscriber, publisher, publisher topic, and the `handlerFunc` that will be executed when a new message appears on the subscriber topic.

Finally, we will run the router.

```golang
func main() {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	//SignalsHandler will gracefully shutdown Router when SIGTERM is received
	router.AddPlugin(plugin.SignalsHandler)
	//The Recoverer middleware handles panics from handlers
	router.AddMiddleware(middleware.Recoverer)

	postgresDB := createPostgresConnection()
	log.Println("added postgres connection and created weather_info table")
	subscriber := createSubscriber(logger)
	publisher := createPublisher(postgresDB)

	router.AddHandler(
		"weather_info_inserter",
		weather_api_topic, //subscriber topic
		subscriber,
		weather_insert_topic, //publisher topic
		publisher,
		dbHandler{db: postgresDB}.checkRecordExists,
	)

	if err = router.Run(context.Background()); err != nil {
		panic(err)
	}
}
```
## Composing our Docker Container

Almost done! In this section we'll create a docker-compose file in the `weather_data` directory to run our application.

The docker-compose file contains three services. The first service is the `producer` which requires the `WAPIKEY` environment variable used to call the Weather API. It will also need the `ENSIGN_CLIENT_ID`, and `ENSIGN_CLIENT_SECRET` to use Ensign. The second service is the `consumer` which needs the `POSTGRES_USER`, `POSTGRES_DB`, and `POSTGRES_PASSWORD` environment variables in order to connect to the database and it will also need `ENSIGN_CLIENT_ID`, and `ENSIGN_CLIENT_SECRET` to use Ensign.  The third service is the `postgres` database, which is a docker image that will also require the same environment variables as the consumer.  You will notice that the container name is `weather_db`, which is the host name that the consumer application uses to connect to the database and it has also got the same Postgres environment variables as the consumer.

```yaml
version: '3'
services:
  producer:
    image: golang:1.19
    restart: unless-stopped
    volumes:
    - .:/app
    - $GOPATH/pkg/mod:/go/pkg/mod
    working_dir: /app/producer/
    command: go run main.go
    environment:
      WAPIKEY: ${WAPIKEY}
	  ENSIGN_CLIENT_ID: ${ENSIGN_CLIENT_ID}
      ENSIGN_CLIENT_SECRET: ${ENSIGN_CLIENT_SECRET}

  consumer:
    image: golang:1.19
    restart: unless-stopped
    depends_on:
    - postgres
    volumes:
    - .:/app
    - $GOPATH/pkg/mod:/go/pkg/mod
    working_dir: /app/consumer/
    command: go run main.go db.go
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
	  ENSIGN_CLIENT_ID: ${ENSIGN_CLIENT_ID}
      ENSIGN_CLIENT_SECRET: ${ENSIGN_CLIENT_SECRET}

  postgres:
    image: postgres:12
    restart: unless-stopped
    ports:
      - 5432:5432
    container_name: weather_db
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
```

## Let's Gooooooooo

We made it to the end! Once you have all of the code in place, run the following commands on the terminal in the `producer` and `consumer` directories:

```bash
go mod init
go mod tidy
```
This will create the `go.mod` and the `go.sum` files in both directories.  Next, move up to the `weather_data` directory and run the following command:

```bash
docker-compose up
```

You will see all the applications running and messages printing to the screen.

![weather_app](/img/weather_app.png)

On a separate terminal window, run the following command to view the contents of the `weather_info` table:

```bash
docker-compose exec weather_db psql -U $POSTGRES_USER -d $POSTGRES_DB -c 'select * from weather_info;'
```

![database_record](/img/database_record.png)

## Next Steps

Hopefully running this example gives you a general idea on how to build an event-driven application using Watermill and Ensign. You can modify this example slightly and have the Ensign consumer do the entire work of checking and inserting new weather records into the database (replace the handler with a `NoPublisherHandler`), but remember that loose coupling is the name of the game with event driven architectures!  You can also challenge yourself by creating a consumer that takes the records produced by the publisher and updates a front end application with the latest weather data.

Imagine all the possibilities that event-driven architectures open up!  You no longer have to worry about multiple users hitting your database and bringing it down by running a massive `select * from` query or worse, dropping a table!  You can simply publish different data streams from your database to meet all of your various end users data requirements.  Now you have even better control over data access as you can select which fields to publish and you can also mask PII data prior to publishing.  No longer do end users need to build their applications using mock data only to find all sorts of data issues in production.  They can directly work with production data and deploy their applications faster and with less headache.

Let us know ([info@rotational.io](mailto:info@rotational.io)) what you end up making with Ensign!