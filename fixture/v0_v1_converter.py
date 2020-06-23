#!/usr/bin/env python3
import json, sys, io
from hashlib import md5
from functools import reduce 
from jsonschema import Draft7Validator
from datetime import datetime

# WARN: this is a example code for result oriented programming

data = json.load(sys.stdin)
schema = dict()
with open('hostdata.v1.strict.schema.json') as json_file:
    schema = json.load(json_file)
v = Draft7Validator(schema)

data["AgentVersion"] = data["Version"]
del data["Version"]
del data["HostType"]
data["Tags"] = []
del data["Databases"]
del data["Schemas"]
data["SchemaVersion"] = 1
del data["HostDataSchemaVersion"]

data["Info"]["CPUFrequency"] = data["Info"]["CPUModel"][data["Info"]["CPUModel"].index("@ ") + 2:]
data["Info"]["CPUSockets"] = data["Info"]["Socket"]
del data["Info"]["Socket"]
del data["Info"]["Environment"]
del data["Info"]["Location"]
if data["Info"]["CPUSockets"] > data["Info"]["CPUCores"]:
    data["Info"]["CoresPerSocket"] = 1
else:
    data["Info"]["CoresPerSocket"] = int(data["Info"]["CPUCores"] / data["Info"]["CPUSockets"])
if data["Info"]["CPUThreads"] == 1:
    data["Info"]["ThreadsPerCore"] = 1
else:
    data["Info"]["ThreadsPerCore"] = 2
data["Info"]["KernelVersion"] = data["Info"]["Kernel"]
data["Info"]["KernelVersion"] = data["Info"]["Kernel"]
if "Red Hat Enterprise Linux Server" in data["Info"]["OS"]:
    data["Info"]["OSVersion"] = data["Info"]["OS"].split(" ")[6]
    data["Info"]["OS"] = "Red Hat Enterprise Linux"
    data["Info"]["Kernel"] = data["Info"]["Kernel"]
elif "Ubuntu" in data["Info"]["OS"]:
    data["Info"]["OSVersion"] = data["Info"]["OS"].split(" ")[-1]
    data["Info"]["OS"] = "Ubuntu"
    data["Info"]["Kernel"] = data["Info"]["Kernel"]
else:
    print("I don't know which operating system is ", data["Info"]["OS"])
data["ClusterMembershipStatus"] = {
    "OracleClusterware": data["Info"]["OracleCluster"],
    "VeritasClusterServer": data["Info"]["VeritasCluster"],
    "SunCluster": data["Info"]["SunCluster"],
    "HACMP": data["Info"]["AixCluster"],
} 
del data["Info"]["OracleCluster"]
del data["Info"]["VeritasCluster"]
del data["Info"]["SunCluster"]
del data["Info"]["AixCluster"]
if data["Info"]["Virtual"]:
    data["Info"]["HardwareAbstraction"] = "VIRT"
else:
    data["Info"]["HardwareAbstraction"] = "PH"
del data["Info"]["Virtual"]
data["Info"]["HardwareAbstractionTechnology"] = data["Info"]["Type"] 
del data["Info"]["Type"]

data["Filesystems"] = data["Extra"]["Filesystems"]
del data["Extra"]["Filesystems"]
for f in data["Filesystems"]:
    if "T" in f["Size"]:
        f["Size"] = int(float(f["Size"][:-1])*1024*1024*1024*1024)
    elif "G" in f["Size"]:
        f["Size"] = int(float(f["Size"][:-1])*1024*1024*1024)
    elif "M" in f["Size"]:
        f["Size"] = int(float(f["Size"][:-1])*1024*1024)
    elif "K" in f["Size"]:
        f["Size"] = int(float(f["Size"][:-1])*1024)
    elif f["Size"] == "0":
        f["Size"] = 0
    f["AvailableSpace"] = f["Available"]
    del f["Available"]
    if "T" in f["AvailableSpace"]:
        f["AvailableSpace"] = int(float(f["AvailableSpace"][:-1])*1024*1024*1024*1024)
    elif "G" in f["AvailableSpace"]:
        f["AvailableSpace"] = int(float(f["AvailableSpace"][:-1])*1024*1024*1024)
    elif "M" in f["AvailableSpace"]:
        f["AvailableSpace"] = int(float(f["AvailableSpace"][:-1])*1024*1024)
    elif "K" in f["AvailableSpace"]:
        f["AvailableSpace"] = int(float(f["AvailableSpace"][:-1])*1024)
    elif f["AvailableSpace"] == "0":
        f["AvailableSpace"] = 0
    f["UsedSpace"] = f["Used"]
    del f["Used"]
    if "T" in f["UsedSpace"]:
        f["UsedSpace"] = int(float(f["UsedSpace"][:-1])*1024*1024*1024*1024)
    elif "G" in f["UsedSpace"]:
        f["UsedSpace"] = int(float(f["UsedSpace"][:-1])*1024*1024*1024)
    elif "M" in f["UsedSpace"]:
        f["UsedSpace"] = int(float(f["UsedSpace"][:-1])*1024*1024)
    elif "K" in f["UsedSpace"]:
        f["UsedSpace"] = int(float(f["UsedSpace"][:-1])*1024)
    elif f["UsedSpace"] == "0":
        f["UsedSpace"] = 0
    del f["UsedPerc"]
    f["Type"] = f["FsType"]
    del f["FsType"]


data["Features"] = {}
if "Databases" in data["Extra"] and data["Extra"]["Databases"] != None and len(data["Extra"]["Databases"]) > 0 :
    if not "Oracle" in data["Features"]:
        data["Features"]["Oracle"] = {}
    data["Features"]["Oracle"]["Database"] = {
        "Databases": data["Extra"]["Databases"]
    }
    del data["Extra"]["Databases"]
    for db in data["Features"]["Oracle"]["Database"]["Databases"]:
        db["InstanceNumber"] = int(db["InstanceNumber"])
        if db["Archivelog"] == "NOARCHIVELOG":
            db["Archivelog"] = False
        elif db["Archivelog"] == "ARCHIVELOG":
            db["Archivelog"] = True
        else:
            print("I don't know the value of archivelog ", db["Archivelog"])
        db["BlockSize"] = int(db["BlockSize"])
        db["CPUCount"] = int(db["CPUCount"])
        db["SGATarget"] = float(db["SGATarget"])
        db["PGATarget"] = float(db["PGATarget"])
        db["MemoryTarget"] = float(db["MemoryTarget"])
        db["SGAMaxSize"] = float(db["SGAMaxSize"])
        db["SegmentsSize"] = float(db["SegmentsSize"])
        db["DatafileSize"] = float(db["Used"])
        del db["Used"]
        db["Allocated"] = float(db["Allocated"])
        db["Elapsed"] = float(db["Elapsed"])
        db["DBTime"] = float(db["DBTime"])
        db["DailyCPUUsage"] = float(db["DailyCPUUsage"])
        if db["Work"] == "N/A":
            db["Work"] = None
        else:
            db["Work"] = float(db["Work"])
        for tb in db["Tablespaces"]:
            del tb["Database"]
            tb["MaxSize"] = float(tb["MaxSize"])
            tb["Total"] = float(tb["Total"])
            tb["Used"] = float(tb["Used"])
            tb["UsedPerc"] = float(tb["UsedPerc"])
        for sc in db["Schemas"]:
            del sc["Database"]
        db["PSUs"] = db["LastPSUs"]
        del db["LastPSUs"]
        if "Features2" in db:
            db["FeatureUsageStats"] = db["Features2"]
            del db["Features2"]
        else:
            db["FeatureUsageStats"] = []
        for fus in db["FeatureUsageStats"]:
            fus["FirstUsageDate"] = datetime.strptime(fus["FirstUsageDate"], "%Y-%m-%d %H:%M:%S").strftime("%Y-%m-%dT%H:%M:%SZ")
            fus["LastUsageDate"] = datetime.strptime(fus["LastUsageDate"], "%Y-%m-%d %H:%M:%S").strftime("%Y-%m-%dT%H:%M:%SZ")
        db["IsCDB"] = False
        db["PDBs"] = []
        db["Services"] = []
        for pt in db["Patches"]:
            del pt["Database"]
            if pt["PatchID"] == "":
                pt["PatchID"] = -1
            else:
                pt["PatchID"] = int(pt["PatchID"])
            pt["Date"] = datetime.strptime(pt["Date"], "%d-%b-%Y").strftime("%F")
        for sa in db["SegmentAdvisors"]:
            if sa["Reclaimable"] == "<1":
                sa["Reclaimable"] = 0.5
            else:
                sa["Reclaimable"] = float(sa["Reclaimable"])
        for addm in db["ADDMs"]:
            addm["Benefit"] = float(addm["Benefit"])
        for ba in db["Backups"]:
            ba["WeekDays"] = ba["WeekDays"].split(",")
            ba["AvgBckSize"] = float(ba["AvgBckSize"])
            

if "Clusters" in data["Extra"] and data["Extra"]["Clusters"] != None and len(data["Extra"]["Clusters"]) > 0 :
    data["Clusters"] = data["Extra"]["Clusters"]
    del data["Extra"]["Clusters"]
    for cl in data["Clusters"]:
        cl["FetchEndpoint"] = "???"
        for vm in cl["VMs"]:
            del vm["ClusterName"]
            vm["VirtualizationNode"] = vm["PhysicalHost"]
            del vm["PhysicalHost"]
else:
    data["Clusters"] = None

if "Exadata" in data["Extra"] and data["Extra"]["Exadata"] != None:
    if not "Oracle" in data["Features"]:
        data["Features"]["Oracle"] = {}
    data["Features"]["Oracle"]["Exadata"] = data["Extra"]["Exadata"]
    del data["Extra"]["Exadata"]
    data["Features"]["Oracle"]["Exadata"]["Components"] = data["Features"]["Oracle"]["Exadata"]["Devices"] 
    del data["Features"]["Oracle"]["Exadata"]["Devices"]
    for com in data["Features"]["Oracle"]["Exadata"]["Components"]:
        com["SwVersion"] = com["ExaSwVersion"]
        del com["ExaSwVersion"]
        com["SwReleaseDate"] = com["SwVersion"].split(".")[-1]
        com["SwReleaseDate"] = "20" + com["SwReleaseDate"][0:2] + "-" + com["SwReleaseDate"][2:4] + "-" + com["SwReleaseDate"][4:6]
        if com["CPUEnabled"] != "-":
            com["RunningCPUCount"] = int(com["CPUEnabled"].split("/")[0])
            com["TotalCPUCount"] = int(com["CPUEnabled"].split("/")[1])
        else:   
            com["RunningCPUCount"] = None
            com["TotalCPUCount"] = None
        del com["CPUEnabled"]
        if com["Memory"] != "-":
            com["Memory"] = int(com["Memory"][:-2])
        else:
            com["Memory"] = None
        if com["TempActual"] != "-":
            com["TempActual"] = float(com["TempActual"])
        else:
            com["TempActual"] = None
        if com["PowerCount"] != "-":
            com["RunningPowerSupply"] = int(com["PowerCount"].split("/")[0])
            com["TotalPowerSupply"] = int(com["PowerCount"].split("/")[1])
        else:   
            com["RunningPowerSupply"] = None
            com["TotalPowerSupply"] = None
        del com["PowerCount"]
        if com["FanCount"] != "-":
            com["RunningFanCount"] = int(com["FanCount"].split("/")[0])
            com["TotalFanCount"] = int(com["FanCount"].split("/")[1])
        else:   
            com["RunningFanCount"] = None
            com["TotalFanCount"] = None
        del com["FanCount"]
        com["CellsrvServiceStatus"] = com["CellsrvService"]
        del com["CellsrvService"]
        com["MsServiceStatus"] = com["MsService"]
        del com["MsService"]
        com["RsServiceStatus"] = com["RsService"]
        del com["RsService"]
        if com["FlashcacheMode"] == "-":
            com["FlashcacheMode"] = None
        if com["Status"] == "-":
            com["Status"] = None
        if com["CellsrvServiceStatus"] == "-":
            com["CellsrvServiceStatus"] = None
        if com["MsServiceStatus"] == "-":
            com["MsServiceStatus"] = None
        if com["RsServiceStatus"] == "-":
            com["RsServiceStatus"] = None
        if com["PowerStatus"] == "-":
            com["PowerStatus"] = None
        if com["FanStatus"] == "-":
            com["FanStatus"] = None
        if com["TempStatus"] == "-":
            com["TempStatus"] = None
        if com["CellDisks"] != None:
            for cd in com["CellDisks"]:        
                cd["ErrCount"] = int(cd["ErrCount"])
                cd["UsedPerc"] = int(cd["UsedPerc"])

del data["Extra"]

for error in sorted(v.iter_errors(data), key=str):
    print(error.message)

json.dump(data, sys.stdout)
print()