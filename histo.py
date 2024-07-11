#!/usr/bin/env python3
import sys
import matplotlib.pyplot as plt

def main():
    timings = []

    # Read timings from stdin
    for line in sys.stdin:
        try:
            line = line.strip()
            if len(line) == 0:
                continue

            # Convert each line to an integer (milliseconds) and append to the list
            timings.append(int(line.strip()))

        except Exception as e: print(e)

    # Create a histogram
    plt.hist(timings, bins=30, edgecolor='black')
    plt.title('Histogram of Millisecond Timings')
    plt.xlabel('Time (ms)')
    plt.ylabel('Frequency')

    # Display the histogram
    plt.show()

if __name__ == "__main__":
    main()

