#!/usr/bin/env python3
import json
import sys
import io
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
    digest = reduce(lambda x, y: (x+y) % mod, map(lambda x: x[0]*256 + ord(
        x[1]), enumerate(list(md5(p.encode("ascii")).hexdigest()))))
    return digest


def assoc_list(p: str, l: list):
    return l[assoc(p, len(l))] + "-" + md5(p.encode("ascii")).hexdigest()


data = json.load(sys.stdin)


try:
    data["hostname"] = assoc_list(data["hostname"], list0)
except Exception as ex:
    pass
try:
    data["info"]["hostname"] = assoc_list(data["info"]["hostname"], list0)
except Exception as ex:
    pass
try:
    data["databases"] = " ".join(
        map(lambda x: assoc_list(x, list2),  data["databases"].split(" ")))
except Exception as ex:
    pass
try:
    data["schemas"] = " ".join(
        map(lambda x: assoc_list(x, list2),  data["schemas"].split(" ")))
except Exception as ex:
    pass
try:
    for fs in data["extra"]["filesystems"]:
        try:
            fs["filesystem"] = assoc_list(fs["filesystem"], list4)
        except Exception:
            pass

        try:
            fs["mountedOn"] = assoc_list(fs["mountedOn"], list4)
        except Exception:
            pass
except Exception:
    pass

try:
    for db in data["extra"]["databases"]:
        try:
            db["name"] = assoc_list(db["name"], list2)
        except Exception as ex:
            pass
        try:
            db["uniqueName"] = assoc_list(db["uniqueName"], list2)
        except Exception as ex:
            pass
        try:
            for patch in db["patches"]:
                patch["database"] = assoc_list(patch["database"], list2)
        except Exception as ex:
            pass
        try:
            for ts in db["tablespaces"]:
                ts["database"] = assoc_list(ts["database"], list2)
                ts["name"] = assoc_list(ts["name"], list2)
        except Exception as ex:
            pass
        try:
            for sc in db["schemas"]:
                sc["database"] = assoc_list(sc["database"], list2)
                sc["user"] = assoc_list(sc["user"], list2)
        except Exception as ex:
            pass
        try:
            for sa in db["segmentAdvisors"]:
                sa["segmentOwner"] = assoc_list(sa["segmentOwner"], list3)
                sa["segmentName"] = assoc_list(sa["segmentName"], list2)
                sa["recommendation"] = md5(
                    sa["recommendation"].encode("ascii")).hexdigest()
        except Exception as ex:
            pass
except Exception as ex:
    pass


try:
    for db in data["extra"]["databases"]:
        try:
            db["name"] = assoc_list(db["name"], list2)
        except Exception as ex:
            pass
        try:
            db["uniqueName"] = assoc_list(db["uniqueName"], list2)
        except Exception as ex:
            pass
        try:
            for patch in db["patches"]:
                patch["database"] = assoc_list(patch["database"], list2)
        except Exception as ex:
            pass
        try:
            for ts in db["tablespaces"]:
                ts["database"] = assoc_list(ts["database"], list2)
                ts["name"] = assoc_list(ts["name"], list2)
        except Exception as ex:
            pass
        try:
            for sc in db["schemas"]:
                sc["database"] = assoc_list(sc["database"], list2)
                sc["user"] = assoc_list(sc["user"], list2)
        except Exception as ex:
            pass
        try:
            for sa in db["segmentAdvisors"]:
                sa["segmentOwner"] = assoc_list(sa["segmentOwner"], list3)
                sa["segmentName"] = assoc_list(sa["segmentName"], list2)
                sa["recommendation"] = md5(
                    sa["recommendation"].encode("ascii")).hexdigest()
        except Exception as ex:
            pass
except Exception as ex:
    pass

try:
    for cl in data["clusters"]:
        try:
            if cl["name"] != "not_in_cluster":
                cl["name"] = assoc_list(cl["name"], list1)
            for vm in cl["vms"]:
                try:
                    vm["name"] = assoc_list(vm["name"], list0)
                    vm["hostname"] = assoc_list(vm["hostname"], list0)
                    vm["virtualizationNode"] = assoc_list(
                        vm["virtualizationNode"], list0)
                    vm["clusterName"] = assoc_list(vm["clusterName"], list1)
                except Exception:
                    pass
            if cl["fetchEndpoint"] == "":
                fetchEndpoint = "http://" + cl["name"] + ".test"
                cl["fetchEndpoint"] = fetchEndpoint[:64]

        except Exception:
            pass
except Exception:
    pass

try:
    for dev in data["features"]["oracle"]["exadata"]["components"]:
        try:
            dev["hostname"] = assoc_list(dev["hostname"], list0)
        except Exception:
            pass
        try:
            for cd in dev["cellDisks"]:
                try:
                    cd["name"] = assoc_list(cd["name"], list0)
                except Exception:
                    pass
        except Exception:
            pass
except Exception:
    pass

json.dump(data, sys.stdout)
