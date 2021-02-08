import argparse
import pika
import waggle.message as message
import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from prometheus_client import start_http_server, Counter
import os
import logging


RABBITMQ_URL = os.environ.get("RABBITMQ_URL", "amqp://localhost")
RABBITMQ_EXCHANGE = os.environ.get("RABBITMQ_EXCHANGE", "waggle.msg")
RABBITMQ_QUEUE = os.environ.get("RABBITMQ_QUEUE", "influxdb-messages")
INFLUXDB_URL = os.environ["INFLUXDB_URL"]
INFLUXDB_TOKEN = os.environ["INFLUXDB_TOKEN"]
INFLUXDB_BUCKET = os.environ.get("INFLUXDB_BUCKET", "waggle")
INFLUXDB_ORG = os.environ.get("INFLUXDB_ORG", "waggle")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--debug", action="store_true")
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

    logging.info("connecting to influxdb at %s", INFLUXDB_URL)
    client = influxdb_client.InfluxDBClient(
        url=INFLUXDB_URL,
        token=INFLUXDB_TOKEN,
        org=INFLUXDB_ORG)
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

        writer.write(bucket=INFLUXDB_BUCKET, org=INFLUXDB_ORG, record=record)
        ch.basic_ack(method.delivery_tag)
        logging.debug("proccessed message %s", msg)
        waggle_messages_total.labels(node=msg.meta["node"], plugin=msg.meta["plugin"]).inc()

    params = pika.URLParameters(RABBITMQ_URL)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(RABBITMQ_QUEUE, durable=True)
    ch.queue_bind(RABBITMQ_QUEUE, RABBITMQ_EXCHANGE, "#")
    ch.basic_consume(RABBITMQ_QUEUE, message_handler)
    ch.start_consuming()


if __name__ == "__main__":
    main()
