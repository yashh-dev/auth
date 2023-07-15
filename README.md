# auth service
This is the microservice that handles everything auth related.


## Events
These are the events this microservice is listening to.

- `auth.password.initial`: event which triggers user and verification_token creation.
- `auth.password.change`: event which triggers user password change.
- `auth.login`: event which triggers session creation.
- `auth.session.exists`: event which checks if current session is valid.
- `auth.session.get-user`: event which returns the current user of active session.
- `auth.verify`: event which verifies the user

## Deployment (without miauw stack)

1. Image Build
```sh
$ docker build -t miauw/auth .
```
2. Run Image
```
$ docker run -d miauw/auth
```
