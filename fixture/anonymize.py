#!/usr/bin/env python3
#!/usr/bin/env python3
import json, sys, io
from hashlib import md5
from functools import reduce 

list0 = []
with open("fixture/list0.txt") as fp:
    for line in fp:
        list0.append(line.strip())
list1 = []
with open("fixture/list1.txt") as fp:
    for line in fp:
        list1.append(line.strip())
list2 = []
with open("fixture/list2.txt") as fp:
    for line in fp:
        list2.append(line.strip())
list3 = []
with open("fixture/list3.txt") as fp:
    for line in fp:
        list3.append(line.strip())


# A injective and suriective function...
def assoc(p: str, mod: int):
    digest = reduce(lambda x, y: (x+y) % mod, map(lambda x: x[0]*256 + ord(x[1]), enumerate(list(md5(p.encode("ascii")).hexdigest()))))
    return digest

data = json.load(sys.stdin)

try:
    data["Hostname"] = list0[assoc(data["Hostname"], len(list0))]
except Exception as ex:
    pass
try:
    data["Info"]["Hostname"] = list0[assoc(data["info"]["Hostname"], len(list0))]
except Exception as ex:
    pass
try:
    data["Databases"] = " ".join(map(lambda x: list2[assoc(x, len(list2))],  data["Databases"].split(" ")))
except Exception as ex:
    pass
try:
    data["Schemas"] = " ".join(map(lambda x: list2[assoc(x, len(list2))],  data["Schemas"].split(" ")))
except Exception as ex:
    pass

try:
    for db in data["Extra"]["Databases"]:
        try:
            db["Name"] = list2[assoc(db["Name"], len(list2))]
        except Exception as ex:
            pass
        try:
            db["UniqueName"] = list2[assoc(db["UniqueName"], len(list2))]
        except Exception as ex:
            pass
        try:
            for patch in db["Patches"]:
                patch["Database"] = list2[assoc(patch["Database"], len(list2))]
        except Exception as ex:
            pass
        try:
            for ts in db["Tablespaces"]:
                ts["Database"] = list2[assoc(ts["Database"], len(list2))]
                ts["Name"] = list2[assoc(ts["Name"], len(list2))]
        except Exception as ex:
            pass
        try:
            for sc in db["Schemas"]:
                sc["Database"] = list2[assoc(sc["Database"], len(list2))]
                sc["User"] = list2[assoc(sc["User"], len(list2))]
        except Exception as ex:
            pass
        try:
            for sa in db["SegmentAdvisors"]:
                sa["SegmentOwner"] = list3[assoc(sa["SegmentOwner"], len(list3))]
                sa["SegmentName"] = list2[assoc(sa["SegmentName"], len(list2))]
                sa["Recommendation"] = md5(sa["Recommendation"].encode("ascii")).hexdigest()
        except Exception as ex:
            pass
except Exception as ex:
    pass


json.dump(data, sys.stdout)