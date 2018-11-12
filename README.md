# lunch-lambda
A lambda service to notify me what's in lunch at office daily.

Too lazy to open up the portal. :sleeping:

## How to use?
So basically it's deployed on AWS Lambda (but can easily be used as CLI, with
minor changes). You can easily build this for Lamda using the following
command:

```bash
$ GOOS=linux go build .
```

To set it up, you will have to provision your Lambda function with the
following environment variables:

Name | Description | Note
-|-|-
`COGS_USERNAME` | Your username for 10Pearls COGS (the portal)
`COGS_PASSWORD` | Your password for the same ^
`CLIENT_NUM` | Number of clients where you need push notification
`CLIENT_{x}` | Name of the push notification service (currently supported: `pushbullet` only)
`CLIENT_{x}_TOKEN` | Access Token for push notification service | Valid if `CLIENT_{x}` is `pushbullet` only

Note: `{x}` must start from `1` till `CLIENT_NUM`.

PS: I know array starts from zero. :sweat_smile:

## Why did I make this?
I was getting bored. :neutral_face:
