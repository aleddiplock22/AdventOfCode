from pathlib import Path

with Path("../data/2/test.txt").open("r") as f:
    raw_test_data = f.read()

with Path("../data/2/input.txt").open("r") as f:
    raw_data = f.read()

def part1(data: str):
    return sum(map(lambda z: sum(y for y in z if len(str(y)) % 2 == 0 and str(y)[len(str(y))//2:] == str(y)[:len(str(y))//2]), map(lambda x: range(x[0], x[1]+1), map(lambda x: tuple(map(int, (x.split("-")))), data.split(",")))))

def part2(data: str):
    data = tuple(map(int, (x.split("-"))) for x in data.split(","))

    total = 0
    for s, e in data:
        for val in range(s, e+1):
            val_str = str(val)
            length = len(val_str)
            window = 1

            for window in range(1, length // 2+1):
                if not (length % window == 0):
                    continue
                if val_str[:window] * (length // window) == val_str:
                    total += val
                    break
    return total


if __name__ == "__main__":
    assert part1(raw_test_data) == 1227775554
    print("Part 1:", part1(raw_data)) # 17077011375

    assert part2(raw_test_data) == 4174379265
    print("Part 2:", part2(raw_data)) # 36037497037