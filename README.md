# Cypto.com assignment

## Contents
- [Quick start](#quick-start)
- Decisions
    - use RDBMS to make trade data consistency
    - use repository pattern to encapsulate DB access
- how long you spent on the test
    - around 6-8 hours

## Quick start
- [How to use](#how-to-use)
- [Others](#others)

There are two methods that can launch the application.
1. via the Makefile to setup the application (**recommend**)
```txt
make run
```
2. through the docker-compose to buildup the application manually
```txt
docker-compose up --build
```

## How to use
- The server is listening to 8080 port, you can change it in the docker-compose.yaml
## authorization
- add header **{"Authorization", {token}}** in your HTTP.Header
- for simplicity, we have some hardcode tokens to handle it
```txt
['Tim', 'Alex', 'Arthur', 'Ray', 'HD', 'peko', 'miko', 'rushia', 'gura', 'Ame']
```

### withdraw
```txt
POST: localhost:8080/api/v1/wallet/withdraw

Header: {
    "Authorization": {{token}},
    "Content-Type": "application/json"
}

RequestBody: {
	"amount": integer (required)
}

Response:
	200: OK
	400: BadRequest
	401: Unauthorized
	500: serverError 

```

### Deposit
```txt
POST: localhost:8080/api/v1/wallet/deposit

Header: {
    "Authorization": {{token}},
    "Content-Type": "application/json"
}

RequestBody: {
	"amount": integer (required)
}

Response:
	200: OK
	400: BadRequest
	401: Unauthorized
	500: serverError 

```

### Transfer
```txt
POST: localhost:8080/api/v1/wallet/transfer

Header: {
    "Authorization": {{token}},
    "Content-Type": "application/json"
}

RequestBody: {
	"toAccount": string (required)
	"amount": integer (required)
}

Response:
	200: OK
	400: BadRequest
	401: Unauthorized
	500: serverError 
```

### GetAccount
```txt
GET: localhost:8080/api/v1/wallet/account

Header: {
    "Authorization": {{token}}
}

Response:
	200: OK
	401: Unauthorized
	500: serverError 
```


## Others
1. build images
```
make build
```
2. clean up the environment
```
make clean
```