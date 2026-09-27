def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


print("--- SIMPLE CALC ---")
check(2 + 2, 4)
check(2 + 3 * 4, 14)
check((2 + 3) * 4, 20)
check(-10, -10)
check(5 - 3 - 1, 1)
check(1 + 2 * 3 - 4, 3)
check(-1+2, 1)
check(8 / 2, 4)
check((5), 5)

print("--- STRINGS ---")
check("hola", "hola")
check("hola" + "chau", "holachau")

print("--- BOOLEAN LOGIC ---")
check(True or False, True)
check(True and False, False)
check(not True, False)
check(True or None, True)
check(True and None, None)
check((1 - (2 * 3) < 4) == True, True)
