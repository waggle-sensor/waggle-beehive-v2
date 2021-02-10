# RabbitMQ Config Notes

## message-generator in definitions.json

This user is meant to publish generated messages into the data pipeline to help test Beehive.
To do this, the user needs the `impersonator` tag so it can make up RabbitMQ user IDs with messages
being dropped. For more info, see: https://www.rabbitmq.com/validated-user-id.html
