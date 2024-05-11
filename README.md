# Run

Script run.sh can be used to start the application. Application will use following ports:

- 26257
- 8080
- 5672
- 15672
- 3000

# API
GET and POST methods are allowed for /messages enpoint. Using the GET method, the client retrieves all existing messages and new ones are streamed. With POST method new message will be created.

Example GET:
```cURL
curl http://localhost:3000/messages
```

Example POST:
```cURL
curl -d '{"content": "New Message"}' http://localhost:3000/messages
```
