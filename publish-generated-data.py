import pika
import json
import waggle.message as message
import time
import random

params = pika.URLParameters("amqp://service:service@localhost")
conn = pika.BlockingConnection(params)
ch = conn.channel()

while True:
    msg = message.Message(
        name="env.temperature.gen",
        timestamp=time.time_ns(),
        value=random.uniform(0.0, 5.0),
        meta={"node": "0000000000000001", "plugin": "metsense:1.0.2"})
    body = message.dump(msg)
    ch.basic_publish("waggle.msg", routing_key="", body=body)

    msg = message.Message(
        name="sys.uptime",
        timestamp=time.time_ns(),
        value=time.time(),
        meta={"node": "0000000000000001", "plugin": "status:1.0.0"})
    body = message.dump(msg)
    ch.basic_publish("waggle.msg", routing_key="", body=body)

    msg = message.Message(
        name="sys.uptime",
        timestamp=time.time_ns(),
        value=time.time()+1.4,
        meta={"node": "0000000000000002", "plugin": "status:1.0.0"})
    body = message.dump(msg)
    ch.basic_publish("waggle.msg", routing_key="", body=body)

    msg = message.Message(
        name="sys.uptime",
        timestamp=time.time_ns(),
        value=time.time()+2.3,
        meta={"node": "0000000000000003", "plugin": "status:1.0.0"})
    body = message.dump(msg)
    ch.basic_publish("waggle.msg", routing_key="", body=body)

    time.sleep(1)

# we can also exclude "zero values" and assume they're empty.
