import ast
import sys
from pathlib import Path

EXAMPLE_FILEPATH = Path("example.txt")
INPUT_FILEPATH = Path("input.txt")

def read_input(filepath: Path):
    with filepath.open("r") as f:
        contents = f.read()
    contents = contents.split("\n\n")
    contents = tuple(
        tuple(ast.literal_eval(y) for y in x.split("\n")) for x in contents
    )
    return contents

def compare_packets(p1, p2):
    if isinstance(p1, list) and isinstance(p2, list):
        if len(p1) == len(p2):
            for x,y in zip(p1, p2):
                if isinstance(x, list) or isinstance(y, list):
                    if compare_packets(x, y):
                        continue
                    else:
                        return False
                if x < y:
                    return True
                if x > y:
                    return False
            return True

        if len(p1) < len(p2):
            # If we run out of comparisons, LHS ran out of elements first, so we're in order
            return compare_packets(p1, p2[:len(p1)])
        if len(p1) > len(p2):
            for x,y in zip(p1, p2):
                if isinstance(x, list):
                    return compare_packets(x, [y])
                if isinstance(y, list):
                    return compare_packets([x], y)
                if x < y:
                    return True
                if y > x:
                    return False
            # RHS ran out elements to compare, so we're out of order 
            return False
    
    if isinstance(p1, list) and not isinstance(p2, list):
        return compare_packets(p1, [p2])
    
    if not isinstance(p1, list) and isinstance(p2, list):
        return compare_packets([p1], p2)
    
    raise AssertionError(f"Shouldn't get here? {p1=} {p2=}")

def part1(filepath: Path):
    inputs = read_input(filepath)

    count = 0
    for i, (p1, p2) in enumerate(inputs, start=1):
        if compare_packets(p1, p2):
            count += i
    return count

if __name__ == "__main__":
    sys.setrecursionlimit(9999)
    print(f"[P1 Example]\n\tExpected: 13\n\tAnswer: {part1(EXAMPLE_FILEPATH)}")
    print(f"[Part 1]\n\tAnswer: {part1(INPUT_FILEPATH)}")  # rip answer is too low, I do not enjoy these ones enough to fix :|