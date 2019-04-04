# Pure-Go implementation of sfsdb

```go

import (
    "github.com/AlmightyFloppyFish/sfsdb-go"
    "fmt"
)

// Must be exported
type User struct {
    Name string
    Id int
}

func main() {
    db := sfsdb.New("db", 10) // Cache up to 10 objects
    user := User { "some_name", 13142 }

    db.Save("some_key", user)

    var retrieved User
    db.Load("some_key", &retrieved)

    fmt.Println(user, "is same as", retrieved)
}
```
