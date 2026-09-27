def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


print("--- FLOATS ---")
check(5 + 0.5, 5.5)
check(1.5 * 2, 3)
check(8 / 2, 4.0)
check(7 / 2, 3.5)
check(-7 / 2, -3.5)

print("--- MODULO ---")
check(5 % 2, 1)
check(-7 % 3, 2)
check(7 % -3, -2)
