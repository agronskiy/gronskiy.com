#!/usr/bin/env python3

import glob
import json
import os
import string
import sys
from pathlib import Path

import click
from fontTools import subset


def _get_cast_chars(f):
    for line in f:
        json_line = json.loads(line)
        if not isinstance(json_line, list):
            continue
        if len(json_line) != 3:
            print(f"Unexpected number of array elements: {json_line}\n")
            continue
        for char in json_line[2]:
            yield char


@click.command()
@click.option("--input-font", required=True, help="Path to input font file.")
@click.option("--output", required=True, help="File to output subsetted font to.")
def main(input_font, output):
    # Build list of unicode codepoints that the asciinema casts require
    cast_chars = set()
    for cast_file in glob.iglob(
        os.path.join(Path(__file__).absolute().parent.parent.parent, "**", "*.cast"),
        recursive=True,
    ):
        print("Addint file", cast_file)
        with open(cast_file, "r") as f:
            cast_chars = cast_chars.union(set(_get_cast_chars(f)))

    # Exclude ascii chars (0-127)? If the font you are including is only
    # for glyphs then you don't really need them.
    cast_chars -= {
        c
        for chars in [string.ascii_letters, string.digits, string.punctuation]
        for c in chars
    }

    # Exclude whitespace, fonts don't register these characters so
    # leaving them in here will break fontTools.subset
    cast_chars -= {c for c in string.whitespace}
    return subset.main(
        args=[
            input_font,
            "--unicodes={}".format(
                ",".join(cast_chars)
                .encode("unicode-escape")
                .decode("ascii")
                .replace("\\u", "U+")
            ),
            f"--output-file={output}",
            "--timing",
        ]
    )


if __name__ == "__main__":
    sys.exit(main())
