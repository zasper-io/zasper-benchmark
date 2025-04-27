import json
import matplotlib.pyplot as plt  # type: ignore
from datetime import datetime
import argparse
import os
os.makedirs("plots", exist_ok=True)


parser = argparse.ArgumentParser(description="Generate visualizations for Zasper and Jupyter.")
parser.add_argument("--n", type=int, required=True, help="Number of kernels")
parser.add_argument("--delay", type=str, required=True, help="delay between requests")

args = parser.parse_args()
n = args.n
delay = args.delay

# File names depend on n
zasper_file = f"data/{delay}ms/benchmark_results_zasper_{n}kernels.json"
jupyter_file = f"data/{delay}ms/benchmark_results_jupyter_{n}kernels.json"


# Load benchmark results for Zasper and Jupyter Server from JSON files
with open(zasper_file, 'r') as f:
    zasper_data = json.load(f)

with open(jupyter_file, 'r') as f:
    jupyter_data = json.load(f)

# Function to extract normalized time and usage metrics
def extract_metrics(data):
    timestamps = [datetime.fromisoformat(entry['timestamp']) for entry in data]
    start_time = timestamps[0]
    normalized_time = [(ts - start_time).total_seconds() for ts in timestamps]
    cpu_usage = [entry['cpu_usage'] for entry in data]
    memory_usage = [entry['memory_usage_mb'] for entry in data]
    messages_sent_count = [entry['messages_sent_count'] for entry in data]
    messages_received_count = [entry['messages_received_count'] for entry in data]
    message_sent_throughput = [entry['message_sent_throughput'] for entry in data]
    message_received_throughput = [entry['message_received_throughput'] for entry in data]
    return normalized_time, cpu_usage, memory_usage, messages_sent_count, messages_received_count, message_sent_throughput, message_received_throughput

# Extract data
zasper_time, zasper_cpu_usage, zasper_memory_usage, zasper_messages_sent_count, zasper_messages_received_count, zasper_message_sent_throughput, zasper_message_received_throughput = extract_metrics(zasper_data)
jupyter_time, jupyter_cpu_usage, jupyter_memory_usage, jupyter_messages_sent_count, jupyter_messages_received_count, jupyter_message_sent_throughput, jupyter_message_received_throughput = extract_metrics(jupyter_data)

# Create two subplots: one for CPU usage and one for Memory usage
fig, ((ax1, ax2), (ax3, ax4), (ax5, ax6)) = plt.subplots(3, 2, figsize=(20, 10))

# Plot CPU usage
ax1.plot(zasper_time, zasper_cpu_usage, label="Zasper CPU Usage", color='#583BD8', marker='o')
ax1.plot(jupyter_time, jupyter_cpu_usage, label="Jupyter Server CPU Usage", color='#E46E2E', marker='s')
ax1.set_xlabel('Time (seconds since start)')
ax1.set_ylabel('CPU Usage (%)')
ax1.set_title(f'CPU Usage Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax1.legend()
ax1.grid(True)

# Plot Memory usage
ax2.plot(zasper_time, zasper_memory_usage, label="Zasper Memory Usage (MB)", color='#583BD8', marker='o')
ax2.plot(jupyter_time, jupyter_memory_usage, label="Jupyter Server Memory Usage (MB)", color='#E46E2E', marker='s')
ax2.set_xlabel('Time (seconds since start)')
ax2.set_ylabel('Memory Usage (MB)')
ax2.set_title(f'Memory Usage Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax2.legend()
ax2.grid(True)


# Plot Message Sent
ax3.plot(zasper_time, zasper_messages_sent_count, label="Zasper Message Sent Count", color='#583BD8', marker='o')
ax3.plot(jupyter_time, jupyter_messages_sent_count, label="Jupyter Server Message Sent Count", color='#E46E2E', marker='s')
ax3.set_xlabel('Time (seconds since start)')
ax3.set_ylabel('Messages Sent')
ax3.set_title(f'Messages Sent Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax3.legend()
ax3.grid(True)

# Plot Message Receieved
ax4.plot(zasper_time, zasper_messages_received_count, label="Zasper Message Receieved Count", color='#583BD8', marker='o')
ax4.plot(jupyter_time, jupyter_messages_received_count, label="Jupyter Server Message Receieved Count", color='#E46E2E', marker='s')
ax4.set_xlabel('Time (seconds since start)')
ax4.set_ylabel('Messages Received')
ax4.set_title(f'Messages Received Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax4.legend()
ax4.grid(True)


# Plot Message Sent Throughput
ax5.plot(zasper_time, zasper_message_sent_throughput, label="Zasper Message Sent Throughput", color='#583BD8', marker='o')
ax5.plot(jupyter_time, jupyter_message_sent_throughput, label="Jupyter Server Message Sent Throughput", color='#E46E2E', marker='s')
ax5.set_xlabel('Time (seconds since start)')
ax5.set_ylabel('Messages Sent/second')
ax5.set_title(f'Message Sent Throughput Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax5.legend()
ax5.grid(True)

# Plot Message Receieved Throughput
ax6.plot(zasper_time, zasper_message_received_throughput, label="Zasper Message Receieved Throughput", color='#583BD8', marker='o')
ax6.plot(jupyter_time, jupyter_message_received_throughput, label="Jupyter Server Message Receieved Throughput", color='#E46E2E', marker='s')
ax6.set_xlabel('Time (seconds since start)')
ax6.set_ylabel('Messages Sent/second')
ax6.set_title(f'Message Receieved Throughput Comparison: Zasper vs Jupyter Server | {n} kernels | 10 RPS per kernel')
ax6.legend()
ax6.grid(True)


# Show the plot
plt.tight_layout()
plt.savefig(f"plots/{delay}ms/benchmark_result_{n}kernels.png")
# plt.show()
