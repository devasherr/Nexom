# Simple Go ORM

This is a lightweight, chainable ORM (Object-Relational Mapping) implementation in Go, built to provide a fluent and expressive interface for SQL operations without relying on heavy third-party ORM libraries.

## Features

- Fluent interface for building queries
- Type-safe method chaining
- Supports SELECT, INSERT, UPDATE, DELETE operations
- Automatic parameter escaping

## Installation

```bash
go get github.com/devasherr/nexom
```

# Usage

### Initialization

```go
package main

import (
    "github.com/devasherr/nexom"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Initialize the driver
    norm := nexom.New("mysql", "username:password@/dbname")
    defer norm.db.Close()

    // Create an ORM instance for a table
    users := norm.NewOrm("users")
}
```

### SELECT Queries

```go
// Simple select
result, err := users.Select("id", "name", "email").Exec()

// Select with WHERE
result, err := users.Select().Where("id = ?", "1").Exec()

// Select with multiple conditions
result, err := users.Select("name", "email").Where("status = ? AND age > ? OR created_at > ?", "active", "25", "2025-01-01").Exec()
```

### INSERT Queries

```go
// Insert with columns and values
result, err := users.Insert("name", "email", "age").Values("John Doe", "john@example.com", "30").Exec()
```

### UPDATE Queries

```go
// Update with SET and WHERE
result, err := users.Update().
    Set(nexom.M{
        "name": "Jane Doe",
        "email": "jane@example.com",
    }).
    Where("id = ?", "1").
    Exec()

// Update with multiple conditions
result, err := users.Update().
    Set(nexom.M{
        "status": "inactive",
    }).
    Where("last_login = ? AND active = ? OR banned = ?", "< 2023-01-01", "false", "true").Exec()
```

### DELETE Queries

```go
// Simple delete
result, err := users.Delete().Where("id = ?", "1").Exec()

// Delete with multiple conditions
result, err := users.Delete().Where("status = ? AND last_login < ?", "inactive", "2022-01-01").Exec()
```

### DROP TABLE

```go
// Drop table
result, err := users.Drop().Exec()
```
# Context
All operations support context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// SELECT with context
result, err := users.Select().Where("id = ?", "1").ExecContext(ctx)

// INSERT with context
result, err := users.Insert("name", "email").Values("John", "john@example.com").ExecContext(ctx)

// UPDATE with context
result, err := users.Update().Set(nexom.M{"name": "John"}).Where("id = ?", "1").ExecContext(ctx)

// DELETE with context
result, err := users.Delete().Where("id = ?", "1").ExecContext(ctx)

// DROP with context
result, err := users.Drop().ExecContext(ctx)
```

## CONTRIBUTIONS

Pull requests and issues are welcome! If you'd like to contribute, feel free to fork and improve.
