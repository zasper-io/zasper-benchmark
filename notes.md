# Notes

```
/Users/prasunanand/Library/Python/3.9/bin/jupyter lab
```

```
[I 2025-04-12 04:56:20.124 ServerApp] Shutting down 5 extensions
[I 2025-04-12 04:56:20.125 ServerApp] Shutting down 2 kernels
[I 2025-04-12 04:56:20.126 ServerApp] Discarding 81000 buffered messages for c9b1c25a-98e1-4ac3-b863-0cf5688daf3d:6e84f277-656d-4d62-9820-e8b9b88acc05
[I 2025-04-12 04:56:20.136 ServerApp] Kernel shutdown: c9b1c25a-98e1-4ac3-b863-0cf5688daf3d
[I 2025-04-12 04:56:20.136 ServerApp] Discarding 4805 buffered messages for e8ec36f7-b65e-4598-9d9f-40330ebdcf1e:388991e0-8e9f-4c0c-b9fa-728f9155a09e
[I 2025-04-12 04:56:20.138 ServerApp] Kernel shutdown: e8ec36f7-b65e-4598-9d9f-40330ebdcf1e
```

At (64 kernels, 10RPS per kernel) and (16 kernels, 100RPS per kernel)  Jupyter Server completely gave up
```
[W 2025-04-26 22:48:39.098 ServerApp] Write error on <socket.socket fd=637, family=AddressFamily.AF_INET6, type=SocketKind.SOCK_STREAM, proto=0, laddr=('::1', 8888, 0, 0), raddr=('::1', 57168, 0, 0)>: [Errno 55] No buffer space available
[W 2025-04-26 22:48:39.099 ServerApp] Write error on <socket.socket fd=161, family=AddressFamily.AF_INET6, type=SocketKind.SOCK_STREAM, proto=0, laddr=('::1', 8888, 0, 0), raddr=('::1', 56615, 0, 0)>: [Errno 55] No buffer space available
[W 2025-04-26 22:48:39.099 ServerApp] Write error on <socket.socket fd=198, family=AddressFamily.AF_INET6, type=SocketKind.SOCK_STREAM, proto=0, laddr=('::1', 8888, 0, 0), raddr=('::1', 56658, 0, 0)>: [Errno 55] No buffer space available
Task exception was never retrieved
future: <Task finished name='Task-82495' coro=<WebSocketProtocol13.write_message.<locals>.wrapper() done, defined at /Users/prasunanand/Library/Python/3.9/lib/python/site-packages/tornado/websocket.py:1086> exception=WebSocketClosedError()>
Traceback (most recent call last):
  File "/Users/prasunanand/Library/Python/3.9/lib/python/site-packages/tornado/websocket.py", line 1088, in wrapper
    await fut
tornado.iostream.StreamClosedError: Stream is closed

During handling of the above exception, another exception occurred:

```
```
[I 2025-04-26 22:48:39.134 ServerApp] Starting buffering for 3677e004-a553-479c-8cb9-f0da390eee27:1371dd36-816c-4fa0-a63b-fc7429bfd43b
Task exception was never retrieved
future: <Task finished name='Task-82551' coro=<WebSocketProtocol13.write_message.<locals>.wrapper() done, defined at /Users/prasunanand/Library/Python/3.9/lib/python/site-packages/tornado/websocket.py:1086> exception=WebSocketClosedError()>
Traceback (most recent call last):
  File "/Users/prasunanand/Library/Python/3.9/lib/python/site-packages/tornado/websocket.py", line 1088, in wrapper
    await fut
tornado.iostream.StreamClosedError: Stream is closed

During handling of the above exception, another exception occurred:

Traceback (most recent call last):
  File "/Users/prasunanand/Library/Python/3.9/lib/python/site-packages/tornado/websocket.py", line 1090, in wrapper
    raise WebSocketClosedError()
tornado.websocket.WebSocketClosedError
```

```
{"level":"info","time":1745735833,"message":"Error writing message: write tcp [::1]:8048->[::1]:51161: write: no buffer space available"}
{"level":"info","time":1745735834,"message":"Error writing message: write tcp [::1]:8048->[::1]:50991: write: no buffer space available"}
{"level":"error","error":"writev tcp 127.0.0.1:51485->127.0.0.1:5679: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"writev tcp 127.0.0.1:51136->127.0.0.1:5647: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"writev tcp 127.0.0.1:51024->127.0.0.1:5230: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"zmq4: read/write on closed connection","time":1745735834,"message":"failed to send message"}
```

At (16 kernels, 100RPS per kernel)  Zasper completely gave up because of Zeromq,

```
{"level":"info","time":1745735833,"message":"Error writing message: write tcp [::1]:8048->[::1]:51161: write: no buffer space available"}
{"level":"info","time":1745735834,"message":"Error writing message: write tcp [::1]:8048->[::1]:50991: write: no buffer space available"}
{"level":"error","error":"writev tcp 127.0.0.1:51485->127.0.0.1:5679: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"writev tcp 127.0.0.1:51136->127.0.0.1:5647: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"writev tcp 127.0.0.1:51024->127.0.0.1:5230: writev: no buffer space available","time":1745735834,"message":"failed to send message"}
{"level":"error","error":"zmq4: read/write on closed connection","time":1745735834,"message":"failed to send message"}
```