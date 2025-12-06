from pathlib import Path
from math import prod

with Path("../data/6/test.txt").open("r") as f:
    test_data = [line.strip().split() for line in f.readlines()]

with Path("../data/6/input.txt").open("r") as f:
    data = [line.strip().split() for line in f.readlines()]

    
with Path("../data/6/test.txt").open("r") as f:
    raw_test_data = [line.removesuffix("\n") for line in f.readlines()]

with Path("../data/6/input.txt").open("r") as f:
    raw_data = [line.removesuffix("\n") for line in f.readlines()]

def solve_problem(problem: list[str]) -> int:
    method = problem[-1]
    # this additional n != "" logic is because of wrong assumption in part2 
    # that all numbers are L=4 digits long, when many are L-2=2 or L-1=3 digits long.
    nums = map(int, [n for n in problem[:-1] if n != ""])
    match method:
        case "*":
            return prod(nums)
        case "+":
            return sum(nums)

def part1(data):
    problems = list(zip(*data))
    return sum(map(solve_problem, problems))

def part2(data):
    problems = []
    L = len(data)
    current_problem = [''] * L
    idx = 0
    for line in zip(*[[*line] for line in data]):
        if all(char == " " for char in line):
            problems.append(current_problem)
            current_problem = [''] * L
            idx = 0
            continue
        if idx == 0:
            assert line[L-1] in ("+", "*")
            current_problem[L-1] = line[L-1]

        current_problem[idx] += "".join(line[:L-1])
        idx += 1
    
    # add the final problem manually
    problems.append(current_problem)
    
    return sum(map(solve_problem, problems))

if __name__ == "__main__":
    test1 = part1(test_data)
    ans1 = part1(data)
    print(f"[Part 1] Test: {test1} | Answer: {ans1}")

    test2 = part2(raw_test_data)
    ans2 = part2(raw_data)
    print(f"[Part 2] Test: {test2} | Answer: {ans2}")