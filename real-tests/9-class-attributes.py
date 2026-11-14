def check(value, expected):
    if value == expected:
        print("OK")
    else:
        print("ERROR")

outer = 4
class C:
    value = outer + 1
    def read(self):
        return outer

check(C.value, 5)
check(C.read(None), 4)
C.value = 6
check(C.value, 6)

class Parent:
    value = 1
class Left(Parent): pass
class Right:
    other = 2
class Child(Left, Right): pass
check(Child.value, 1)
check(Child.other, 2)
Parent.value = 3
check(Child.value, 3)

def make():
    value = 7
    class Inner:
        value = 9
        def read(self):
            return value
    return Inner

Inner = make()
check(Inner.read(None), 7)
check(Inner == Inner, True)

def change_enclosing():
    value = 1
    class Inner:
        nonlocal value
        value = 2
    return value

check(change_enclosing(), 2)
