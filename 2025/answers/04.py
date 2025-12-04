from pathlib import Path

data_dir = Path("../data/4/")
with data_dir.joinpath("input.txt").open("r") as f:
    grid = [[*line.strip()] for line in f.readlines()]

with data_dir.joinpath("test.txt").open("r") as f:
    grid_test = [[*line.strip()] for line in f.readlines()]

def part1(grid: list[list[str]]):
    R = len(grid)
    C = len(grid[0])
    total = 0
    for r in range(R):
        for c in range(C):
            if grid[r][c] != "@":
                continue
            nbrs = 0
            for dr, dc in [(-1, 0), (-1, 1), (-1, -1), (1, 0), (1, -1), (1, 1), (0, 1), (0, -1)]:
                nr, nc = r + dr, c + dc
                if not (nr >= 0 and nc >= 0 and nr < R and nc < C):
                    continue
                if grid[nr][nc] == "@":
                    nbrs += 1
                if nbrs > 3:
                    break
            if nbrs < 4:
                total += 1
    return total


def part2(grid: list[list[str]]):
    R = len(grid)
    C = len(grid[0])
    total = 0
    while True:
        change = False
        for r in range(R):
            for c in range(C):
                if grid[r][c] != "@":
                    continue
                nbrs = 0
                for dr, dc in [(-1, 0), (-1, 1), (-1, -1), (1, 0), (1, -1), (1, 1), (0, 1), (0, -1)]:
                    nr, nc = r + dr, c + dc
                    if not (nr >= 0 and nc >= 0 and nr < R and nc < C):
                        continue
                    if grid[nr][nc] == "@":
                        nbrs += 1
                    if nbrs > 3:
                        break
                if nbrs < 4:
                    change = True
                    grid[r][c] = "."
                    total += 1
        if not change:
            break
    return total

if __name__ == "__main__":
    print("P1 Test:", part1(grid_test))
    print("P1:", part1(grid))

    print("P2 Test:", part2(grid_test))
    print("P2:", part2(grid))