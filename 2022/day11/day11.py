from typing import Callable
from math import prod

# MONKEYS
EXAMPLE_FILEPATH = "example.txt"
INPUT_FILEPATH = "input.txt"

global IS_PART_2
IS_PART_2 = False

DIVISORS = set()
global BIG_DIVISOR

class Monkey:
    def __init__(self, id: int, items: list[int], operation_definition: str, test_definition: str) -> None:
        self.id = id
        self.items: list[int] = items
        self.operation: Callable[[int], int] = self.parse_operation(operation_definition)
        self.test: Callable[[int], int] = self.parse_test(test_definition)
        self.counter = 0
    
    def parse_operation(self, definition: str) -> Callable[[int], int]:
        if "*" in definition:
            val = definition.split("* ")[-1]
            try:
                val = int(val)
            except ValueError:
                assert "* old" in definition
                val = None

            def f(x: int) -> int:
                if val is None:
                    return x * x
                return x * val
            
            return f

        if "+" in definition:
            val = definition.split("+ ")[-1]
            try:
                val = int(val)
            except ValueError:
                assert "+ old" in definition
                val = None
            
            def f(x: int) -> int:
                if val is None:
                    return x + x
                return x + val
            
            return f

    def parse_test(self, definition: str) -> Callable[[int], int]:
        lines = (l.strip() for l in definition.split("\n"))
        lines = iter(lines)
        div_by = int(next(lines).strip().split("divisible by ")[-1])
        DIVISORS.add(div_by)
        if_true = int(next(lines).strip().split("throw to monkey ")[-1])
        if_false = int(next(lines).strip().split("throw to monkey ")[-1])
        def f(x: int) -> int:
            if x % div_by == 0:
                return if_true
            return if_false
        return f

    def add_item(self, x: int):
        self.items.append(x)
    
    @property
    def is_empty(self):
        return len(self.items) == 0

    def pop_to(self, divisor = None):
        self.counter += 1
        item = self.items[0]
        self.items = self.items[1:]
        item = self.operation(item)
        if not IS_PART_2:
            item //= 3
        if divisor is not None:
            item %= divisor
        return item, self.test(item)


def parse_input(filepath: str) -> list[Monkey]:
    with open(filepath, "r") as f:
        monkey_blocks = f.read().split("\n\n")
    monkeys = []
    for block in monkey_blocks:
        lines = block.splitlines()
        id = int(lines[0].split("Monkey ")[-1][:-1])
        items = list(map(int, (lines[1].split("items: ")[-1].split(", "))))
        operation_defn = lines[2]
        test_defn = "\n".join(lines[3:])
        monkeys.append(Monkey(id, items, operation_defn, test_defn))
    
    return monkeys

def solve(filepath: str) -> int:
    N = 10000 if IS_PART_2 else 20
    monkeys: list[Monkey] = parse_input(filepath)
    BIG_DIVISOR = prod(DIVISORS)
    for _ in range(N):
        for monkey in monkeys:
            while not monkey.is_empty:
                item, to_monkey = monkey.pop_to(BIG_DIVISOR)
                monkeys[to_monkey].add_item(item)
    counters = [monkey.counter for monkey in monkeys]
    counters.sort(reverse=True)
    
    return counters[0]*counters[1]




if __name__ == "__main__":
    assert solve(EXAMPLE_FILEPATH) == 10605, solve(EXAMPLE_FILEPATH)
    print("Part 1:", solve(INPUT_FILEPATH))

    """
    For part 2 we use the fact that:
    (a + b) mod x = ( (a mod x) + (b mod x) ) mod x

    and:
    (a * b) mod x = ( (a mod x) * (b mod x) ) mod x 


    so because all our tests are just modulo checks ('is divisble by')
    and we have this associative property,
    we can keep our values scaled down as follows:
    
        divisor = product(all divisors in tests)  #  or least common multiple! but they're all prime here so same thing
        value = value % divisor
    
    thus our values stay capped by divisor, so we dont get massive numbers slowing everything down!
    """


    IS_PART_2 = True
    solve(EXAMPLE_FILEPATH) == 2713310158
    print("Part 2:", solve(INPUT_FILEPATH))