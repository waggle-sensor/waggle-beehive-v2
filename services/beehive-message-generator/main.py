import argparse
from os import getenv
import pika
import json
import waggle.message as message
import time
import random
import ssl


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

    while True:
        msg = message.Message(
            name="env.temperature.gen",
            timestamp=time.time_ns(),
            value=random.uniform(0.0, 5.0),
            meta={"node": "0000000000000001", "plugin": "metsense:1.0.2"})
        body = message.dump(msg)
        properties = pika.BasicProperties(user_id="node-0000000000000001")
        ch.basic_publish("waggle.msg", routing_key="", properties=properties, body=body)

        msg = message.Message(
            name="sys.uptime",
            timestamp=time.time_ns(),
            value=time.time(),
            meta={"node": "0000000000000001", "plugin": "status:1.0.0"})
        body = message.dump(msg)
        properties = pika.BasicProperties(user_id="node-0000000000000001")
        ch.basic_publish("waggle.msg", routing_key="", properties=properties, body=body)

        msg = message.Message(
            name="sys.uptime",
            timestamp=time.time_ns(),
            value=time.time()+1.4,
            meta={"node": "0000000000000002", "plugin": "status:1.0.0"})
        body = message.dump(msg)
        properties = pika.BasicProperties(user_id="node-0000000000000002")
        ch.basic_publish("waggle.msg", routing_key="", properties=properties, body=body)

        msg = message.Message(
            name="sys.uptime",
            timestamp=time.time_ns(),
            value=time.time()+2.3,
            meta={"node": "0000000000000003", "plugin": "status:1.0.0"})
        body = message.dump(msg)
        properties = pika.BasicProperties(user_id="node-0000000000000003")
        ch.basic_publish("waggle.msg", routing_key="", properties=properties, body=body)

        time.sleep(1)


if __name__ == "__main__":
    main()
