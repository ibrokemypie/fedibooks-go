# fedibooks-go

Work in progress ebooks bot. Uses markov chains to generate random text from followed users' post histories.

## Features currently missing

- Getting message history from remote instances.

- Reply handling.

- Max post length.

## Usage

Create a user for the bot and follow the users you want to source from with it. Run ``fedibooks-go`` for the first time to authenticate as the bot user through oauth.  

Optionally ``fedibooks-go`` can be used with the ``-c`` flag to use a specified config file, the default is ``./config.yaml``.

## Configuration

Sample Config:

```yaml
get_posts_interval: 30
history:
  file_path: ./history.gob
  max_length: 100000
instance:
  access_token: xxxxxxxx
  url: https://mastodon.social
learn_from_cw: false
make_post_interval: 30
post_visibility: unlisted
```

|Setting|Description|
|---|---|
|get_posts_interval|interval in minutes to try to retrieve new posts from followed users|
|make_post_interval|interval in minutes to generate a new post|
|learn_from_cw|whether to learn from sensitive/CW'd posts|
|post_visibility|visibility of generated posts|
|instance.access_token|OAuth access token. Generated on first run|
|instance.url|Bot's instance URL. Generated on first run|
|history.file_path|Path and name of history file on disk|
|history.max_length|Maximum number of statuses to store and generate from|
