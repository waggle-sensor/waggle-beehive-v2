import argparse
import pika
import waggle.message as message
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from prometheus_client import start_http_server, Counter
from os import getenv
import logging


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--debug", action="store_true")
    parser.add_argument("--rabbitmq_host",default=getenv("RABBITMQ_HOST", "localhost"))
    parser.add_argument("--rabbitmq_port", default=getenv("RABBITMQ_PORT", "5672"), type=int)
    parser.add_argument("--rabbitmq_username", default=getenv("RABBITMQ_USERNAME", "guest"))
    parser.add_argument("--rabbitmq_password", default=getenv("RABBITMQ_PASSWORD", "guest"))
    parser.add_argument("--rabbitmq_ssl_cacertfile", default=getenv("RABBITMQ_SSL_CACERTFILE"))
    parser.add_argument("--rabbitmq_exchange", default=getenv("RABBITMQ_EXCHANGE", "waggle.msg"))
    parser.add_argument("--rabbitmq_queue", default=getenv("RABBITMQ_QUEUE", "influxdb-messages"))
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

    # setup and run metrics for prometheus monitoring
    waggle_messages_total = Counter("waggle_messages_total", "Total messages for a node and plugin", ["node", "plugin"])
    waggle_errors_total = Counter("waggle_errors_total", "Total number of errors")
    start_http_server(9123)

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
            waggle_errors_total.inc()
            return

        try:
            record = (influxdb_client.Point(msg.name)
                .tag("node", msg.meta["node"])
                .tag("plugin", msg.meta["plugin"])
                .field("value", msg.value))
        except KeyError as key:
            ch.basic_ack(method.delivery_tag)
            logging.warning("message missing meta %s", key)
            waggle_errors_total.inc()
            return

        writer.write(bucket=args.influxdb_bucket, org=args.influxdb_org, record=record)
        ch.basic_ack(method.delivery_tag)
        logging.debug("proccessed message %s", msg)
        waggle_messages_total.labels(node=msg.meta["node"], plugin=msg.meta["plugin"]).inc()

    params = pika.URLParameters(args.rabbitmq_url)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(args.rabbitmq_queue, durable=True)
    ch.queue_bind(args.rabbitmq_queue, args.rabbitmq_exchange, "#")
    ch.basic_consume(args.rabbitmq_queue, message_handler)
    ch.start_consuming()


if __name__ == "__main__":
    main()
