import json
import os
import time

import pika

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
EXCHANGE = "southpark"
QUEUE = "southpark_messages"
BINDING_KEY = "message.#"


def connect():
    """Connect to RabbitMQ, retrying while it boots."""
    params = pika.URLParameters(RABBITMQ_URL)
    for attempt in range(1, 31):
        try:
            return pika.BlockingConnection(params)
        except pika.exceptions.AMQPConnectionError:
            print(f"RabbitMQ not ready (attempt {attempt}/30), retrying in 2s...", flush=True)
            time.sleep(2)
    raise SystemExit("Could not connect to RabbitMQ, giving up.")


def setup(channel):
    """Declare the same exchange, queue and binding as the Go API."""
    channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
    channel.queue_declare(queue=QUEUE, durable=True)
    channel.queue_bind(queue=QUEUE, exchange=EXCHANGE, routing_key=BINDING_KEY)
    channel.basic_qos(prefetch_count=1)


def on_message(channel, method, properties, body):
    """Called by pika for every message in the queue."""
    try:
        msg = json.loads(body)
        author = msg["author"]
        text = msg["body"]
    except (json.JSONDecodeError, KeyError) as err:
        print(f"Skipping bad message: {err}", flush=True)
        channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
        return

    print(f'[{method.routing_key}] {author}: "{text}"', flush=True)
    channel.basic_ack(delivery_tag=method.delivery_tag)


def main():
    connection = connect()
    channel = connection.channel()
    setup(channel)

    channel.basic_consume(queue=QUEUE, on_message_callback=on_message)
    print(f"Listening on '{QUEUE}'. Press Ctrl+C to stop.", flush=True)

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()