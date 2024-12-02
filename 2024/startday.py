import shutil
import sys

from pathlib import Path

"""
Usage:

> python startday.py 2

Result:

Creates 02.go, and example.txt and input.txt files in inputs/2/
"""


day = sys.argv[1]
if len(day) == 1:
    #single digit
    day = "0" + day

daylines = f'''package main

func day{day}(part2 bool) Solution {{
	if !part2 {{
		return Solution{{
			"{day}",
			"example part 1",
			"input part 1",
		}}
	}} else {{
		return Solution{{
			"{day}",
			"example part 2",
			"input part 2",
		}}
    }}
}}
'''

if __name__ == "__main__":
    assert Path.cwd().stem == "2024", "you need to `cd 2024`"
    with Path(f"{day}.go").open("w") as f:
        f.write(daylines)
    
    inputs_dir = Path.cwd().joinpath("inputs").joinpath(sys.argv[1])
    inputs_dir.mkdir(exist_ok=True)

    with inputs_dir.joinpath("example.txt").open("w") as f:
        f.write("")

    with inputs_dir.joinpath("input.txt").open("w") as f:
        f.write("")

    print(f"\nNow go delete day{day} func from placeholder.go...")

