# acmeii
irc client for plan9 acme & ii

To use, start ii:
```
ii
```

Then start acmeii in acme:
```
acmeii
```

A new window will open for each server, channel, or nick you've
joined or are chatting with.

To build:
```
go install github.com/raylai/acmeii
go install github.com/raylai/acmeii/acmeiiwin
```

## Todo

- `pledge(2)` for OpenBSD
- Exec `ii` (with `tcpclient` as needed), take same flags as `ii`
(so calling `ii` and `acmeii` has same semantics, except `acmeii`
opens new windows)
- `ii` itself doesn't have any output, so combine with `/irc/$server`
window; that's one less window to manage
- Use 9fans.net/go/acme package to manipulate windows rather than
calling `win`; this will enable more fine-grained controls of
windows, such as clearing history, `rm`ing channels we don't care
about anymore, etc.
- Reconnection?
- Write man page

Enjoy! That's an order.
