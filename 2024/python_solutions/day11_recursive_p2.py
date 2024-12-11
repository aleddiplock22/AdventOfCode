"""
Credit: Ryan
"""


from functools import cache
from datetime import datetime

times = []


for _ in range(10):
    t1 = datetime.now()

    with open("../inputs/11/input.txt") as f:
        stones = [int(x) for x in f.read().split()]

    @cache
    def recurse_blink(stones, steps_to_go):
        if steps_to_go == 0:
            return len(stones)
        return sum(recurse_blink(modify_stone(stone), steps_to_go-1) for stone in stones)

    def modify_stone(stone):
        if stone == 0:
            return (1,)
        elif len(str(stone)) % 2 == 0:
            l = len(str(stone))
            return (stone // (10**(l//2)), stone % (10**(l//2)))
        else:
            return (stone*2024,)
    
    t2 = datetime.now()
    times.append((t2-t1).microseconds)

print((sum(times)/10)/1000)
