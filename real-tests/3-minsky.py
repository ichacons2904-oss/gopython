def minsky_even(n):
    r0 = n
    pc = 0
    halt = False

    while True and not halt:
        # Instruction 0: DECJZ(r0, pc2, pc1)
        if pc == 0:
            if r0 == 0:
                pc = 2
            else:
                r0 = r0 - 1
                pc = 1

        # Instruction 1: DECJZ(r0, pc3, pc0)
        elif pc == 1:
            if r0 == 0:
                pc = 3
            else:
                r0 = r0 - 1
                pc = 0

        # Instruction 2: INC(r0, pc3)
        elif pc == 2:
            r0 = r0 + 1
            pc = 3

        # Instruction 3: HALT
        elif pc == 3:
            halt = True

    return r0


def check(value, expected):
    if value != expected:
        print("ERROR")
    else:
        print("OK")


print("--- MINSKY MACHINE ---")
check(minsky_even(0), 1)
check(minsky_even(1), 0)
check(minsky_even(2), 1)
check(minsky_even(3), 0)
check(minsky_even(4), 1)
check(minsky_even(5), 0)
