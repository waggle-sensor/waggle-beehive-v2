import argparse
import pika
import waggle.message as message
import os
import json
import sys
import ssl


RABBITMQ_HOST = os.environ.get("RABBITMQ_HOST", "localhost")
RABBITMQ_PORT = int(os.environ.get("RABBITMQ_PORT", "5671"))
RABBITMQ_USERNAME = os.environ.get("RABBITMQ_USERNAME", "guest")
RABBITMQ_PASSWORD = os.environ.get("RABBITMQ_PASSWORD", "guest")
RABBITMQ_SSL_CACERTFILE = os.environ.get("RABBITMQ_SSL_CACERTFILE")
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
    credentials = pika.PlainCredentials(RABBITMQ_USERNAME, RABBITMQ_PASSWORD)

    if RABBITMQ_SSL_CACERTFILE is not None:
        context = ssl.create_default_context(cafile=RABBITMQ_SSL_CACERTFILE)
        # HACK this allows the host and baked in host to be configured independently
        context.check_hostname = False
        ssl_options = pika.SSLOptions(context, RABBITMQ_HOST)
    else:
        ssl_options = None

    params = pika.ConnectionParameters(
        host=RABBITMQ_HOST,
        port=RABBITMQ_PORT,
        credentials=credentials,
        ssl_options=ssl_options,
        retry_delay=60,
        socket_timeout=10.0,
    )

    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(RABBITMQ_QUEUE, durable=True)
    ch.queue_bind(RABBITMQ_QUEUE, RABBITMQ_EXCHANGE, "#")
    ch.basic_consume(RABBITMQ_QUEUE, message_handler)
    ch.start_consuming()


if __name__ == "__main__":
    main()
