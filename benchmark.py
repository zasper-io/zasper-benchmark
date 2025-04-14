
# Paths to your JSON files
zasper_file = "benchmark_results_zasper.json"
jupyter_file = "benchmark_results_new.json"

import json
import matplotlib.pyplot as plt  # type: ignore
from datetime import datetime

# Load benchmark results for Zasper and JupyterLab from JSON files
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
    return normalized_time, cpu_usage, memory_usage

# Extract data
zasper_time, zasper_cpu_usage, zasper_memory_usage = extract_metrics(zasper_data)
jupyter_time, jupyter_cpu_usage, jupyter_memory_usage = extract_metrics(jupyter_data)

# Create two subplots: one for CPU usage and one for Memory usage
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 10))

# Plot CPU usage
ax1.plot(zasper_time, zasper_cpu_usage, label="Zasper CPU Usage", color='blue', marker='o')
ax1.plot(jupyter_time, jupyter_cpu_usage, label="JupyterLab CPU Usage", color='green', marker='s')
ax1.set_xlabel('Time (seconds since start)')
ax1.set_ylabel('CPU Usage (%)')
ax1.set_title('CPU Usage Comparison: Zasper vs JupyterLab')
ax1.legend()
ax1.grid(True)

# Plot Memory usage
ax2.plot(zasper_time, zasper_memory_usage, label="Zasper Memory Usage (MB)", color='blue', marker='o')
ax2.plot(jupyter_time, jupyter_memory_usage, label="JupyterLab Memory Usage (MB)", color='green', marker='s')
ax2.set_xlabel('Time (seconds since start)')
ax2.set_ylabel('Memory Usage (MB)')
ax2.set_title('Memory Usage Comparison: Zasper vs JupyterLab')
ax2.legend()
ax2.grid(True)

# Show the plot
plt.tight_layout()
plt.show()
