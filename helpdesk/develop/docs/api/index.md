# API Documentation

## Response

```typescript
interface Response {
    status: bool,
    ?data: any,
    ?error: {
        type: string,
        message: string,
    },
}

```

## Routes


### Authorization 

POST: `/api/auth/token`

payload: request's ip address

response:

```typescript
interface Token
{
    token: string,
    refresh_token: string,
}
```

### Registration

POST: `/api/auth/register`

payload: request's ip address +
```typescript
interface RegisterForm {
    name: string,
    phone: string
}
```
response:
```typescript
interface Token 
{
    token: string,
    refresh_token: string,
}
```

### Client's routes

#### Get Profile Information

GET: `/api/helpdesk/profile`

payload: none

headers: token

response:

```typescript
interface Profile
{
    name: string,
    phone: string,
    company: {
        name: string,
    },
    branch: {
        name: string,
        description: string,
        address: string,
        contacts: string,
    }
}
```

#### Get User's Telegram Pass

GET: `/api/helpdesk/telegram`

payload: none

headers: token

response: `telegram token`

#### Get User's Tasks

GET: `/api/helpdesk/tasks`

payload: none

headers: token

response: `Task[]`

```typescript
interface Task {
    id: string,
    name: string,
    subject: string,
    status: string,
    created_at: Date,
    activity_at: Date,
}
```

#### Get Branch's Tasks

GET: `/api/helpdesk/tasks?filter=branch`

payload: none

headers: token

response: `Task[]`

```typescript
interface Task {
    id: string,
    name: string,
    subject: string,
    status: string,
    created_at: Date,
    activity_at: Date,
    user: {
        name: string,
        phone: string,
    }
}
```

#### Get Company's Tasks

GET: `/api/helpdesk/tasks?filter=company`

payload: none

headers: token

response: `Task[]`

```typescript
interface Task {
    id: string,
    name: string,
    subject: string,
    status: string,
    created_at: Date,
    activity_at: Date,
    branch: {
        name: string,
        description: string,
        address: string,
        contacts: string,
    },
    user: {
        name: string,
        phone: string,
    }
}
```


#### Get Task

GET: `/api/helpdesk/tasks/:id`

payload: none

headers: token

response: `Task`

```typescript
interface Task {
    id: string,
    name: string,
    subject: string,
    status: string,
    created_at: Date,
    activity_at: Date,
    company: {
        name: string,
    }
    branch: {
        name: string,
        description: string,
        address: string,
        contacts: string,
    },
    user: {
        name: string,
        phone: string,
    },
    comments: {
        id: string,
        content: string,
        ?user: {
            name: string,
            phone: string,
        },
        direction: "to" | "from",
        created_at: Date,
    }[],
}
```

#### Get Task's comments

GET: `/api/helpdesk/tasks/:id/comments`

payload: none

headers: token

response: `Comment[]`

```typescript
interface Comment {
    id: string,
    content: string,
    ?user: {
        name: string,
        phone: string,
    } | undefined,
    direction: "to" | "from",
    created_at: Date
}
```

#### Create New Task

POST: `/api/helpdesk/tasks`

payload:
```typescript
interface NewTask {
    name: string,
    subject: string,
}
```

headers: token

response: `task id`

#### Comment Task

POST: `/api/helpdesk/tasks/:id/comments`

payload:
```typescript
interface NewComment {
    content: string
}
```

headers: token

response: `comment id`
