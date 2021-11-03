
# user-service

Is a background process for sshd server to create user accounts via an API.

## API Resources

```text
GET  /user          # list all existsing users with prefix ("node-")
POST /user/{user}   # create user
```

Usernames must match `^[0-9a-f]{16}`, prefix `node-` will appended to username if missing. 