
# encrypted-udp

This is an experiment to both learn and demonstrate an implementation of a rudimentary encrypted UDP chat client and server.

It would be possible to add RSA encrypted handshakes, but the PKI would have to be something like "compiled with" as opposed to "verified over TLS".

If we added RSA encryption we could also support AES-GCM, but sizes would change, making for a lot of additional logic.

I also did not account for the address size and timestamps for message senders, which may also be useful to add.

The client and server implementation(s) are resilient, meaning if the server goes down the client will automatically "reconnect" (establish new handshake credentials) at the cost of a lost message or two.

This uses no third party packages besides `golang.org/x/crypto` for NaCl.

The clients array may not be concurrently safe, so new users connecting could create a race condition when iterating the list to send a chat message.  It might be more appropriate to use channels for something like this.

While the contents of the messages are encrypted, the signature and type are not, which leaves room for some trouble, but to keep things remotely sane I did not account for a version that encrypts the whole message.

I think a web interface and API for the client would make this more demonstrable, but I don't think I'll put the time or effort into that.

I would also like to create an experiment in the future that uses DHT, the bit torrent protocol.  I think it would be pretty interesting to support server selection and automatic discovery.  That's something that could certainly be a handy demonstration for future projects.

In that same token, a demonstration of very basic NAT punch through for peer-to-peer connections might be interesting.


# references

- [key exchange examples](https://github.com/golang/crypto/blob/master/ssh/kex.go)
- [GoLang Encrypt string to base64 and vice versa using AES encryption.](https://gist.github.com/manishtpatel/8222606)
- [Latency computation](https://gamedev.stackexchange.com/questions/105196/how-to-measure-packet-latency)
- [running average](https://github.com/RobinUS2/golang-moving-average)
- [exponential average](https://stackoverflow.com/questions/12538959/algorithm-for-counting-things-by-the-second-and-maintaining-a-running-average)
- [bluntly (nodejs project)](https://github.com/danoctavian/bluntly)
- [Golang : UDP client server read write example](https://socketloop.com/tutorials/golang-udp-client-server-read-write-example)
- [CREATING A SECURE SERVER IN GOLANG](https://austburn.me/blog/golang-server.html)
- [Cryptographic Best Practices](https://gist.github.com/atoponce/07d8d4c833873be2f68c34f9afc5a78a)
- [Go Lang NACL Cryptography](https://8gwifi.org/docs/go-nacl.jsp)
- [AES-NI SSL Performance](https://calomel.org/aesni_ssl_performance.html)
- [What is the difference between UDP hole punching and UPnP?](https://superuser.com/questions/617263/what-is-the-difference-between-udp-hole-punching-and-upnp)
- [UDP AGAINST ROUTERS : HOLE PUNCHING](http://jwhsmith.net/2014/03/udp-routers-hole-punching/)
- [Game servers: UDP vs TCP=](https://1024monkeys.wordpress.com/2014/04/01/game-servers-udp-vs-tcp/)
