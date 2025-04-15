# develop

## llm
`https://chat.qwen.ai/`

## generate input text for LLM
```shell 
repomix --ignore "**/*test.go,**/suite.go"
# or, to include tests
repomix

cat repomix-output.md |tail +38|pbcopy
```

# run
## run backend
```shell
go run main.go
# Server running on http://localhost:8080
```

## test python
```shell 
go run main.go & sleep 2;python debug_api.py ;kill $(lsof -t -i :8080)
```


## test (fails for now)
```shell
  go test -coverpkg=./... -coverprofile=coverage.out ./...
```

## update requirements
```shell
go mod tidy
```

## send api requests
```shell
python3 debug_api.py
```

## build for any platform
```shell
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build.exe main.go
```

## run dev 
```shell
go run main.go
```


## see db contents
```shell
    sqlite3 main.db
```

# Examples
## 1. Register a user:
```shell
curl -X POST -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "password": "securepassword123"
}' http://localhost:8080/register

```
## 2. Login to get JWT token:
```shell
curl -X POST -H "Content-Type: application/json" -d '{
  "identity": "john_doe",
  "password": "securepassword123"
}' http://localhost:8080/login
```

## 3. Create a house (replace <TOKEN> with your JWT):

```shell
curl -X POST -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{
    "name": "Test Estate",
    "type": "flat_yearly",
    "price": 2500,
    "address": "123 Main St",
    "in_rent": True,
    "user_id": 0,
}' http://localhost:8080/api/estates
```
## 4. Update a house:

```shell
curl -X PUT -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{
  "name": "Winter Cottage",
  "price": 420000
}' http://localhost:8080/api/estates/0
```

## 5. Delete a house:
```shell
curl -X DELETE -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/estates/0
```
