# PAGERDUTY AUTO ACKNOWLEDGER

because pagerduty notifications can get annoying to have to manually acknowledge each time. this lets you view your notifications async without your boss yelling at you :)

Usage:
build the docker image, and run like so:
```
docker run --env PD_AUTH_TOKEN=<token> --env PD_EMAIL=<email> --env PD_USER_ID=<user-id> --detach <container-id>
```

You NEED to have the pagerduty authorization info. You can get this by following Pagerduty API docs, it's not too hard.

Ideally, you would deploy this on an EC2 instance to run forever during your oncall shift.
