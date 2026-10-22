import os
import random
import time

import requests

API_URL = os.getenv("API_URL", "http://localhost:8080/messages")
INTERVAL = float(os.getenv("INTERVAL_SECONDS", "3"))
KENNY_FAIL_CHANCE = 0.15

MESSAGES = [
    ("Cartman", "Respect my authoritah!"),
    ("Cartman", "Screw you guys, I'm going home!"),
    ("Cartman", "I'm not fat, I'm big-boned!"),
    ("Stan", "Oh my God, they killed Kenny!"),
    ("Kyle", "You bastards!"),
    ("Butters", "Oh hamburgers!"),
    ("Kenny", "Mmmph mmmph mmph!"),
    ("Mr Mackey", "Unhandled errors are bad, mkay?"),
    ("Randy", "I'm not drunk, I'm load testing the message queue!"),
    ("Towelie", "Don't forget to bring a towel... and a dead-letter queue."),
    ("Cartman", "This queue will do what I say, because I am the producer!"),
]


def send_random():
    """Send one random message; sometimes Kenny sends a broken one."""
    if random.random() < KENNY_FAIL_CHANCE:
        resp = requests.post(API_URL, json={"author": "Kenny"}, timeout=5)
        print(f"Kenny tried to talk without a body -> {resp.status_code} {resp.text.strip()}"
              f"  (Oh my God, they killed Kenny's message!)", flush=True)
        return

    author, body = random.choice(MESSAGES)
    resp = requests.post(API_URL, json={"author": author, "body": body}, timeout=5)
    print(f"sent {author}: {body} -> {resp.status_code}", flush=True)


def main():
    print(f"Spamming {API_URL} every {INTERVAL}s. Ctrl+C to stop.", flush=True)
    while True:
        try:
            send_random()
        except requests.RequestException as err:
            print(f"API not reachable ({err.__class__.__name__}), retrying...", flush=True)
        time.sleep(INTERVAL)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("Spammer stopped. Screw you guys, I'm going home!")