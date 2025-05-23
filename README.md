# Software Engineering School 5.0 // Case Task - Weather API Application

## Table of Contents

- [Description](#description)
- [Deployment](#deployment)
- [Tests coverage](#tests-coverage)
- [Running locally](#running-locally)
  - [Running in Docker (mocked APIs)](#running-in-docker)
  - [Configuring the full application](#configuring-the-full-application)
  - [Migrating the database](#migrating-the-database)
- [Known limitations, issues and possible improvements](#known-limitations-issues-and-possible-improvements)
- [Contacts](#contacts)

## Description

Weather API application allows users to subscribe to weather updates for a chosen city. 

![Weather API application](./assets/pictures/page.png)

An example of the weather notification email:

![Weather notification email](./assets/pictures/notification.png)

Application is built using Go 1.24 programming language.

Integrations:
- `Weather API` (https://www.weatherapi.com/) — for weather data;
- `Mailjet` (https://www.mailjet.com/) — for sending emails.

Frameworks and libraries (most significant):
- `chi` — HTTP shenanigans;
- `pgdb`, `squirrel` — for PostgreSQL;
- `spf13/cobra` — for CLI;
- `html/template` — for HTML templates (email bodies) and static HTML page;
- `distributed_lab` — kit for building microservices (configuration, logging, middlewares, etc.).
- `mockery` — for mocking interfaces (testing purposes);

The implementation consists of an HTTP server that handles the required endpoints to manage user subscriptions
and a background worker that periodically fetches pending subscriptions to notify and sends emails with weather updates.

## Deployment

The application is deployed on the private VM and can be accessed via the following links:
- http://34.59.198.232:8080 — for the web-page;
- http://34.59.198.232:8090 — for the API.

**Note: There is no static IP address assigned to the VM, so the IP address may change.
See [contacts](#contacts) section to report a possible app failure**

## Tests coverage

Tests cover only API server (integration-like with mocked interfaces); they are located in the one [file](./internal/api/server_test.go).
The common test cases map practice is used in a pair with the `httptest` package to test the HTTP server.
The mocked interfaces are used to test the application without the need for real API calls or database access.

As a possible improvement, the tests can be split into several files and test the application on a handler func level.

## Running locally

### Running in Docker
There is a [docker-compose.yml](./build/docker-compose.yml) file supposed to spin up the application alongside PostgreSQL database.
It is located in the `/build` directory.

The configured application includes:
- `PostgreSQL` database for storing the data;
- application itself with the next open ports:
  - `8090` for HTTP API;
  - `8080` for `index.html` page;

To run the application, execute the following command in the root directory of the project:

```bash
docker-compose -f build/docker-compose.yml up
```

**NOTE! The Docker Compose setup is configured to use the mock APIs for Weather API and Mailjet (no one wants the keys to be exposed, right?)**

### Configuring the full application
To set up the application without mock APIs, you simply need to 
1. Add your API keys to the [/build/config.yaml](./build/config.yaml) file:
```yaml
weather_api:
  api_key: YOUR_API_KEY

mailjet:
  api_key: YOUR_API_KEY
  secret_key: YOUR_API_KEY
  from_email: YOUR_EMAIL
```
2. Remove the `--mocks=true` flag from the application entrypoint in the [docker-compose.yml](./build/docker-compose.yml)

See [Contacts](#contacts) section to get the API keys.

### Migrating the database

To migrate the database, you can use the `migrate up/down` commands provided by the CLI.

When running with Docker Compose, the `docker-compose.yml` setup already includes the migration step, so you don't need to worry about it.


## Known limitations, issues and possible improvements
- confirmation/unsubscription tokens have no expiration time (although this is not defined by the specification provided);
- there is no confirmation/unsubscription link in the email body (although this is not defined by the specification provided);
- the application does not support multiple subscriptions for the same email address (although this is not defined by the specification provided);
- the spec defines `Subscription` model, but it never uses it, so do I;
- the spec doesn't define `500 Internal Server Error` response, but I've included it in the code;
- there is no usage of batch processing for sending emails, batch querying the data from the weather API, bulk updates of db records.
Instead, there is a simple concurrent processing of the data using weather data caching, semaphore and transactional queries to fully control the processing flow;
- database schema is simplified to the one table with all the data in it (however, everything still looks okay);
- **`index.html` page is AI-generated (I can't stand frontend development)**

## Contacts

If there are any questions related to the project (application cannot be started properly, links do not work, API keys required),
please contact [Max](https://t.me/slbmax)