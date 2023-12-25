# Task

## Declaration

### All tasks

**GET: `/api/v2/tasks`**

Response:

```json
[
    {
      id: "1234",
      company: {
        ...
      },
      branch: {
        ...
      },
      user: {
        ...
      },
      name: "Title",
      subject: "Subject",
    },
    ...
]
```

### Create task (General)

**POST `/api/v2/tasks`**

*Must specify complete schema in order to successfully create task*

Request body:

```json
{
  id: "1234",
  company: {
    ...
  },
  branch: {
    ...
  },
  user: {
    ...
  },
  name: "Title",
  subject: "Subject",
}
```


### Create task for user

**POST: `/api/v2/users/:user_id/tasks`**

Request body:

*Other fields are inferred from user*

```json
{
  name: "Title",
  subject: "Subject",
}
```

```go
type Task struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Subject          string           `json:"subject"`
	Status           string           `json:"status"`
	TimeCreated      time.Time        `json:"created_at"`
	LastActivity     time.Time        `json:"activity_at"`
	Company          Company          `json:"company"`
	Branch           Branch           `json:"branch"`
	User             User             `json:"user"`
	Comments         []Comment        `json:"comments"`
	BeforeCreateHook BeforeCreateHook `json:"-" bson:"-"`
}
```
