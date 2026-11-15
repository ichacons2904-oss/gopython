def check(value, expected):
    if value == expected:
        print("OK")
    else:
        print("ERROR")

check((), tuple())
check((1,), tuple([1]))
check((1), 1)
values = 1, 2, 3
check(values[-1], 3)
check(len(values), 3)
check(2 in values, True)
check(4 not in values, True)
check(values + (4,), (1, 2, 3, 4))
check(2 * (1,), (1, 1))
check(list(values), [1, 2, 3])
check(tuple({"a": 1, "b": 2}), ("a", "b"))
check(not (), True)
def pair():
    return 3, 4
check(pair(), (3, 4))
total = 0
for value in 1, 2, 3:
    total = total + value
check(total, 6)
items = [1]
container = (items,)
items.append(2)
check(container[0], [1, 2])
