
# udp

This is a rudimentary UDP client-server demonstration, which passes "chat" messages to a server and distributes them to all clients.

It is the foundation for alternative implementations:

- UDP /w TLS encryption
- Peer server /w NAT Punch Through
- web based client interface /w additional metrics such as latency

All three are common, valid patterns used for network communication where TCP is for whatever reason not an option (_games /w packet loss latency, systems that need to accept dropped traffic, or custom prioritization_).
