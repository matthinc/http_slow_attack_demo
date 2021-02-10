#!/usr/bin/env python3

import socket


TCP_IP = "target"
TCP_PORT = 7000
BUFFER_SIZE = 8
HEADER = "POST /echo HTTP/1.1\r\nHost: Attacker\r\nContent-Type: text/plain\r\nContent-Length: 32\r\n\r\n"
MESSAGE = "Hello, World! From the attacker."

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((TCP_IP, TCP_PORT))
s.send((HEADER + MESSAGE).encode("utf-8"))
i = 1000
while i > 0:
    data = s.recv(BUFFER_SIZE)
    i -= 1
s.close()

print("received data:" + str(data))
