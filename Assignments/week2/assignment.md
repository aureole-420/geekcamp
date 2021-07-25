## Question
> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码

## Answer
> 是。如果没有error wrapping, DAO layer的upstream caller会缺乏error handling的context. 我们可以直接用golang1.13的errors.Wrap()方法来返回wrapped error.

Yes, `sql.ErrNoRows` error should be wrapped before returning to upstream caller of DAO layer. Otherwise, this error is not handled in DAO layer so it is supposed to be handled by the upstream caller. But the upstream caller basically obtains only a raw error, i.e., it lacks the context to decide how to handle the error.

We can use the `Wrap()` method provided in standard `errors` library or `github.com/pkg/errors` package.

Assuming DAO has exposed a function named `GetUser` which retrieve the user object based on the username.


```golang
import (
    "context"
    "database/sql"
    "github.com/pkg/errors"
)

type DBService struct { // DAO layer
    db *sql.DB
}

//....function to initialize DBService

func(dbs *DBService) GetUser(cxt context.Context, username string) (User, error) {
    db := dbs.db

    rows, err := db.QueryContext(cxt, "SELETE * FROM users WHERE name=?", username)
    if err != nil {
        return nil, errors.Wrap(err, fmt.sprintf("Failing in getting user object, username=%s", username))
    }

    ... // construct result; 
    return result, nil
}
```