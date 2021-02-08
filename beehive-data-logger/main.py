import argparse
import pika
import waggle.message as message
import os
import json
import sys


RABBITMQ_URL = os.environ.get("RABBITMQ_URL", "amqp://localhost")
RABBITMQ_EXCHANGE = os.environ.get("RABBITMQ_EXCHANGE", "waggle.msg")
RABBITMQ_QUEUE = os.environ.get("RABBITMQ_QUEUE", "logger-messages")


def message_handler(ch, method, properties, body):
    try:
        msg = message.load(body)
    except Exception:
        ch.basic_ack(method.delivery_tag)
        print("failed to parse message", file=sys.stderr, flush=True)
        return

    log = json.dumps({
        "timestamp": msg.timestamp,
        "name": msg.name,
        "meta": msg.meta,
        "value": msg.value,
    }, separators=(",", ":"))

    print(log, flush=True)
    ch.basic_ack(method.delivery_tag)


def main():
    params = pika.URLParameters(RABBITMQ_URL)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(RABBITMQ_QUEUE, durable=True)
    ch.queue_bind(RABBITMQ_QUEUE, RABBITMQ_EXCHANGE, "#")
    ch.basic_consume(RABBITMQ_QUEUE, message_handler)
    ch.start_consuming()


if __name__ == "__main__":
    main()
