def make_counter():
    count = 0

    def increment():
        nonlocal count
        count = count + 1
        return count

    return increment


counter = make_counter()
result = 0
for i in range(2000000):
    result = counter()
print(result)
