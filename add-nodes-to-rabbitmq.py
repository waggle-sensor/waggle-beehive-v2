#!/usr/bin/env python3
import argparse
import json
import subprocess

parser = argparse.ArgumentParser()
parser.add_argument("ids", nargs="*", help="list of node IDs to add")
args = parser.parse_args()

users = [f"node-{id}" for id in args.ids]

definitions = {
    "users": [
        {
            "name": user,
            "tags": "",
            "limits": {}
        } for user in users
    ],
    "permissions": [
        {
            "user": user,
            "vhost": "/",
            "configure": "^$",
            "write": "waggle.msg",
            "read": "^$"
        } for user in users
    ],
}

payload = json.dumps(definitions).encode()

# TODO change to use rabbitmq api instead of kubectl exec
subprocess.run([
    "kubectl", "exec", "-i", "svc/rabbitmq", "--",
    "rabbitmqctl", "--timeout", "300", "import_definitions",
], input=payload)
