def main() -> None:
    with open("./inputs/2/input.txt") as f:
        data = [list(map(int, line.split(" "))) for line in f.readlines()]
    part1 = sum(map(lambda nums: all(1 <= right - left <= 3 for left, right in zip(nums, nums[1:])) or all(-1 >= right - left >= -3 for left, right in zip(nums, nums[1:])), data))
    part2 = sum(map(lambda nums: any(map(lambda n: all(1 <= right - left <= 3 for left, right in zip(n, n[1:])) or all(-1 >= right - left >= -3 for left, right in zip(n, n[1:])), [nums[:i] + nums[i+1:] for i in range(len(nums))])), data))
    print(f"{part1},{part2}")

if __name__ == "__main__":
    main()