# azqlite

azqlite is lightweight wrapper around github.com/Azure/azure-storage-queue-go/azqueue to interact with the Azure Storage Queue service in a simpler and more idiomatic way.

## Install

```
go get github.com/josebalius/azqlite
```

## How to use

### Instantiate a service 
```
storageService, err := azqlite.NewService(azqlite.Config{
	AccountName: "YOUR_AZURE_STORAGE_ACCOUNT_NAME_HERE",
	AccountKey:  "YOUR_AZURE_STORAGE_ACCOUNT_KEY_HERE",
})
```

### Create a queue
```
q, err := storageService.CreateQueue(ctx, "test")
```

### Delete a queue
```
err = s.DeleteQueue(ctx, "test")
```

### Instantiate an existing queue
```
q := s.NewQueue("test")
```

### Get message account
```
c, err := q.MessageCount(ctx)
```

### Enqueue a message
```
m, err := q.Enqueue(ctx, "my message", 1*time.Second, -time.Second)
```

### Dequeue messages
```
messages, err := q.Dequeue(ctx, 30, 1*time.Second)
```

### Peek messages
```
messages, err := q.Peek(ctx, 30)
```

### Delete a message
```
err := q.Delete(ctx, &Message{ID: "1"})
```
