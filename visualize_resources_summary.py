import json
import matplotlib.pyplot as plt
import numpy as np
import argparse

parser = argparse.ArgumentParser(description="Generate visualizations for Zasper and Jupyter.")
parser.add_argument("--delay", type=str, required=True, help="delay between requests")

args = parser.parse_args()
delay = args.delay

# Kernel counts to process
kernel_counts = [2, 4, 8, 16, 32, 64, 100]

# Prepare storage
avg_cpus_zasper = []
avg_mems_zasper = []
max_cpus_zasper = []
max_mems_zasper = []

avg_cpus_jupyter = []
avg_mems_jupyter = []
max_cpus_jupyter = []
max_mems_jupyter = []

# Define helper to add annotation
def add_note(ax, text, position='upper right'):
    positions = {
        'upper right': dict(x=0.95, y=0.95, ha='right', va='top'),
        'upper left': dict(x=0.05, y=0.95, ha='left', va='top'),
        'center': dict(x=0.5, y=0.5, ha='center', va='center'),
        'center left': dict(x=0.05, y=0.5, ha='left', va='center'),
        'center right': dict(x=0.95, y=0.5, ha='right', va='center'),
        'lower right': dict(x=0.95, y=0.05, ha='right', va='bottom'),
        'lower left': dict(x=0.05, y=0.05, ha='left', va='bottom'),
    }
    ax.text(
        s=text,
        transform=ax.transAxes,
        fontsize=9,
        bbox=dict(facecolor='white', edgecolor='black', boxstyle='round,pad=0.3'),
        **positions[position]
    )

# Helper function to compute averages
def compute_averages(data):
    cpu = [entry['cpu_usage'] for entry in data]
    memory = [entry['memory_usage_mb'] for entry in data]

    return (np.mean(cpu), np.mean(memory), np.max(cpu), np.max(memory))

# Loop through the files
for k in kernel_counts:
    filename = f'data/{delay}ms/benchmark_results_zasper_{k}kernels.json'
    with open(filename, 'r') as f:
        data = json.load(f)
    (avg_cpu, avg_memory, max_cpu, max_memory) = compute_averages(data)
    avg_cpus_zasper.append(avg_cpu)
    avg_mems_zasper.append(avg_memory)
    max_cpus_zasper.append(max_cpu)
    max_mems_zasper.append(max_memory)

for k in kernel_counts:
    filename = f'data/{delay}ms/benchmark_results_jupyter_{k}kernels.json'
    with open(filename, 'r') as f:
        data = json.load(f)
    (avg_cpu_jupyter, avg_memory_jupyter,  max_cpu_jupyter, max_memory_jupyter) = compute_averages(data)
    avg_cpus_jupyter.append(avg_cpu_jupyter)
    avg_mems_jupyter.append(avg_memory_jupyter)
    max_cpus_jupyter.append(max_cpu_jupyter)
    max_mems_jupyter.append(max_memory_jupyter)

# Plotting
fig, ((ax1, ax2), (ax3, ax4)) = plt.subplots(2, 2, figsize=(10, 10))

# 1. CPU Usage
ax1.plot(kernel_counts, avg_cpus_zasper, label="Zasper CPU Usage", marker='o', color='#583BD8')
ax1.plot(kernel_counts, avg_cpus_jupyter, label="Jupyter Server CPU Usage", marker='o', color='#E46E2E')
ax1.set_title('Zasper v/s Jupyter: Average CPU Usage (%)')
ax1.set_xlabel('Number of Kernels')
ax1.set_ylabel('CPU Usage (%)')
ax1.legend()
ax1.grid(True)
add_note(ax1, "Lower CPU usage is better", position='center right')

ax2.plot(kernel_counts, max_cpus_zasper, label="Zasper CPU Usage", marker='o', color='#583BD8')
ax2.plot(kernel_counts, max_cpus_jupyter, label="Jupyter Server CPU Usage", marker='o', color='#E46E2E')
ax2.set_title('Zasper v/s Jupyter: Max CPU Usage (%)')
ax2.set_xlabel('Number of Kernels')
ax2.set_ylabel('CPU Usage (%)')
ax2.legend()
ax2.grid(True)
add_note(ax2, "Lower CPU usage is better", position='center right')

# 2. RAM Usage
ax3.plot(kernel_counts, avg_mems_zasper, label="Zasper Memory Usage", marker='o', color='#583BD8')
ax3.plot(kernel_counts, avg_mems_jupyter, label="Jupyter Server Memory Usage", marker='o', color='#E46E2E')
ax3.set_title('Zasper v/s Jupyter: Average Memory Usage (MB)')
ax3.set_xlabel('Number of Kernels')
ax3.set_ylabel('Memory (MB)')
ax3.legend()
ax3.grid(True)
add_note(ax3, "Lower RAM usage is better", position='center right')

ax4.plot(kernel_counts, max_mems_zasper, label="Zasper Memory Usage", marker='o', color='#583BD8')
ax4.plot(kernel_counts, max_mems_jupyter, label="Jupyter Server Memory Usage", marker='o', color='#E46E2E')
ax4.set_title('Zasper v/s Jupyter: Max Memory Usage (MB)')
ax4.set_xlabel('Number of Kernels')
ax4.set_ylabel('Memory (MB)')
ax4.legend()
ax4.grid(True)
add_note(ax4, "Lower RAM usage is better", position='center right')


fig.text(
    0.5, 0.02,
    "* Lower CPU and RAM usage indicates better performance.",
    ha='center',
    fontsize=10,
    bbox=dict(facecolor='white', edgecolor='gray', boxstyle='round,pad=0.4')
)

# Show the plot
# plt.tight_layout()
plt.tight_layout(rect=[0, 0.05, 1, 1])
plt.savefig(f"plots/{delay}ms/summary_resources.png")
# plt.show()
