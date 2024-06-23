import subprocess
import csv
import pandas as pd
import matplotlib
matplotlib.use('Agg')  # This must be done before importing pyplot
import matplotlib.pyplot as plt
import os

# Definitions
file_sizes = ['xsmall', 'small', 'medium']
# file_sizes = ['xsmall', 'small', 'medium', 'large']
thread_counts = [1, 2, 4, 6, 8, 12]
runs = 3
cwd = os.getcwd()
go_script = cwd + '/simulation'  # Adjust with the actual Go script name

# Function to run the Go script command and handle multiple runs
def run_command(command):
    """Runs the given command multiple times and averages the output."""
    seq_results = []
    par_results = []
    for _ in range(runs):
        result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        if result.returncode == 0:
            sequential_output = result.stdout.strip().split(',')[0].split()[-1]
            parallel_output = result.stdout.strip().split(',')[-1].split()[-1]
            seq_results.append(int(sequential_output))
            par_results.append(int(parallel_output))
        else:
            print(f"Error: {result.stderr}")
            return None, None  # Ensure function always returns a tuple
    return sum(seq_results) // len(seq_results), sum(par_results) // len(par_results)

def main():
    for size in file_sizes:
        data = []
        seq_time, _ = run_command(['go', 'run', go_script, size])  # Handle None returns
        if seq_time is None:
            continue
        print(f"Sequential for {size} done in {seq_time} microseconds")

        for threads in thread_counts:
            p_partime, p_seqtime = run_command(['go', 'run', go_script, size, str(threads), 'p'])
            q_partime, q_seqtime = run_command(['go', 'run', go_script, size, str(threads), 'q'])
            if p_partime is None or q_partime is None:
                continue  # Skip this iteration if the command failed
            data.append([threads, seq_time, p_partime, p_seqtime, q_partime, q_seqtime])
            print(f"{threads} for {size} done. Parallel: {p_partime, p_seqtime}, Q Parallel: {q_partime, q_seqtime}")

        if not data:
            continue  # Skip file and plot generation if data collection failed

        # Save data to CSV
        csv_filename = f"{size}_timings.csv"
        with open(csv_filename, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(['Threads', 'Sequential', 'P Sequential', 'P Parallel', 'Q Sequential', 'Q Parallel'])
            writer.writerows(data)

        # Plotting
        df = pd.DataFrame(data, columns=['Threads', 'Sequential', 'P Sequential', 'P Parallel', 'Q Sequential', 'Q Parallel'])
        plot_graphs(df, size)

def plot_graphs(df, size):
    # Theoretical Speedup Graphs
    plt.figure()
    p_parallel_prop = df['P Parallel'].iloc[0] / (df['P Sequential'].iloc[0] + df['P Parallel'].iloc[0])
    q_parallel_prop = df['Q Parallel'].iloc[0] / (df['Q Sequential'].iloc[0] + df['Q Parallel'].iloc[0])
    df['P Speedup'] = 1 / ((1 - p_parallel_prop) + p_parallel_prop / df['Threads'])
    df['Q Speedup'] = 1 / ((1 - q_parallel_prop) + q_parallel_prop / df['Threads'])
    plt.plot(df['Threads'], df['P Speedup'], label='Parallel')
    plt.plot(df['Threads'], df['Q Speedup'], label='Work-Queue Parallel')
    plt.title(f'Theoretical Speedup Comparison for {size}')
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.legend()
    plt.savefig(f"{size}_theoretical_speedup.png")
    plt.close()

    # Actual Speedup Graphs
    plt.figure()
    df['P Actual Speedup'] = df['Sequential'] / (df['P Parallel'] + df['P Sequential'])
    df['Q Actual Speedup'] = df['Sequential'] / (df['Q Parallel'] + df['Q Sequential'])
    plt.plot(df['Threads'], df['P Actual Speedup'], label='Parallel')
    plt.plot(df['Threads'], df['Q Actual Speedup'], label='Work-Queue Parallel')
    plt.title(f'Actual Speedup Comparison for {size}')
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.legend()
    plt.savefig(f"{size}_actual_speedup.png")
    plt.close()

    # # Total Time Graphs
    # plt.figure()
    # plt.plot(df['Threads'], [df['Sequential'].iloc[0]] * len(df['Threads']), label='Sequential')
    # plt.plot(df['Threads'], df['P Sequential'] + df['P Parallel'], label='Parallel')
    # plt.plot(df['Threads'], df['Q Sequential'] + df['Q Parallel'], label='Work-Queue Parallel')
    # plt.title(f'Total Time Comparison for {size}')
    # plt.xlabel('Number of Threads')
    # plt.ylabel('Time in Microseconds')
    # plt.legend()
    # plt.savefig(f"{size}_total_time.png")
    # plt.close()

if __name__ == '__main__':
    main()
