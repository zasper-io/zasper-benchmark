# Benchmarking Zasper v/s Jupyter Server

![](/assets/results.png)

# Objective
The primary goal of this benchmarking exercise is to compare the performance of Zasper against the traditional Jupyter Server. The focus areas for evaluation are:

* CPU Usage
* RAM Usage
* Throughput Performance

Through this comparison, we aim to determine how Zasper performs in a real-world scenario where multiple execute requests are made, with particular interest in resource consumption and efficiency.


# Understanding Jupyter Server Architecture

To establish a baseline, it is important to understand how a Jupyter Server operates internally. Here's a simplified breakdown:

### 1. Session Lifecycle

A new session is initiated when a user opens a Jupyter notebook.

This session launches a kernel, which handles code execution.

### 2. Communication via WebSocket
A WebSocket is established between the Jupyter client and the server, allowing real-time, bi-directional communication.

### 3. Kernel Channels
The Jupyter kernel communicates with the server over five dedicated channels:
* stdin – for user inputs.
* shell – for sending execution requests.
* control – for kernel control messages.
* iopub – for publishing results back to the client.
* heartbeat – for kernel liveliness checks.

📌 For this benchmarking exercise, we focus only on:

* Shell channel – used to send execution requests (e.g., 2+2, print("Hello World!"))
* IOPub channel – used to receive outputs from the kernel (e.g., 4, Hello World!)


## Methodology

The benchmarking setup follows a controlled and repeatable process:

### 1. Session Initialization
A session is created and a WebSocket connection is established using a goroutine.

### 2. Execution Requests
A stream of kernel_execute_request messages is sent over the shell channel. Each message triggers a computation in the kernel.

### 3. Monitoring & Logging
System metrics such as CPU usage, memory consumption, and execution throughput are recorded at 10-second intervals. These are visualized for comparison.


# Results

The graph shows a clear performance difference between Zasper and Jupyter Server across the selected metrics.

### Key observations:

* CPU Usage: Zasper maintained consistently lower CPU usage during heavy load.

* RAM Usage: Memory consumption was significantly lower for Zasper.

* Throughput: Zasper handled more execution requests per unit of time, indicating better scalability under concurrent workloads.

## Benefits of Zasper

### For Individual Users
* Improved Responsiveness: Faster execution of notebook cells.
* Lightweight: Reduced memory usage allows smoother multitasking, especially on lower-spec machines.

### For Enterprises
* Cost Efficiency: Lower resource usage translates to fewer cloud compute instances required.

* Better Scalability: Efficient resource handling allows support for more users and sessions per node.


## Conclusion

This benchmarking study highlights Zasper's performance advantages over the traditional Jupyter Server. Whether for individual developers or large-scale enterprise deployments, Zasper demonstrates meaningful improvements in resource efficiency and execution throughput, making it a promising alternative for interactive computing environments.

# Copyright

Prasun Anand
