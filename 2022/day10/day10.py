example_instructions = [line.strip() for line in open("example.txt").read().splitlines()]
input_instructions = [line.strip() for line in open("input.txt").read().splitlines()]

NOOP = "noop"
ADDX = "addx "

def part1(instructions: list[str]):
    cycles = 0
    signal_strength_total = 0
    x = 1
    for instruction in instructions:
        cycles += 1

        if (cycles - 20) % 40 == 0:
            signal_strength_total += x*cycles

        if instruction == NOOP:
            continue
            
        val = int(instruction.split(ADDX)[-1])
        cycles += 1
        if (cycles - 20) % 40 == 0:
            signal_strength_total += x*cycles
        x += val
    
    return signal_strength_total

def part2(instructions: list[str]):
    cycles = 0
    current_row = ""
    rows: list[str] = []
    x = 1
    for instruction in instructions:
        cycles += 1
        if cycles in (x, x+1, x+2):
            current_row += "#"
        else:
            current_row += "."

        if cycles % 40 == 0:
            cycles = 0
            rows.append(current_row)
            current_row = ""

        if instruction == NOOP:
            continue
            
        val = int(instruction.split(ADDX)[-1])
        cycles += 1
        if cycles in (x, x+1, x+2):
            current_row += "#"
        else:
            current_row += "."
        if cycles % 40 == 0:
            rows.append(current_row)
            cycles = 0
            current_row = ""
        x += val
    
    for row in rows:
        print(row)

assert part1(example_instructions) == 13140
print(part1(input_instructions))

part2(example_instructions)
print()
part2(input_instructions)
