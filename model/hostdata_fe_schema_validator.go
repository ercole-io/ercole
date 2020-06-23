// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package model

var FrontendHostdataSchemaValidator string = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id": "ercole-hostdata",
    "type": "object",
    "required": [
        "Hostname",
        "Location",
        "Environment",
        "Tags",
        "AgentVersion",
        "SchemaVersion",
        "Info",
        "ClusterMembershipStatus",
        "Features",
        "Filesystems"
    ],
    "properties": {
        "Hostname": {
            "type": "string",
            "minLength": 1,
            "maxLength": 253,
            "format": "idn-hostname"
        },
        "Location": {
            "type": "string",
            "minLength": 1,
            "maxLength": 64,
            "pattern": "^[a-zA-Z0-9-]+$"
        },
        "Environment": {
            "type": "string",
            "minLength": 1,
            "maxLength": 16,
            "pattern": "^[A-Z0-9]+$"
        },
        "Tags": {
            "type": "array",
            "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 64,
                "pattern": "^[a-zA-Z0-9-]+$"
            },
            "uniqueItems": true
        },
        "AgentVersion": {
            "type": "string",
            "minLength": 1,
            "maxLength": 64,
            "pattern": "^(([0-9]+([.][0-9]+)*)|(git-[0-9a-f]+)|(latest))$"
        },
        "SchemaVersion": {
            "type": "integer",
            "const": 1
        },
        "Info": {
            "type": "object",
            "required": [
                "Hostname",
                "CPUModel",
                "CPUFrequency",
                "CPUSockets",
                "CPUCores",
                "CPUThreads",
                "ThreadsPerCore",
                "CoresPerSocket",
                "HardwareAbstraction",
                "HardwareAbstractionTechnology",
                "Kernel",
                "KernelVersion",
                "OS",
                "OSVersion",
                "MemoryTotal",
                "SwapTotal"
            ],
            "properties": {
                "Hostname": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 253,
                    "format": "idn-hostname"
                },
                "CPUModel": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64,
                    "pattern": "^[^\n]+$"
                },
                "CPUFrequency": {
                    "type": "string",
                    "minLength": 3,
                    "maxLength": 16,
                    "pattern": "^[0-9]+([.][0-9]+)?[ ]*(?i)(GHz|MHz)$"
                },
                "CPUSockets": {
                    "type": "integer",
                    "minimum": 0
                },
                "CPUCores": {
                    "type": "integer",
                    "minimum": 1
                },
                "CPUThreads": {
                    "type": "integer",
                    "minimum": 1
                },
                "ThreadsPerCore": {
                    "type": "integer",
                    "minimum": 1
                },
                "CoresPerSocket": {
                    "type": "integer",
                    "minimum": 1
                },
                "HardwareAbstraction": {
                    "type": "string",
                    "enum": ["PH", "VIRT"]
                },
                "HardwareAbstractionTechnology": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 32,
                    "pattern": "^[A-Z0-9]+$"
                },
                "Kernel": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "KernelVersion": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "OS": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "OSVersion": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "MemoryTotal": {
                    "type": "number",
                    "minimum": 0
                },
                "SwapTotal": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "ClusterMembershipStatus": {
            "type": "object",
            "properties": {
                "OracleClusterware": {
                    "type": "boolean"
                },
                "VeritasClusterServer": {
                    "type": "boolean"
                },
                "SunCluster": {
                    "type": "boolean"
                },
                "HACMP": {
                    "type": "boolean"
                }
            }
        },
        "Features": {
            "type": "object",
            "properties": {
                "Oracle": {
                    "anyOf": [{
                            "type": "null"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "Database": {
                                    "anyOf": [{
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "Databases"
                                            ],
                                            "properties": {
                                                "Databases": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "InstanceNumber",
                                                            "Name",
                                                            "UniqueName",
                                                            "Status",
                                                            "IsCDB",
                                                            "Version",
                                                            "Platform",
                                                            "Archivelog",
                                                            "Charset",
                                                            "NCharset",
                                                            "BlockSize",
                                                            "CPUCount",
                                                            "SGATarget",
                                                            "PGATarget",
                                                            "MemoryTarget",
                                                            "SGAMaxSize",
                                                            "SegmentsSize",
                                                            "DatafileSize",
                                                            "Allocated",
                                                            "Elapsed",
                                                            "DBTime",
                                                            "DailyCPUUsage",
                                                            "Work",
                                                            "ASM",
                                                            "Dataguard",
                                                            "Patches",
                                                            "Tablespaces",
                                                            "Schemas",
                                                            "Licenses",
                                                            "ADDMs",
                                                            "SegmentAdvisors",
                                                            "PSUs",
                                                            "Backups",
                                                            "FeatureUsageStats",
                                                            "PDBs",
                                                            "Services"
                                                        ],
                                                        "properties": {
                                                            "InstanceNumber": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "Name": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "UniqueName": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "Status": {
                                                                "type": "string",
                                                                "enum": [
                                                                    "OPEN",
                                                                    "MOUNTED"
                                                                ]
                                                            },
                                                            "IsCDB": {
                                                                "type": "boolean"
                                                            },
                                                            "Version": {
                                                                "type": "string",
                                                                "minLength": 8,
                                                                "maxLength": 64
                                                            },
                                                            "Platform": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "Archivelog": {
                                                                "type": "boolean"
                                                            },
                                                            "Charset": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "NCharset": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "BlockSize": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "CPUCount": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "SGATarget": {
                                                                "type": "number"
                                                            },
                                                            "PGATarget": {
                                                                "type": "number"
                                                            },
                                                            "MemoryTarget": {
                                                                "type": "number"
                                                            },
                                                            "SGAMaxSize": {
                                                                "type": "number"
                                                            },
                                                            "SegmentsSize": {
                                                                "type": "number"
                                                            },
                                                            "DatafileSize": {
                                                                "type": "number"
                                                            },
                                                            "Allocated": {
                                                                "type": "number"
                                                            },
                                                            "Elapsed": {
                                                                "anyOf": [{
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "DBTime": {
                                                                "anyOf": [{
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "DailyCPUUsage": {
                                                                "anyOf": [{
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "Work": {
                                                                "anyOf": [{
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "ASM": {
                                                                "type": "boolean"
                                                            },
                                                            "Dataguard": {
                                                                "type": "boolean"
                                                            },
                                                            "Patches": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Version",
                                                                        "PatchID",
                                                                        "Action",
                                                                        "Description",
                                                                        "Date"
                                                                    ],
                                                                    "properties": {
                                                                        "Version": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 16
                                                                        },
                                                                        "PatchID": {
                                                                            "type": "integer"
                                                                        },
                                                                        "Action": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 128
                                                                        },
                                                                        "Description": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 256
                                                                        },
                                                                        "Date": {
                                                                            "type": "string",
                                                                            "format": "date"
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "Tablespaces": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Name",
                                                                        "MaxSize",
                                                                        "Total",
                                                                        "Used",
                                                                        "UsedPerc",
                                                                        "Status"
                                                                    ],
                                                                    "properties": {
                                                                        "Name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "MaxSize": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "Total": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "Used": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "UsedPerc": {
                                                                            "type": "number",
                                                                            "minimum": 0,
                                                                            "maximum": 100
                                                                        },
                                                                        "Status": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "ONLINE",
                                                                                "READ ONLY",
                                                                                "OFFLINE"
                                                                            ]
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "Schemas": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Indexes",
                                                                        "LOB",
                                                                        "Tables",
                                                                        "Total",
                                                                        "User"
                                                                    ],
                                                                    "properties": {
                                                                        "Indexes": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "LOB": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "Tables": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "Total": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "User": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "Licenses": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Name",
                                                                        "Count"
                                                                    ],
                                                                    "properties": {
                                                                        "Name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "Count": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "ADDMs": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Finding",
                                                                        "Recommendation",
                                                                        "Action",
                                                                        "Benefit"
                                                                    ],
                                                                    "properties": {
                                                                        "Finding": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 256
                                                                        },
                                                                        "Recommendation": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        },
                                                                        "Action": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 256
                                                                        },
                                                                        "Benefit": {
                                                                            "type": "number",
                                                                            "minimum": 0,
                                                                            "maximum": 100
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "SegmentAdvisors": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "SegmentOwner",
                                                                        "SegmentName",
                                                                        "SegmentType",
                                                                        "PartitionName",
                                                                        "Reclaimable",
                                                                        "Recommendation"
                                                                    ],
                                                                    "properties": {
                                                                        "SegmentOwner": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "SegmentName": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "SegmentType": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        },
                                                                        "PartitionName": {
                                                                            "type": "string",
                                                                            "maxLength": 32
                                                                        },
                                                                        "Reclaimable": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "Recommendation": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 256
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "PSUs": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Date",
                                                                        "Description"
                                                                    ],
                                                                    "properties": {
                                                                        "Date": {
                                                                            "type": "string",
                                                                            "format": "date"
                                                                        },
                                                                        "Description": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 128
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "Backups": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "BackupType",
                                                                        "Hour",
                                                                        "WeekDays",
                                                                        "AvgBckSize",
                                                                        "Retention"
                                                                    ],
                                                                    "properties": {
                                                                        "BackupType": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "Archivelog",
                                                                                "Full",
                                                                                "Level0",
                                                                                "Level1"
                                                                            ]
                                                                        },
                                                                        "Hour": {
                                                                            "type": "string",
                                                                            "minLength": 5,
                                                                            "maxLength": 5,
                                                                            "pattern": "^[0-9]{2}:[0-9]{2}$"
                                                                        },
                                                                        "WeekDays": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "string",
                                                                                "enum": [
                                                                                    "Monday",
                                                                                    "Tuesday",
                                                                                    "Wednesday",
                                                                                    "Thursday",
                                                                                    "Friday",
                                                                                    "Saturday",
                                                                                    "Sunday"
                                                                                ]
                                                                            },
                                                                            "uniqueItems": true
                                                                        },
                                                                        "AvgBckSize": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "Retention": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 16
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "FeatureUsageStats": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Product",
                                                                        "Feature",
                                                                        "DetectedUsages",
                                                                        "CurrentlyUsed",
                                                                        "FirstUsageDate",
                                                                        "LastUsageDate",
                                                                        "ExtraFeatureInfo"
                                                                    ],
                                                                    "properties": {
                                                                        "Product": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        },
                                                                        "Feature": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        },
                                                                        "DetectedUsages": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "CurrentlyUsed": {
                                                                            "type": "boolean"
                                                                        },
                                                                        "FirstUsageDate": {
                                                                            "type": "string",
                                                                            "format": "date-time"
                                                                        },
                                                                        "LastUsageDate": {
                                                                            "type": "string",
                                                                            "format": "date-time"
                                                                        },
                                                                        "ExtraFeatureInfo": {
                                                                            "type": "string",
                                                                            "maxLength": 64
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "Services": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Name"
                                                                    ],
                                                                    "properties": {
                                                                        "Name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "PDBs": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Name",
                                                                        "Status",
                                                                        "Tablespaces",
                                                                        "Schemas",
                                                                        "Services"
                                                                    ],
                                                                    "properties": {
                                                                        "Name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        },
                                                                        "Status": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "OPEN",
                                                                                "READ WRITE",
                                                                                "MOUNTED"
                                                                            ]
                                                                        },
                                                                        "Tablespaces": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "Name",
                                                                                    "MaxSize",
                                                                                    "Total",
                                                                                    "Used",
                                                                                    "UsedPerc",
                                                                                    "Status"
                                                                                ],
                                                                                "properties": {
                                                                                    "Name": {
                                                                                        "type": "string",
                                                                                        "minLength": 1,
                                                                                        "maxLength": 32
                                                                                    },
                                                                                    "MaxSize": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "Total": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "Used": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "UsedPerc": {
                                                                                        "type": "number",
                                                                                        "minimum": 0,
                                                                                        "maximum": 100
                                                                                    },
                                                                                    "Status": {
                                                                                        "type": "string",
                                                                                        "enum": [
                                                                                            "ONLINE",
                                                                                            "READ ONLY",
                                                                                            "OFFLINE"
                                                                                        ]
                                                                                    }
                                                                                }
                                                                            }
                                                                        },
                                                                        "Schemas": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "Indexes",
                                                                                    "LOB",
                                                                                    "Tables",
                                                                                    "Total",
                                                                                    "User"
                                                                                ],
                                                                                "properties": {
                                                                                    "Indexes": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "LOB": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "Tables": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "Total": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "User": {
                                                                                        "type": "string",
                                                                                        "minLength": 1,
                                                                                        "maxLength": 32
                                                                                    }
                                                                                }
                                                                            }
                                                                        },
                                                                        "Services": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "Name"
                                                                                ],
                                                                                "properties": {
                                                                                    "Name": {
                                                                                        "type": "string",
                                                                                        "minLength": 1,
                                                                                        "maxLength": 32
                                                                                    }
                                                                                }
                                                                            }
                                                                        }
                                                                    }
                                                                }
                                                            }
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    ]
                                },
                                "Exadata": {
                                    "anyOf": [{
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "Components"
                                            ],
                                            "properties": {
                                                "Components": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "Hostname",
                                                            "ServerType",
                                                            "Model",
                                                            "SwVersion",
                                                            "SwReleaseDate",
                                                            "RunningCPUCount",
                                                            "TotalCPUCount",
                                                            "Memory",
                                                            "Status",
                                                            "RunningPowerSupply",
                                                            "TotalPowerSupply",
                                                            "PowerStatus",
                                                            "RunningFanCount",
                                                            "TotalFanCount",
                                                            "FanStatus",
                                                            "TempActual",
                                                            "TempStatus",
                                                            "CellsrvServiceStatus",
                                                            "MsServiceStatus",
                                                            "RsServiceStatus",
                                                            "FlashcacheMode",
                                                            "CellDisks"
                                                        ],
                                                        "properties": {
                                                            "Hostname": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 253,
                                                                "format": "idn-hostname"
                                                            },
                                                            "ServerType": {
                                                                "type": "string",
                                                                "enum": [
                                                                    "DBServer",
                                                                    "IBSwitch",
                                                                    "StorageServer"
                                                                ]
                                                            },
                                                            "Model": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "SwVersion": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "SwReleaseDate": {
                                                                "type": "string",
                                                                "format": "date"
                                                            },
                                                            "RunningCPUCount": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "TotalCPUCount": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "Memory": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "$comment": "Memory in GB",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "Status": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string",
                                                                        "enum": [
                                                                            "online",
                                                                            "offline"
                                                                        ]
                                                                    }
                                                                ]
                                                            },
                                                            "RunningPowerSupply": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "TotalPowerSupply": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "PowerStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "RunningFanCount": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "TotalFanCount": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "FanStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "TempActual": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "TempStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "CellsrvServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "MsServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "RsServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "FlashcacheMode": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string",
                                                                        "enum": [
                                                                            "WriteBack",
                                                                            "WriteThrough"
                                                                        ]
                                                                    }
                                                                ]
                                                            },
                                                            "CellDisks": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "array",
                                                                        "items": {
                                                                            "type": "object",
                                                                            "required": [
                                                                                "ErrCount",
                                                                                "Name",
                                                                                "Status",
                                                                                "UsedPerc"
                                                                            ],
                                                                            "properties": {
                                                                                "ErrCount": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0
                                                                                },
                                                                                "Name": {
                                                                                    "type": "string",
                                                                                    "minLength": 1,
                                                                                    "maxLength": 64
                                                                                },
                                                                                "Status": {
                                                                                    "type": "string"
                                                                                },
                                                                                "UsedPerc": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0,
                                                                                    "maximum": 100
                                                                                }
                                                                            }
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        { "type": "null" }
                    ]
                },
                "Postgresql": {
                    "anyOf": [{
                            "type": "null"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "Postgresql": {
                                    "anyOf": [{
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "WorkMem",
                                                "ArchiveMode",
                                                "ArchivePath",
                                                "MinWalSize",
                                                "MaxWalSize",
                                                "MaxConnections",
                                                "CheckpointCompletionTarget",
                                                "DefaultStatisticsTarget",
                                                "RandomPageCost",
                                                "MaintenanceWorkMem",
                                                "SharedBuffers",
                                                "EffectiveCacheSize",
                                                "EffectiveIOConcurrency",
                                                "MaxWorkerProcesses",
                                                "MaxParallelWorkers",
                                                "Databases"
                                            ],
                                            "properties": {
                                                "WorkMem": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "ArchiveMode": {
                                                    "type": "boolean"
                                                },
                                                "ArchivePath": {
                                                    "type": "string",
                                                    "minLength": 0,
                                                    "maxLength": 128
                                                },
                                                "MinWalSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "MaxWalSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "MaxConnections": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "CheckpointCompletionTarget": {
                                                    "type": "number"
                                                },
                                                "DefaultStatisticsTarget": {
                                                    "type": "number"
                                                },
                                                "RandomPageCost": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "MaintenanceWorkMem": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "SharedBuffers": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "EffectiveCacheSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "EffectiveIOConcurrency": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "MaxWorkerProcesses": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "MaxParallelWorkers": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "Databases": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "Name",
                                                            "Replication",
                                                            "ReplicationDelay",
                                                            "EffectiveCacheSize",
                                                            "Version",
                                                            "Size",
                                                            "Collation",
                                                            "CharacterSet",
                                                            "Schemas"
                                                        ],
                                                        "properties": {
                                                            "Name": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "Replication": {
                                                                "type": "boolean"
                                                            },
                                                            "ReplicationDelay": {
                                                                "type": "boolean"
                                                            },
                                                            "EffectiveCacheSize": {
                                                                "type": "integer",
                                                                "minimum": 0
                                                            },
                                                            "Version": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "Size": {
                                                                "type": "number",
                                                                "minimum": 0
                                                            },
                                                            "Collation": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "CharacterSet": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "Schemas": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "Name",
                                                                        "Size",
                                                                        "TableCount"
                                                                    ],
                                                                    "properties": {
                                                                        "Name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "Size": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "TableCount": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        }
                                                                    }
                                                                }
                                                            }
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                }
            }
        },
        "Filesystems": {
            "type": "array",
            "items": {
                "type": "object",
                "required": [
                    "Filesystem",
                    "Type",
                    "Size",
                    "UsedSpace",
                    "AvailableSpace",
                    "MountedOn"
                ],
                "properties": {
                    "Filesystem": {
                        "type": "string",
                        "minLength": 1,
                        "maxLength": 64
                    },
                    "Type": {
                        "type": "string",
                        "minLength": 1,
                        "maxLength": 16
                    },
                    "Size": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "Size in bytes"
                    },
                    "UsedSpace": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "Used in bytes"
                    },
                    "AvailableSpace": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "AvailableSpace in bytes"
                    },
                    "MountedOn": {
                        "type": "string",
                        "minLength": 1,
                        "maxLength": 64
                    }
                }
            }
        },
        "Clusters": {
            "anyOf": [{
                    "type": "null"
                },
                {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": [
                            "FetchEndpoint",
                            "Type",
                            "Name",
                            "CPU",
                            "Sockets",
                            "VMs"
                        ],
                        "properties": {
                            "FetchEndpoint": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 64
                            },
                            "Type": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 16
                            },
                            "Name": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 128
                            },
                            "CPU": {
                                "type": "integer",
                                "minimum": 0
                            },
                            "Sockets": {
                                "type": "integer",
                                "minimum": 0
                            },
                            "VMs": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "required": [
                                        "Name",
                                        "Hostname",
                                        "CappedCPU",
                                        "VirtualizationNode"
                                    ],
                                    "properties": {
                                        "Name": {
                                            "type": "string",
                                            "minLength": 1,
                                            "maxLength": 128
                                        },
                                        "Hostname": {
                                            "type": "string",
                                            "minLength": 1,
                                            "maxLength": 253,
                                            "format": "idn-hostname"
                                        },
                                        "CappedCPU": {
                                            "type": "boolean"
                                        },
                                        "VirtualizationNode": {
                                            "type": "string",
                                            "minLength": 1,
                                            "maxLength": 253,
                                            "format": "idn-hostname"
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            ]
        }
    }
}
`
