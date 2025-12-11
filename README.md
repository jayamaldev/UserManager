# User Manager

User-Manager is a simple Golang application with CRUD Operations.

## Starting The Application

1. Add a .env file with DB connection details. Use the sample .env template file.


```
DB_HOST=database host name (use database service name from the docker-compose file)
DB_PORT=5432
DB_USER=DB user name ex:<<postgres>>
DB_PASSWORD=DB User Password ex:<<123>>
DB_NAME=DB Name ex:<<user_manager>>

APP_PORT=8080
```

2. Start the Dockerized Application with database. Use the docker compose file for this.

```
docker compose up
``` 


## Usage
### Rest End Points

#### Get All Users
```
GET <<http://localhost:8080>>/users
```
#### Get a Single User
GET <<http://localhost:8080>>/users/<ID>

#### Add a User
POST <<http://localhost:8080>>/users

**Request JSON Body**
```json
{
    "firstName": "Jay",
    "lastName": "sV",
    "email": "mail@maail.com",
    "phone": "+876543219",
    "age": 35,
    "status": "Active"
}
```

#### Delete User
DELETE <<http://localhost:8080>>/users/<ID>

#### Update User
PATHC <<http://localhost:8080>>/users/<ID>

**Request JSON Body**
```json
{
    "firstName": "Jay",
    "lastName": "Vas",
    "email": "mail@maail.com",
    "phone": "0876543219",
    "age": 32,
    "status": "Active"
}
```

## Swagger URL

```
http://localhost:8080/doc/index.html
```

## Linting Report
```
./Report.xml
```

## Test Coverage
To Run the Integration Tests, User-Manager uses TestContainers to run a test database.

### Run Integration Tests
```
go test
```

### Run Unit Tests
```
 go test .\internal\
