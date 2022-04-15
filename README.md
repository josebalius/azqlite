# azqlite [![Go Reference](https://pkg.go.dev/badge/github.com/josebalius/azqlite.svg)](https://pkg.go.dev/github.com/josebalius/azqlite)

azqlite is a lightweight wrapper around [github.com/Azure/azure-storage-queue-go](https://github.com/Azure/azure-storage-queue-go) to interact with the Azure Storage Queue service in a simpler and more idiomatic way.

## Install

```
go get github.com/josebalius/azqlite
```

## How to use

### Instantiate a service 
```go
client, err := azqlite.NewClient(azqlite.Config{
	AccountName: "YOUR_AZURE_STORAGE_ACCOUNT_NAME_HERE",
	AccountKey:  "YOUR_AZURE_STORAGE_ACCOUNT_KEY_HERE",
})
```

### Create a queue
```go
q, err := client.CreateQueue(ctx, "test")
```

### Delete a queue
```go
err = c.DeleteQueue(ctx, "test")
```

### Get an existing queue
```go
q := c.GetQueue("test")
```

### Get message count
```go
c, err := q.MessageCount(ctx)
```

### Enqueue a message
```go
m, err := q.Enqueue(ctx, "my message", 1*time.Second, -time.Second)
```

### Dequeue messages
```go
messages, err := q.Dequeue(ctx, 30, 1*time.Second)
```

### Peek messages
```go
messages, err := q.Peek(ctx, 30)
```

### Delete a message
```go
err := q.Delete(ctx, &Message{ID: "1"})
```
