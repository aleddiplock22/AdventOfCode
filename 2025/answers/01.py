from pathlib import Path

with Path("../data/1/input.txt").open("r") as f:
    data = [line.strip() for line in f.readlines()]

with Path("../data/1/test.txt").open("r") as f:
    data_test = [line.strip() for line in f.readlines()]


def part1(data: list[str]) -> int:
    pos = 50
    counter = 0
    for line in data:
        if line[0] == 'R':
            pos += int(line[1:])
        elif line[0] == 'L':
            pos -= int(line[1:])
        else:
            raise RuntimeError(line)
        yes = False
        while pos < 0:
            pos += 100
        while pos > 99:
            pos -= 100
        if pos == 0:
            counter += 1
    return counter

def part2(data: list[str]) -> int:
    pos = 50
    counter = 0
    for line in data:
        started_zero = pos == 0
        if line[0] == 'R':
            pos += int(line[1:])
        else:
            pos -= int(line[1:])
        yes = False
        while pos < 0:
            if started_zero:
                started_zero = False
            else:
                counter += 1
            pos += 100
        while pos > 99:
            counter += 1
            yes = True
            pos -= 100
        if pos == 0 and not yes:
            counter += 1
    return counter

if __name__ == "__main__":
    assert part1(data_test) == 3
    print(f"Part 1: {part1(data)}")  # 999
    assert (ans := part2(data_test)) == 6, ans
    print(f"Part 2: {part2(data)}")  # 6099