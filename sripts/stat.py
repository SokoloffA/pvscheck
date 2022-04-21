#!/usr/bin/env python3

import sys

if __name__ == "__main__":
    stats = {}

    src = open(sys.argv[1], 'r')
    for line in src:
        if line == "":
            continue

        print(line)

        (file, line, msgType, message) = line.split("\t")

        print("----------------")
        print(message)
        print(message.split(' ', 1))
        print("----------------")
        (vNum, message) = message.split(' ', 1)

        try:
            stats[vNum] += 1
        except KeyError:
            stats[vNum] = 1

    src.close()

    data = {k: v for k, v in sorted(stats.items(), key=lambda item: item[1], reverse=True)}
    for num, count in data.items():
        print(f" * {num}: {count}")
    #keys = stats.keys().sorted()
