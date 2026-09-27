def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


print("--- SCOPE HANDLING ---")
x = 0
y = 0
check(x, 0)
check(y, 0)


def first_block():
    z = 1
    check(z, 1)


def second_block():
    z = 2
    check(z, 2)


def outer_block():
    x = 1
    check(x, 1)
    check(y, 0)
    first_block()
    second_block()


outer_block()
check(x, 0)
check(y, 0)

print("--- SIMPLE FUNCTION USAGE ---")


def fib(n):
    if n <= 1:
        return n
    return fib(n - 2) + fib(n - 1)


def noreturn():
    x = 1


def emptyreturn():
    return


check(fib(20), 6765)
check(noreturn(), None)
check(emptyreturn(), None)

print("--- VARIABLE SHADOWING ---")
a = 1
check(a, 1)


def shadow_inner():
    a = 3
    check(a, 3)


def shadow_outer():
    a = 2
    check(a, 2)
    shadow_inner()
    check(a, 2)


shadow_outer()
check(a, 1)

print("--- LATE BINDING ---")
a = "global"


def ret_a():
    return a


check(ret_a(), "global")
a = "reassigned"
check(ret_a(), "reassigned")
