import argparse
import pika
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS, WritePrecision
from os import getenv
import logging
import ssl
import waggle.message as message


def assert_type(obj, t):
    if not isinstance(obj, t):
        raise TypeError(f"{obj!r} must be {t}")


def assert_maxlen(s, n):
    if len(s) > n:
        raise ValueError(f"len({s!r}) must be <= {n}")


def assert_valid_message(msg):
    assert_type(msg.name, str)
    assert_maxlen(msg.name, 64)
    assert_type(msg.timestamp, int)
    assert_type(msg.value, (int, float, str))
    assert_type(msg.meta, dict)
    for k, v in msg.meta.items():
        assert_type(k, str)
        assert_maxlen(k, 64)
        assert_type(v, str)
        assert_maxlen(v, 64)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--debug", action="store_true")
    parser.add_argument("--rabbitmq_host",default=getenv("RABBITMQ_HOST", "localhost"))
    parser.add_argument("--rabbitmq_port", default=getenv("RABBITMQ_PORT", "5672"), type=int)
    parser.add_argument("--rabbitmq_username", default=getenv("RABBITMQ_USERNAME", ""))
    parser.add_argument("--rabbitmq_password", default=getenv("RABBITMQ_PASSWORD", ""))
    parser.add_argument("--rabbitmq_cacertfile", default=getenv("RABBITMQ_CACERTFILE", ""))
    parser.add_argument("--rabbitmq_certfile", default=getenv("RABBITMQ_CERTFILE", ""))
    parser.add_argument("--rabbitmq_keyfile", default=getenv("RABBITMQ_KEYFILE", ""))
    parser.add_argument("--rabbitmq_exchange", default=getenv("RABBITMQ_EXCHANGE", "waggle.msg"))
    parser.add_argument("--rabbitmq_queue", default=getenv("RABBITMQ_QUEUE", "influx-messages"))
    parser.add_argument("--influxdb_url", default=getenv("INFLUXDB_URL", "http://localhost:8086"))
    parser.add_argument("--influxdb_token", default=getenv("INFLUXDB_TOKEN"))
    parser.add_argument("--influxdb_bucket", default=getenv("INFLUXDB_BUCKET", "waggle"))
    parser.add_argument("--influxdb_org", default=getenv("INFLUXDB_ORG", "waggle"))
    
    args = parser.parse_args()

    logging.basicConfig(
        level=logging.DEBUG if args.debug else logging.INFO,
        format="%(asctime)s %(message)s",
        datefmt="%Y/%m/%d %H:%M:%S")
    # pika logging is too verbose, so we turn it down.
    logging.getLogger("pika").setLevel(logging.CRITICAL)

    logging.info("connecting to influxdb at %s", args.influxdb_url)
    client = influxdb_client.InfluxDBClient(
        url=args.influxdb_url,
        token=args.influxdb_token,
        org=args.influxdb_org)
    logging.info("connected to influxdb")

    writer = client.write_api(write_options=SYNCHRONOUS)

    def message_handler(ch, method, properties, body):
        try:
            msg = message.load(body)
        except Exception:
            ch.basic_ack(method.delivery_tag)
            logging.warning("failed to parse message")
            return
        
        try:
            assert_valid_message(msg)
        except Exception:
            logging.exception("invalid message - dropping message")
            ch.basic_ack(method.delivery_tag)
            return

        record = {
            "measurement": msg.name,
            "tags": msg.meta,
            "fields": {
                "value": msg.value,
            },
            "time": msg.timestamp,
        }

        writer.write(bucket=args.influxdb_bucket, org=args.influxdb_org, record=record, write_precision=WritePrecision.NS)
        ch.basic_ack(method.delivery_tag)
        logging.debug("proccessed message %s", msg)

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
