def check(value, expected):
    if value == expected:
        print("OK")
    else:
        print("ERROR")

ages = {"ana": 30, "beto": 25}
check(ages["ana"], 30)
ages["carla"] = 41
ages["ana"] = 31
check(ages, {"ana": 31, "beto": 25, "carla": 41})
check(len(ages), 3)
check("beto" in ages, True)
check(25 in ages, False)
check("dario" not in ages, True)

check(ages.get("dario"), None)
check(ages.get("dario", 0), 0)
check(list(ages.keys()), ["ana", "beto", "carla"])
check(list(ages.values()), [31, 25, 41])
check(ages.pop("beto"), 25)
check(ages.pop("beto", -1), -1)

alias = ages
copy = dict(ages)
alias["eva"] = 19
check(ages["eva"], 19)
check("eva" in copy, False)

copy.update({"ana": 0, "fede": 50})
check(copy, {"ana": 0, "carla": 41, "fede": 50})

check({1: "uno"}[1.0], "uno")
check({}, dict())
check(not {}, True)

def count_words(words):
    counts = {}
    for word in words:
        counts[word] = counts.get(word, 0) + 1
    return counts

counts = count_words(["a", "b", "a", "c", "a", "b"])
check(counts, {"a": 3, "b": 2, "c": 1})

order = []
for key in counts:
    order.append(key)
check(order, ["a", "b", "c"])

nested = {"list": [1, 2], "dict": {"x": 1}}
nested["list"].append(3)
nested["dict"]["y"] = 2
check(nested, {"list": [1, 2, 3], "dict": {"x": 1, "y": 2}})

copy.clear()
check(len(copy), 0)
