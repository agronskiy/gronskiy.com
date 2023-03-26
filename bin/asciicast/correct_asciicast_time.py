#!/usr/bin/env python3

import json
import sys

import click


@click.command()
@click.option("--input", required=True, help="Input file")
@click.option(
    "--time", required=True, help="Seconds to add (positive) or remove (negative)"
)
@click.option("--output", required=True, help="Output file.")
def main(input, time, output):
    lines = []
    with open(input, "r") as f:
        for line in f:
            json_line = json.loads(line)
            if not isinstance(json_line, list):
                lines.append(line)
                continue
            if len(json_line) != 3:
                print(f"Unexpected number of array elements: {json_line}\n")
                lines.append(line)
                continue
            new_time = float(json_line[0]) + float(time)
            json_line[0] = new_time
            lines.append(json.dumps(json_line))

    with open(output, "w") as f:
        f.write("\n".join(lines))


if __name__ == "__main__":
    sys.exit(main())
