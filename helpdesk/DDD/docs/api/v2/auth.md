# Authorization

## UserFlow

```
user -> /identity
                |
                `-> OK (identified)
                |
                `-> FAIL (user is null) -> /register
                |                                  |
                |                                  `-> OK (registered) -> OK (identified)
                |                                  |
                |                                  `-> FAIL(user exists) -> ?
                |
                `-> FAIL (device's ip is not correspond to any network)
                                                                      |
                                                                      `-> FAIL (data not found)
```
