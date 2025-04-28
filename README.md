# Benchmarking Zasper v/s Jupyter Server

![](/assets/results.png)
(JupyterLab is running 2 Ipython kernels whereas Zasper is running 20 Ipython kernels. Hence Zasper consumes more CPU than Jupyter server at `t=0`)

## System Specifications

* OS : macOS
* CPU : Apple M4, 10-core CPU
* RAM: 16GB

Note: A typical IPython kernel consumes around 80 MB of RAM on avaerage.

![](/assets/idle_ram.png)

(RAM usage on M4 Mac mini)

On my M4 Mac mini, I can see that leftover RAM is around 9 GB , hence the number of kernels that can fit on my machine is 9GB/80MB = 112 ~= 100 Jupyter kernels.

On an M3 Macbook Air which has just 8GB RAM, the leftover RAM tends to be around 1GB RAM , so we can fit ~10 Ipython kernels running on that machine.

Hence, if you want to run the benchmarks make sure that you have enough RAM for the kernels, else you might end up with results that won't make sense.


# Introduction

The primary goal of this benchmarking exercise is to compare the performance of Zasper against the traditional Jupyter Server. The focus areas for evaluation are:

* CPU Usage
* RAM Usage
* Throughput
* Latency
* Resilience

Through this comparison, we aim to determine how Zasper performs in a real-world scenario where multiple execute requests are made, with particular interest in resource consumption and efficiency.


# Understanding Jupyter Server Architecture

To establish a baseline, it is important to understand how a Jupyter Server operates internally. Here's a simplified breakdown:

### 1. Session Lifecycle

A new session is initiated when a user opens a Jupyter notebook.

This session launches a kernel, which handles code execution.


### 2. Kernel Channels
The Jupyter kernel communicates with the server over five dedicated channels:
* stdin – for user inputs.
* shell – for sending execution requests.
* control – for kernel control messages.
* iopub – for publishing results back to the client.
* heartbeat – for kernel liveliness checks.

📌 For this benchmarking exercise, we focus only on:

* Shell channel – used to send execution requests (e.g., `2+2`, `print("Hello World!")`)
* IOPub channel – used to receive outputs from the kernel (e.g., `4`, `Hello World!`)

### 3. Communication via WebSocket
A WebSocket is established between the Jupyter client and the server, allowing real-time, bi-directional communication. The client send the messages over the websocket. When the jupyter_server receives this message it puts this message on a `shell channel` over zeromq. This message when received by the kernel  triggers a computation in the kernel. The kernel emits the output on `iopub channel` over zeromq. This message is received by Jupyter server and the output is put on websocket.

![](/assets/kernel_communication.svg)


## Methodology

The benchmarking setup follows a controlled and repeatable process:

### 1. Session Initialization
A session is created and a WebSocket connection is established using a goroutine.

### 2. Execution Requests
A stream of `execute_request` kernel messages is sent over the websocket.

### 3. Monitoring & Logging
System metrics such as CPU usage, memory consumption, and execution throughput are recorded at 10-second intervals. These are visualized for comparison.

## Steps to run

* Setting up the benchmark code
```
git clone https://github.com/zasper-io/zasper-benchmark
cd zasper-benchmark
# Install go dependencides
go mod tidy
# Install Python dependencies
pip install -r requirements.txt
```


* Collecting data for zasper

1. Start Zasper

2. Start the monitoring code
```
go run .
```
Run with `--debug` flag to see the `requests` and `responses` happening in real time.

```
go run . --debug
```

```
prasunanand@Prasuns-Mac-mini zasper-benchmark % go run .
prasunanand@Prasuns-Mac-mini zasper-benchmark % go run .
====================================================================
*******            Measuring performance                     *******
====================================================================
Target: zasper
PID: 70049
Number of kernels: 2
Output file: data/benchmark_results_zasper_2kernels.json
====================================================================
Creating kernel sessions ⏳
Sessions created:  ✅
Start sending requests: ⏳
Kernel messages sent:  ✅
====================================================================
*******                   Summary                            *******
====================================================================
Messages sent: 38
Messages received: 192
====================================================================
```
The program writes the output to `data/benchmark_results_zasper_2kernels.json` file.

* Collecting data for Jupyterlab

1. Start JupyterLab.
2. You need to get `api_token` and `xsrf_token` and paste it in the `.env` file.
3. Start the monitoring code
```
go run .
```
The program writes the output to `benchmark_results_jupyterlab.json`

* Visualize the data

```
python visualize.py
```


# Results

The graph shows a clear performance difference between Zasper and Jupyter Server across the selected metrics.

### 2 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_2kernels.png)

Lower CPU uasage and RAM usage is better.
Higher Message sent and  Message receievd is better
Higher Message sent per second and Message received per second  is better.

### 4 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_4kernels.png)

### 8 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_8kernels.png)

### 16 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_16kernels.png)

### 32 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_32kernels.png)

### 64 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_64kernels.png)

Note the messages received throughput starts to drop here.

### 100 kernels | 10 RPS per kernel

![](/plots/100ms/benchmark_result_100kernels.png)

Note the messages received throughput drops to 0.

### Resource Usage summary | 10 RPS per kernel
![](/plots/100ms/summary_resources.png)


### 2 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_2kernels.png)

### 4 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_4kernels.png)

### 8 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_8kernels.png)

### 16 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_16kernels.png)

### 32 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_32kernels.png)

### 64 kernels | 100 RPS per kernel

![](/plots/10ms/benchmark_result_64kernels.png)

### Resource Usage summary | 100 RPS per kernel
![](/plots/10ms/summary_resources.png)


### Observations

Zasper consumer lesser CPU and lesser Memory in all cases.

For 64 kernels, Jupyter server starts loosing kernel connections. For 100 kernels, Jupyter server starts to hang.
Message received throughput falls to 0.

**Note that when you are running a Jupyterlab session on the cloud and JupyterLab session completely hangs and becomes unresponsive, this is what is happening.**

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


### Key observations:

* CPU Usage: Zasper maintained consistently lower CPU usage during heavy load.
* RAM Usage: Memory consumption was significantly lower for Zasper.
* Throughput: Zasper handled more execution requests per unit of time, indicating better scalability under concurrent workloads.
* Latency: Under extremely high load, the latency drops for both JupyterLab and Zasper. However Zasper has much lower latency compared to Jupyter Server.
* Resilience: Zasper is a lot more resilient compared to Jupyterlab, and can easily recover.

## Why Zasper Outperforms Jupyter Server

Go is a compiled language with native support for concurrency and multi-core scalability, whereas Python is an interpreted language that primarily runs on a single core. This fundamental difference gives **Zasper**, built in Go, a significant performance advantage over **Jupyter Server**, which is built in Python.

Jupyter Server uses the **Tornado** web server, which is built around Python’s **asyncio** framework for handling asynchronous requests. In contrast, Zasper leverages Go’s **Gorilla** server, which utilizes Go’s lightweight **goroutines** for concurrency. While both are asynchronous in nature, goroutines are much more efficient and cheaper to schedule compared to Python’s event-loop-based coroutines.

In Jupyter Server, submitting a request to the ZeroMQ channels involves packaging an asynchronous function into the asyncio event loop, along with futures and callbacks. The loop must then schedule and manage these functions—an operation that introduces overhead. Zasper, on the other hand, creates goroutines with minimal scheduling cost, making the process significantly faster.

While Python’s asyncio and Go’s goroutines share similar architectural goals, Go's model is much closer to the hardware. It schedules coroutines across multiple CPU threads seamlessly, while Python is limited by the **Global Interpreter Lock (GIL)**, preventing true multi-core parallelism.

When request handling slows down in Jupyter Server, memory usage climbs, CPU gets overwhelmed, and the garbage collector (GC) starts to intervene—often resulting in degraded performance or even crashes.

Zasper is designed around the principle of **“Use More to Save More.”** As request volume increases, Zasper’s efficiency becomes more apparent. Its architecture thrives under load, delivering better throughput and stability at scale.


## Benefits of Zasper

### For Individual Users
* Improved Responsiveness: Faster execution of notebook cells.
* Lightweight: Reduced memory usage allows smoother multitasking, especially on lower-spec machines.

### For Enterprises
* Cost Efficiency: Lower resource usage translates to fewer cloud compute instances required.
* Better Scalability: Efficient resource handling allows support for more users and sessions per node.

## Conclusion

This benchmarking study highlights Zasper's performance advantages over the traditional Jupyter Server. Whether for individual developers or large-scale enterprise deployments, Zasper demonstrates meaningful improvements in resource efficiency and execution throughput, making it a promising alternative for interactive computing environments.


# Thanks to Jupyter Community

Zasper would not exist without the incredible work of the Jupyter community. Zasper uses the Jupyter wire protocol and draws inspiration from its architecture. Deep thanks to all Jupyter contributors for laying the groundwork. Data Science Notebooks would not have existed without them.

# Copyright

Prasun Anand
