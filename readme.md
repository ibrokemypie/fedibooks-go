# fedibooks-go

Work in progress ebooks bot. Uses markov chains to generate random text from followed users' post histories.

## Features currently missing

- Getting message history from remote instances.

- Max thread depth

## Usage

Create a user for the bot and follow the users you want to source from with it. Run ``fedibooks-go`` for the first time to authenticate as the bot user through oauth.  

Optionally ``fedibooks-go`` can be used with the ``-c`` flag to use a specified config file, the default is ``./config.yaml``.

## Configuration

Sample Config:

```yaml
history:
  file_path: ./history.gob
  get_interval: 30
  learn_from_cw: false
  max_length: 100000
instance:
  access_token: xxxxxxxx
  url: https://mastodon.social
post:
  make_interval: 30
  max_words: 30
  visibility: unlisted
```

|Setting|Description|
|---|---|
|history.get_interval|Interval in minutes to try to retrieve new posts from followed users|
|history.learn_from_cw|Whether to learn from sensitive/CW'd posts|
|history.file_path|Path and name of history file on disk|
|history.max_length|Maximum number of statuses to store and generate from|
|post.visibility|Visibility of generated posts|
|post.make_interval|Interval in minutes to generate a new post|
|post.max_words|Maximum numqber of words to generate per post|
|instance.access_token|OAuth access token. Generated on first run|
|instance.url|Bot's instance URL. Generated on first run|

