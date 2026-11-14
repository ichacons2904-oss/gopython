def check(value, expected):
    if value == expected:
        print("OK")
    else:
        print("ERROR")

class Counter:
    def __init__(self, initial):
        self.value = initial
        self.items = []
    def add(self, amount):
        self.value = self.value + amount
        self.items.append(amount)
        return self.value

class Child(Counter): pass
first = Child(2)
second = Child(10)
saved = first.add
check(saved(3), 5)
check(second.value, 10)
check(first.items[0], 3)
check(Counter.add(second, 4), 14)
def replacement(self, amount):
    return 100
Counter.add = replacement
first = Child(0)
check(saved(1), 6)
check(first.add(1), 100)
first.add = replacement
check(first.add(second, 1), 100)
