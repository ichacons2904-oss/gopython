def check(value, expected):
    if value == expected:
        print("OK")
    else:
        print("ERROR")

check(2 not in [1], True)
check(1 not in [1], False)
check(2 not in [1] == [2], False)
check(not 2 in [1] == [2], True)
check(list(), [])
check(list(range(3)), [0, 1, 2])

items = [1, 2, 1]
alias = items
items.insert(-1, 3)
check(alias, [1, 2, 3, 1])
items.remove(1)
check(items, [2, 3, 1])

copy = list(items)
copy.insert(0, 9)
check(items, [2, 3, 1])
check(copy, [9, 2, 3, 1])
