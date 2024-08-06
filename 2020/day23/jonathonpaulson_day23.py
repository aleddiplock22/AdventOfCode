import fileinput
import re
import sys

#start = '389125467'
start = '624397158'

def solve():
    n = int(1e6)
    N = [None for i in range(n+1)]
    X = [int(x) for x in start]
    for i in range(len(X)):
        N[X[i]] = X[(i+1)%len(X)]
    
    N[X[-1]] = len(X)+1
    for i in range(len(X)+1, n+1):
        N[i] = i+1
    N[-1] = X[0]

    current = X[0]
    nmoves = int(1e7)
    for _ in range(nmoves):
        pickup = N[current]
        N[current] = N[N[N[pickup]]]

        dest = n if current==1 else current-1
        while dest in [pickup, N[pickup], N[N[pickup]]]:
            dest = n if dest==1 else dest-1

        N[N[N[pickup]]] = N[dest]
        N[dest] = pickup
        current = N[current]

    return N[1]*N[N[1]]
    

print(solve())