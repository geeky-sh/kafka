My attempt at implementing Kafka functionalities using Golang

### Functional Requirements:
- Create Topic
- Delete Topic
- Add Consumer
- Start Consuming

## Non-functional requirements:
- Ability to publish messages in parallel
- Ability to consume messages depending on the offset

## Notes
- Whenever a publisher adds a message to the topic, notify all its consumers that a message has arrived.
- Use mutex/channel to publish messages

### Terms
- Topic
- Publisher
- Consumer

## TODO later
- Max retention period
- fix casing
- test for error scenarios


## How to test
`go test`
