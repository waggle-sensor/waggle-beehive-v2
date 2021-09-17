import argparse
from os import getenv
import pika
import pika.credentials
import json
import sys
import ssl
import wagglemsg as message


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
    parser = argparse.ArgumentParser()
    parser.add_argument("--rabbitmq_host",default=getenv("RABBITMQ_HOST", "localhost"))
    parser.add_argument("--rabbitmq_port", default=getenv("RABBITMQ_PORT", "5672"), type=int)
    parser.add_argument("--rabbitmq_username", default=getenv("RABBITMQ_USERNAME", ""))
    parser.add_argument("--rabbitmq_password", default=getenv("RABBITMQ_PASSWORD", ""))
    parser.add_argument("--rabbitmq_cacertfile", default=getenv("RABBITMQ_CACERTFILE", ""))
    parser.add_argument("--rabbitmq_certfile", default=getenv("RABBITMQ_CERTFILE", ""))
    parser.add_argument("--rabbitmq_keyfile", default=getenv("RABBITMQ_KEYFILE", ""))
    parser.add_argument("--rabbitmq_exchange", default=getenv("RABBITMQ_EXCHANGE", "waggle.msg"))
    parser.add_argument("--rabbitmq_queue", default=getenv("RABBITMQ_QUEUE", "logger-messages"))
    args = parser.parse_args()

    if args.rabbitmq_username != "":
        credentials = pika.PlainCredentials(args.rabbitmq_username, args.rabbitmq_password)
    else:
        credentials = pika.credentials.ExternalCredentials()

    if args.rabbitmq_cacertfile != "":
        context = ssl.create_default_context(cafile=args.rabbitmq_cacertfile)
        # HACK this allows the host and baked in host to be configured independently
        context.check_hostname = False
        if args.rabbitmq_certfile != "":
            context.load_cert_chain(args.rabbitmq_certfile, args.rabbitmq_keyfile)
        ssl_options = pika.SSLOptions(context, args.rabbitmq_host)
    else:
        ssl_options = None

    params = pika.ConnectionParameters(
        host=args.rabbitmq_host,
        port=args.rabbitmq_port,
        credentials=credentials,
        ssl_options=ssl_options,
        retry_delay=60,
        socket_timeout=10.0)

    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(args.rabbitmq_queue, durable=True)
    ch.queue_bind(args.rabbitmq_queue, args.rabbitmq_exchange, "#")
    ch.basic_consume(args.rabbitmq_queue, message_handler)
    ch.start_consuming()


if __name__ == "__main__":
    main()
