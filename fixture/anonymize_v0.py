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
list4 = []
with open("fixture/list4.txt") as fp:
    for line in fp:
        list4.append(line.strip())

# A injective and suriective function...
def assoc(p: str, mod: int):
    digest = reduce(lambda x, y: (x+y) % mod, map(lambda x: x[0]*256 + ord(x[1]), enumerate(list(md5(p.encode("ascii")).hexdigest()))))
    return digest

def assoc_list(p: str, l: list):
    return l[assoc(p, len(l))] + "-" + md5(p.encode("ascii")).hexdigest()

data = json.load(sys.stdin)

try:
    data["Hostname"] = assoc_list(data["Hostname"], list0)
except Exception as ex:
    pass
try:
    data["Info"]["Hostname"] = assoc_list(data["Info"]["Hostname"], list0)
except Exception as ex:
    pass
try:
    data["Databases"] = " ".join(map(lambda x: assoc_list(x, list2),  data["Databases"].split(" ")))
except Exception as ex:
    pass
try:
    data["Schemas"] = " ".join(map(lambda x: assoc_list(x, list2),  data["Schemas"].split(" ")))
except Exception as ex:
    pass
try:
    for fs in data["Extra"]["Filesystems"]:
        try:
            fs["Filesystem"] = assoc_list(fs["Filesystem"], list4)
        except Exception:
            pass

        try:
            fs["MountedOn"] = assoc_list(fs["MountedOn"], list4)
        except Exception:
            pass
except Exception:
    pass

try:
    for db in data["Extra"]["Databases"]:
        try:
            db["Name"] = assoc_list(db["Name"], list2)
        except Exception as ex:
            pass
        try:
            db["UniqueName"] = assoc_list(db["UniqueName"], list2)
        except Exception as ex:
            pass
        try:
            for patch in db["Patches"]:
                patch["Database"] = assoc_list(patch["Database"], list2)
        except Exception as ex:
            pass
        try:
            for ts in db["Tablespaces"]:
                ts["Database"] = assoc_list(ts["Database"], list2)
                ts["Name"] = assoc_list(ts["Name"], list2)
        except Exception as ex:
            pass
        try:
            for sc in db["Schemas"]:
                sc["Database"] = assoc_list(sc["Database"], list2)
                sc["User"] = assoc_list(sc["User"], list2)
        except Exception as ex:
            pass
        try:
            for sa in db["SegmentAdvisors"]:
                sa["SegmentOwner"] = assoc_list(sa["SegmentOwner"], list3)
                sa["SegmentName"] = assoc_list(sa["SegmentName"], list2)
                sa["Recommendation"] = md5(sa["Recommendation"].encode("ascii")).hexdigest()
        except Exception as ex:
            pass
except Exception as ex:
    pass


try:
    for db in data["Extra"]["Databases"]:
        try:
            db["Name"] = assoc_list(db["Name"], list2)
        except Exception as ex:
            pass
        try:
            db["UniqueName"] = assoc_list(db["UniqueName"], list2)
        except Exception as ex:
            pass
        try:
            for patch in db["Patches"]:
                patch["Database"] = assoc_list(patch["Database"], list2)
        except Exception as ex:
            pass
        try:
            for ts in db["Tablespaces"]:
                ts["Database"] = assoc_list(ts["Database"], list2)
                ts["Name"] = assoc_list(ts["Name"], list2)
        except Exception as ex:
            pass
        try:
            for sc in db["Schemas"]:
                sc["Database"] = assoc_list(sc["Database"], list2)
                sc["User"] = assoc_list(sc["User"], list2)
        except Exception as ex:
            pass
        try:
            for sa in db["SegmentAdvisors"]:
                sa["SegmentOwner"] = assoc_list(sa["SegmentOwner"], list3)
                sa["SegmentName"] = assoc_list(sa["SegmentName"], list2)
                sa["Recommendation"] = md5(sa["Recommendation"].encode("ascii")).hexdigest()
        except Exception as ex:
            pass
except Exception as ex:
    pass

try:
    for cl in data["Extra"]["Clusters"]:
        try:
            if cl["Name"] != "not_in_cluster":
                cl["Name"] = assoc_list(cl["Name"], list1)
            for vm in cl["VMs"]:
                try:
                    vm["Name"] = assoc_list(vm["Name"], list0)
                    vm["Hostname"] = assoc_list(vm["Hostname"], list0)
                    vm["VirtualizationNode"] = assoc_list(vm["VirtualizationNode"], list0)
                    vm["ClusterName"] = assoc_list(vm["ClusterName"], list1)
                except Exception:
                    pass
        except Exception:
            pass
except Exception:
    pass

try:
    for dev in data["Extra"]["Exadata"]["Devices"]:
        try:
            dev["Hostname"] = assoc_list(dev["Hostname"], list0)
        except Exception:
            pass
        try:
            for cd in dev["CellDisks"]:
                try:
                    cd["Name"] = assoc_list(cd["Name"], list0)
                except Exception:
                    pass
        except Exception:
            pass
except Exception:
    pass

json.dump(data, sys.stdout)