#! /usr/bin/env python

from requests import post
from sys import argv as orig_argv, stdin
from uuid import uuid4


def argv(i: int, default: str = "") -> str:
    return orig_argv[i] if len(orig_argv) > i else default


def read() -> str:
    return stdin.read()


if __name__ == "__main__":
    host = argv(1, "localhost")
    port = argv(2, "12345")
    key = argv(3, "loggy")
    environment = argv(4, "dev")
    app_version = argv(5, "1.0.0")
    device_name = argv(6, "send.py")
    message = read()

    res = post(
        f"http://{host}:{port}/logs",
        json={
            "key": key,
            "environment": environment,
            "app_version": app_version,
            "device_name": device_name,
            "message": message,
        },
    )
    print(res.json())
