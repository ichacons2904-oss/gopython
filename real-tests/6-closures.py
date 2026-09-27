def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


print("--- NESTED FUNCTIONS ---")


def make_counter():
    i = 0

    def count():
        nonlocal i
        i = i + 1
        return i

    return count


x_counter = make_counter()
y_counter = make_counter()
check(x_counter(), 1)
check(x_counter(), 2)
check(x_counter(), 3)
check(y_counter(), 1)
check(y_counter(), 2)
check(y_counter(), 3)
check(x_counter(), 4)
check(y_counter(), 4)

print("--- CLOSURE READS ENCLOSING SCOPE ---")


def outer():
    a = "outer"

    def ret_a():
        return a

    check(ret_a(), "outer")
    a = "changed"
    check(ret_a(), "changed")


outer()
