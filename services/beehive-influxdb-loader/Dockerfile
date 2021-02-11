FROM python:3-alpine

RUN pip3 install --no-cache-dir https://github.com/waggle-sensor/pywaggle/archive/v0.40.4.zip influxdb_client

COPY main.py main.py

ENTRYPOINT [ "python", "main.py" ]
