EXAMPLE_FILEPATH = "day8-example.txt"
INPUT_FILEPATH = "day8-input.txt"

def parse_input(filepath: str) -> list[list[int]]:
    with open(filepath, "r") as f:
        return [
            [int(chr) for chr in line.strip()]
            for line in f.readlines()
        ]

def part1(filepath: str) -> int:
    grid = parse_input(filepath)
    count = 0
    for r in range(0, len(grid)):
        for c in range(0, len(grid[0])):
            height = grid[r][c]
            if (
                all(grid[r][x] < height for x in range(c)) or  # All trees to the left
                all(grid[y][c] < height for y in range(r)) or  # All trees above
                all(grid[r][x] < height for x in range(c+1, len(grid[r]))) or  # all trees to the right
                all(grid[y][c] < height for y in range(r+1, len(grid)))  # all trees below
            ):
                count += 1
    return count

def find_scenic_score_of_point(row: int, col: int, grid: list[list[int]]) -> int:
    height = grid[row][col]

    left = 0
    for c in range(col-1, -1, -1):  # backwards stepping from left of point!
        left += 1
        if not (grid[row][c] < height):
            break
    
    right = 0
    for c in range(col+1, len(grid[row])):
        right += 1
        if not (grid[row][c] < height):
            break

    above = 0
    for r in range(row-1, -1, -1):
        above += 1
        if not (grid[r][col] < height):
            break
    
    below = 0
    for r in range(row+1, len(grid)):
        below += 1
        if not (grid[r][col] < height):
            break

    return left * right * above * below

def part2(filepath: str):
    grid = parse_input(filepath)
    return max(
        find_scenic_score_of_point(r, c, grid)
        for r in range(0, len(grid)) 
        for c in range(0, len(grid[r]))
    )

if __name__ == "__main__":
    example_p1 = part1(EXAMPLE_FILEPATH)
    assert example_p1 == 21, f"Got {example_p1} for Part 1 Example :("
    part1_ans = part1(INPUT_FILEPATH)

    example_p2 = part2(EXAMPLE_FILEPATH)
    assert example_p2 == 8, f"Got {example_p2} for Part 2 Example :("
    part2_ans = part2(INPUT_FILEPATH)

    print("---AoC 2022 Day 8---")
    print(f"Part 1: {part1_ans}")
    print(f"Part 2: {part2_ans}")
