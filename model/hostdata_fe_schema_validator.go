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

//TODO with go1.16 use go:embed https://github.com/golang/go/issues/41191

var FrontendHostdataSchemaValidator string = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$id": "ercole-hostdata",
    "type": "object",
    "required": [
        "hostname",
        "location",
        "environment",
        "tags",
        "agentVersion",
        "schemaVersion",
        "info",
        "clusterMembershipStatus",
        "features",
        "filesystems"
    ],
    "properties": {
        "hostname": {
            "type": "string",
            "minLength": 1,
            "maxLength": 253,
            "format": "idn-hostname"
        },
        "location": {
            "type": "string",
            "minLength": 1,
            "maxLength": 64
        },
        "environment": {
            "type": "string",
            "minLength": 1,
            "maxLength": 16
        },
        "tags": {
            "type": "array",
            "items": {
                "type": "string",
                "minLength": 1,
                "maxLength": 128
            },
            "uniqueItems": true
        },
        "agentVersion": {
            "type": "string",
            "minLength": 1,
            "maxLength": 64
        },
        "schemaVersion": {
            "type": "integer",
            "const": 1
        },
        "info": {
            "type": "object",
            "required": [
                "hostname",
                "cpuModel",
                "cpuFrequency",
                "cpuSockets",
                "cpuCores",
                "cpuThreads",
                "threadsPerCore",
                "coresPerSocket",
                "hardwareAbstraction",
                "hardwareAbstractionTechnology",
                "kernel",
                "kernelVersion",
                "os",
                "osVersion",
                "memoryTotal",
                "swapTotal"
            ],
            "properties": {
                "hostname": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 253,
                    "format": "idn-hostname"
                },
                "cpuModel": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64,
                    "pattern": "^[^\n]+$"
                },
                "cpuFrequency": {
                    "type": "string",
                    "minLength": 0,
                    "maxLength": 16
                },
                "cpuSockets": {
                    "type": "integer",
                    "minimum": 0
                },
                "cpuCores": {
                    "type": "integer",
                    "minimum": 1
                },
                "cpuThreads": {
                    "type": "integer",
                    "minimum": 1
                },
                "threadsPerCore": {
                    "type": "integer",
                    "minimum": 1
                },
                "coresPerSocket": {
                    "type": "integer",
                    "minimum": 1
                },
                "hardwareAbstraction": {
                    "type": "string",
                    "enum": [
                        "PH",
                        "VIRT"
                    ]
                },
                "hardwareAbstractionTechnology": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 32,
                    "pattern": "^[A-Z0-9]+$"
                },
                "kernel": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 128
                },
                "kernelVersion": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "os": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "osVersion": {
                    "type": "string",
                    "minLength": 1,
                    "maxLength": 64
                },
                "memoryTotal": {
                    "type": "number",
                    "minimum": 0
                },
                "swapTotal": {
                    "type": "number",
                    "minimum": 0
                }
            }
        },
        "clusterMembershipStatus": {
            "type": "object",
            "properties": {
                "oracleClusterware": {
                    "type": "boolean"
                },
                "veritasClusterServer": {
                    "type": "boolean"
                },
                "sunCluster": {
                    "type": "boolean"
                },
                "hacmp": {
                    "type": "boolean"
                }
            }
        },
        "features": {
            "type": "object",
            "properties": {
                "oracle": {
                    "anyOf": [
                        {
                            "type": "null"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "database": {
                                    "anyOf": [
                                        {
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "databases",
                                                "unlistedRunningDatabases"
                                            ],
                                            "properties": {
                                                "databases": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "instanceNumber",
                                                            "instanceName",
                                                            "name",
                                                            "uniqueName",
                                                            "status",
                                                            "dbID",
                                                            "role",
                                                            "isCDB",
                                                            "version",
                                                            "platform",
                                                            "archivelog",
                                                            "charset",
                                                            "nCharset",
                                                            "blockSize",
                                                            "cpuCount",
                                                            "sgaTarget",
                                                            "pgaTarget",
                                                            "memoryTarget",
                                                            "sgaMaxSize",
                                                            "segmentsSize",
                                                            "datafileSize",
                                                            "allocable",
                                                            "elapsed",
                                                            "dbTime",
                                                            "dailyCPUUsage",
                                                            "work",
                                                            "asm",
                                                            "dataguard",
                                                            "patches",
                                                            "tablespaces",
                                                            "schemas",
                                                            "licenses",
                                                            "addms",
                                                            "segmentAdvisors",
                                                            "psus",
                                                            "backups",
                                                            "featureUsageStats",
                                                            "pdbs",
                                                            "services"
                                                        ],
                                                        "properties": {
                                                            "instanceNumber": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "instanceName": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "name": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "uniqueName": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "status": {
                                                                "type": "string",
                                                                "enum": [
                                                                    "OPEN",
                                                                    "MOUNTED"
                                                                ]
                                                            },
                                                            "dbID": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "role": {
                                                                "type": "string",
                                                                "enum": [
                                                                    "PRIMARY",
                                                                    "LOGICAL STANDBY",
                                                                    "PHYSICAL STANDBY",
                                                                    "SNAPSHOT STANDBY"
                                                                ]
                                                            },
                                                            "isCDB": {
                                                                "type": "boolean"
                                                            },
                                                            "version": {
                                                                "type": "string",
                                                                "minLength": 8,
                                                                "maxLength": 64
                                                            },
                                                            "platform": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "archivelog": {
                                                                "type": "boolean"
                                                            },
                                                            "charset": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "nCharset": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "blockSize": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "cpuCount": {
                                                                "type": "integer",
                                                                "minimum": 1
                                                            },
                                                            "sgaTarget": {
                                                                "type": "number"
                                                            },
                                                            "pgaTarget": {
                                                                "type": "number"
                                                            },
                                                            "memoryTarget": {
                                                                "type": "number"
                                                            },
                                                            "sgaMaxSize": {
                                                                "type": "number"
                                                            },
                                                            "segmentsSize": {
                                                                "type": "number"
                                                            },
                                                            "datafileSize": {
                                                                "type": "number"
                                                            },
                                                            "allocable": {
                                                                "type": "number"
                                                            },
                                                            "elapsed": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "dbTime": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "dailyCPUUsage": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "work": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "asm": {
                                                                "type": "boolean"
                                                            },
                                                            "dataguard": {
                                                                "type": "boolean"
                                                            },
                                                            "patches": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "version",
                                                                        "patchID",
                                                                        "action",
                                                                        "description",
                                                                        "date"
                                                                    ],
                                                                    "properties": {
                                                                        "version": {
                                                                            "type": "string"
                                                                        },
                                                                        "patchID": {
                                                                            "type": "integer"
                                                                        },
                                                                        "action": {
                                                                            "type": "string"
                                                                        },
                                                                        "description": {
                                                                            "type": "string"
                                                                        },
                                                                        "date": {
                                                                            "type": "string",
                                                                            "format": "date"
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "tablespaces": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "name",
                                                                        "maxSize",
                                                                        "total",
                                                                        "used",
                                                                        "usedPerc",
                                                                        "status"
                                                                    ],
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string"
                                                                        },
                                                                        "maxSize": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "total": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "used": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "usedPerc": {
                                                                            "type": "number",
                                                                            "minimum": 0,
                                                                            "maximum": 100
                                                                        },
                                                                        "status": {
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
                                                            "schemas": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "indexes",
                                                                        "lob",
                                                                        "tables",
                                                                        "total",
                                                                        "user"
                                                                    ],
                                                                    "properties": {
                                                                        "indexes": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "lob": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "tables": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "total": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "user": {
                                                                            "type": "string"
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "licenses": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "name",
                                                                        "count"
                                                                    ],
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "count": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "addms": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "finding",
                                                                        "recommendation",
                                                                        "action",
                                                                        "benefit"
                                                                    ],
                                                                    "properties": {
                                                                        "finding": {
                                                                            "type": "string"
                                                                        },
                                                                        "recommendation": {
                                                                            "type": "string"
                                                                        },
                                                                        "action": {
                                                                            "type": "string"
                                                                        },
                                                                        "benefit": {
                                                                            "type": "number",
                                                                            "minimum": 0,
                                                                            "maximum": 100
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "segmentAdvisors": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "segmentOwner",
                                                                        "segmentName",
                                                                        "segmentType",
                                                                        "partitionName",
                                                                        "reclaimable",
                                                                        "recommendation"
                                                                    ],
                                                                    "properties": {
                                                                        "segmentOwner": {
                                                                            "type": "string"
                                                                        },
                                                                        "segmentName": {
                                                                            "type": "string"
                                                                        },
                                                                        "segmentType": {
                                                                            "type": "string"
                                                                        },
                                                                        "partitionName": {
                                                                            "type": "string"
                                                                        },
                                                                        "reclaimable": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "recommendation": {
                                                                            "type": "string"
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "psus": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "date",
                                                                        "description"
                                                                    ],
                                                                    "properties": {
                                                                        "date": {
                                                                            "type": "string",
                                                                            "format": "date"
                                                                        },
                                                                        "description": {
                                                                            "type": "string"
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "backups": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "backupType",
                                                                        "hour",
                                                                        "weekDays",
                                                                        "avgBckSize",
                                                                        "retention"
                                                                    ],
                                                                    "properties": {
                                                                        "backupType": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "Archivelog",
                                                                                "Full",
                                                                                "Level0",
                                                                                "Level1",
                                                                                "Incr Lvl 0",
                                                                                "Incr Lvl 1"
                                                                            ]
                                                                        },
                                                                        "hour": {
                                                                            "type": "string",
                                                                            "minLength": 5,
                                                                            "maxLength": 5,
                                                                            "pattern": "^[0-9]{2}:[0-9]{2}$"
                                                                        },
                                                                        "weekDays": {
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
                                                                        "avgBckSize": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "retention": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 16
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "featureUsageStats": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "product",
                                                                        "feature",
                                                                        "detectedUsages",
                                                                        "currentlyUsed",
                                                                        "firstUsageDate",
                                                                        "lastUsageDate",
                                                                        "extraFeatureInfo"
                                                                    ],
                                                                    "properties": {
                                                                        "product": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 128
                                                                        },
                                                                        "feature": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 128
                                                                        },
                                                                        "detectedUsages": {
                                                                            "type": "integer",
                                                                            "minimum": 0
                                                                        },
                                                                        "currentlyUsed": {
                                                                            "type": "boolean"
                                                                        },
                                                                        "firstUsageDate": {
                                                                            "type": "string",
                                                                            "format": "date-time"
                                                                        },
                                                                        "lastUsageDate": {
                                                                            "type": "string",
                                                                            "format": "date-time"
                                                                        },
                                                                        "extraFeatureInfo": {
                                                                            "type": "string",
                                                                            "maxLength": 128
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "services": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "name"
                                                                    ],
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 32
                                                                        }
                                                                    }
                                                                }
                                                            },
                                                            "pdbs": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "name",
                                                                        "status",
                                                                        "tablespaces",
                                                                        "schemas",
                                                                        "services"
                                                                    ],
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string"
                                                                        },
                                                                        "status": {
                                                                            "type": "string",
                                                                            "enum": [
                                                                                "OPEN",
                                                                                "READ WRITE",
                                                                                "READ ONLY",
                                                                                "MOUNTED"
                                                                            ]
                                                                        },
                                                                        "tablespaces": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "name",
                                                                                    "maxSize",
                                                                                    "total",
                                                                                    "used",
                                                                                    "usedPerc",
                                                                                    "status"
                                                                                ],
                                                                                "properties": {
                                                                                    "name": {
                                                                                        "type": "string"
                                                                                    },
                                                                                    "maxSize": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "total": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "used": {
                                                                                        "type": "number",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "usedPerc": {
                                                                                        "type": "number",
                                                                                        "minimum": 0,
                                                                                        "maximum": 100
                                                                                    },
                                                                                    "status": {
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
                                                                        "schemas": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "indexes",
                                                                                    "lob",
                                                                                    "tables",
                                                                                    "total",
                                                                                    "user"
                                                                                ],
                                                                                "properties": {
                                                                                    "indexes": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "lob": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "tables": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "total": {
                                                                                        "type": "integer",
                                                                                        "minimum": 0
                                                                                    },
                                                                                    "user": {
                                                                                        "type": "string"
                                                                                    }
                                                                                }
                                                                            }
                                                                        },
                                                                        "services": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "object",
                                                                                "required": [
                                                                                    "name"
                                                                                ],
                                                                                "properties": {
                                                                                    "name": {
                                                                                        "type": "string"
                                                                                    }
                                                                                }
                                                                            }
                                                                        }
                                                                    }
                                                                }
                                                            }
                                                        }
                                                    }
                                                },
                                                "unlistedRunningDatabases": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "string",
                                                        "minLength": 1
                                                    },
                                                    "uniqueItems": true
                                                }
                                            }
                                        }
                                    ]
                                },
                                "exadata": {
                                    "anyOf": [
                                        {
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "components"
                                            ],
                                            "properties": {
                                                "components": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "hostname",
                                                            "serverType",
                                                            "model",
                                                            "swVersion",
                                                            "swReleaseDate",
                                                            "runningCPUCount",
                                                            "totalCPUCount",
                                                            "memory",
                                                            "status",
                                                            "runningPowerSupply",
                                                            "totalPowerSupply",
                                                            "powerStatus",
                                                            "runningFanCount",
                                                            "totalFanCount",
                                                            "fanStatus",
                                                            "tempActual",
                                                            "tempStatus",
                                                            "cellsrvServiceStatus",
                                                            "msServiceStatus",
                                                            "rsServiceStatus",
                                                            "flashcacheMode",
                                                            "cellDisks"
                                                        ],
                                                        "properties": {
                                                            "hostname": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 253,
                                                                "format": "idn-hostname"
                                                            },
                                                            "serverType": {
                                                                "type": "string",
                                                                "enum": [
                                                                    "DBServer",
                                                                    "IBSwitch",
                                                                    "StorageServer"
                                                                ]
                                                            },
                                                            "model": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "swVersion": {
                                                                "type": "string"
                                                            },
                                                            "swReleaseDate": {
                                                                "type": "string",
                                                                "format": "date"
                                                            },
                                                            "runningCPUCount": {
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
                                                            "totalCPUCount": {
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
                                                            "memory": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "integer",
                                                                        "$comment": "memory in GB",
                                                                        "minimum": 1
                                                                    }
                                                                ]
                                                            },
                                                            "status": {
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
                                                            "runningPowerSupply": {
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
                                                            "totalPowerSupply": {
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
                                                            "powerStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "runningFanCount": {
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
                                                            "totalFanCount": {
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
                                                            "fanStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "tempActual": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "number"
                                                                    }
                                                                ]
                                                            },
                                                            "tempStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "cellsrvServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "msServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "rsServiceStatus": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "string"
                                                                    }
                                                                ]
                                                            },
                                                            "flashcacheMode": {
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
                                                            "cellDisks": {
                                                                "anyOf": [
                                                                    {
                                                                        "type": "null"
                                                                    },
                                                                    {
                                                                        "type": "array",
                                                                        "items": {
                                                                            "type": "object",
                                                                            "required": [
                                                                                "errCount",
                                                                                "name",
                                                                                "status",
                                                                                "usedPerc"
                                                                            ],
                                                                            "properties": {
                                                                                "errCount": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0
                                                                                },
                                                                                "name": {
                                                                                    "type": "string",
                                                                                    "minLength": 1,
                                                                                    "maxLength": 64
                                                                                },
                                                                                "status": {
                                                                                    "type": "string"
                                                                                },
                                                                                "usedPerc": {
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
                        }
                    ]
                },
                "postgresql": {
                    "anyOf": [
                        {
                            "type": "null"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "postgresql": {
                                    "anyOf": [
                                        {
                                            "type": "null"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "workMem",
                                                "archiveMode",
                                                "archivePath",
                                                "minWalSize",
                                                "maxWalSize",
                                                "maxConnections",
                                                "checkpointCompletionTarget",
                                                "defaultStatisticsTarget",
                                                "randomPageCost",
                                                "maintenanceWorkMem",
                                                "sharedBuffers",
                                                "effectiveCacheSize",
                                                "effectiveIOConcurrency",
                                                "maxWorkerProcesses",
                                                "maxParallelWorkers",
                                                "databases"
                                            ],
                                            "properties": {
                                                "workMem": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "archiveMode": {
                                                    "type": "boolean"
                                                },
                                                "archivePath": {
                                                    "type": "string",
                                                    "minLength": 0,
                                                    "maxLength": 128
                                                },
                                                "minWalSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "maxWalSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "maxConnections": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "checkpointCompletionTarget": {
                                                    "type": "number"
                                                },
                                                "defaultStatisticsTarget": {
                                                    "type": "number"
                                                },
                                                "randomPageCost": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "maintenanceWorkMem": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "sharedBuffers": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "effectiveCacheSize": {
                                                    "type": "number",
                                                    "minimum": 0
                                                },
                                                "effectiveIOConcurrency": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "maxWorkerProcesses": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "maxParallelWorkers": {
                                                    "type": "integer",
                                                    "minimum": 0
                                                },
                                                "databases": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "name",
                                                            "replication",
                                                            "replicationDelay",
                                                            "effectiveCacheSize",
                                                            "version",
                                                            "size",
                                                            "collation",
                                                            "characterSet",
                                                            "schemas"
                                                        ],
                                                        "properties": {
                                                            "name": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "replication": {
                                                                "type": "boolean"
                                                            },
                                                            "replicationDelay": {
                                                                "type": "boolean"
                                                            },
                                                            "effectiveCacheSize": {
                                                                "type": "integer",
                                                                "minimum": 0
                                                            },
                                                            "version": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "size": {
                                                                "type": "number",
                                                                "minimum": 0
                                                            },
                                                            "collation": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "characterSet": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "schemas": {
                                                                "type": "array",
                                                                "items": {
                                                                    "type": "object",
                                                                    "required": [
                                                                        "name",
                                                                        "size",
                                                                        "tableCount"
                                                                    ],
                                                                    "properties": {
                                                                        "name": {
                                                                            "type": "string",
                                                                            "minLength": 1,
                                                                            "maxLength": 64
                                                                        },
                                                                        "size": {
                                                                            "type": "number",
                                                                            "minimum": 0
                                                                        },
                                                                        "tableCount": {
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
                },
                "microsoft": {
                    "anyOf": [
                        {
                            "type": "null"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "sqlServer": {
                                    "anyOf": [
                                        {
                                            "type": "object"
                                        },
                                        {
                                            "type": "object",
                                            "required": [
                                                "instances",
                                                "features",
                                                "patches"
                                            ],
                                            "properties": {
                                                "instances": {
                                                    "type": "array",
                                                    "required": [
                                                        "status",
                                                        "name",
                                                        "displayName",
                                                        "serverName",
                                                        "databaseID",
                                                        "stateDesc",
                                                        "version",
                                                        "platform",
                                                        "collationName",
                                                        "edition",
                                                        "editionType",
                                                        "productCode",
                                                        "licensingInfo",
                                                        "databases"
                                                    ],
                                                    "properties": {
                                                        "status": {
                                                            "type": "string",
                                                            "enum": [
                                                                "Running",
                                                                "Stopped",
                                                                "ContinuePending",
                                                                "Paused",
                                                                "PausePending",
                                                                "StartPending",
                                                                "StopPending"
                                                            ]
                                                        },
                                                        "name": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 64
                                                        },
                                                        "displayName": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 64
                                                        },
                                                        "serverName": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 64
                                                        },
                                                        "databaseID": {
                                                            "type": "integer",
                                                            "minimum": 1
                                                        },
                                                        "stateDesc": {
                                                            "type": "string",
                                                            "enum": [
                                                                "ONLINE",
                                                                "RESTORING",
                                                                "RECOVERING",
                                                                "RECOVERY_PENDING",
                                                                "SUSPECT",
                                                                "EMERGENCY",
                                                                "OFFLINE",
                                                                "COPYING",
                                                                "OFFLINE_SECONDARY"
                                                            ]
                                                        },
                                                        "version": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 32
                                                        },
                                                        "platform": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 16
                                                        },
                                                        "collationName": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 32
                                                        },
                                                        "edition": {
                                                            "type": "string",
                                                            "enum": [
                                                                "ENT",
                                                                "EXP",
                                                                "STD",
                                                                "BI",
                                                                "DEV",
                                                                "WEB",
                                                                "AZU"
                                                            ]
                                                        },
                                                        "editionType": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 64
                                                        },
                                                        "productCode": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 64,
                                                            "pattern": "^\\{[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\\}$"
                                                        },
                                                        "licensingInfo": {
                                                            "type": "string",
                                                            "minLength": 1,
                                                            "maxLength": 512
                                                        },
                                                        "databases": {
                                                            "type": "array",
                                                            "items": {
                                                                "type": "object",
                                                                "required": [
                                                                    "databaseID",
                                                                    "name",
                                                                    "collationName",
                                                                    "status",
                                                                    "recoveryModel",
                                                                    "blockSize",
                                                                    "schedulersCount",
                                                                    "affinityMask",
                                                                    "minServerMemory",
                                                                    "maxServerMemory",
                                                                    "ctp",
                                                                    "maxDop",
                                                                    "alloc",
                                                                    "backups",
                                                                    "schemas",
                                                                    "tablespaces"
                                                                ],
                                                                "properties": {
                                                                    "databaseID": {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    },
                                                                    "name": {
                                                                        "type": "string",
                                                                        "minLength": 1,
                                                                        "maxLength": 64
                                                                    },
                                                                    "collationName": {
                                                                        "type": "string",
                                                                        "minLength": 1,
                                                                        "maxLength": 32
                                                                    },
                                                                    "status": {
                                                                        "type": "string",
                                                                        "enum": [
                                                                            "ONLINE",
                                                                            "RESTORING",
                                                                            "RECOVERING",
                                                                            "RECOVERY_PENDING",
                                                                            "SUSPECT",
                                                                            "EMERGENCY",
                                                                            "OFFLINE",
                                                                            "COPYING",
                                                                            "OFFLINE_SECONDARY"
                                                                        ]
                                                                    },
                                                                    "recoveryModel": {
                                                                        "type": "string",
                                                                        "enum": [
                                                                            "FULL",
                                                                            "BULK_LOGGED",
                                                                            "SIMPLE"
                                                                        ]
                                                                    },
                                                                    "blockSize": {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    },
                                                                    "schedulersCount": {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    },
                                                                    "affinityMask": {
                                                                        "type": "integer",
                                                                        "minimum": 0
                                                                    },
                                                                    "minServerMemory": {
                                                                        "type": "integer",
                                                                        "minimum": 1,
                                                                        "$comment": "size in bytes"
                                                                    },
                                                                    "maxServerMemory": {
                                                                        "type": "integer",
                                                                        "minimum": 1,
                                                                        "$comment": "size in bytes"
                                                                    },
                                                                    "ctp": {
                                                                        "type": "integer",
                                                                        "minimum": 1
                                                                    },
                                                                    "maxDop": {
                                                                        "type": "integer",
                                                                        "minimum": 0
                                                                    },
                                                                    "alloc": {
                                                                        "type": "number",
                                                                        "minimum": 0
                                                                    },
                                                                    "backups": {
                                                                        "type": "array",
                                                                        "items": {
                                                                            "type": "object",
                                                                            "required": [
                                                                                "backupType",
                                                                                "hour",
                                                                                "weekDays",
                                                                                "avgBckSize"
                                                                            ],
                                                                            "properties": {
                                                                                "backupType": {
                                                                                    "type": "string",
                                                                                    "enum": [
                                                                                        "Database",
                                                                                        "Log",
                                                                                        "File or filegroup",
                                                                                        "Differential database",
                                                                                        "Differential file",
                                                                                        "Differential partial",
                                                                                        "Partial"
                                                                                    ]
                                                                                },
                                                                                "hour": {
                                                                                    "type": "string",
                                                                                    "minLength": 5,
                                                                                    "maxLength": 5,
                                                                                    "pattern": "^[0-9]{2}:[0-9]{2}$"
                                                                                },
                                                                                "weekDays": {
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
                                                                                "avgBckSize": {
                                                                                    "type": "number",
                                                                                    "minimum": 0,
                                                                                    "$comment": "size in GB"
                                                                                }
                                                                            }
                                                                        }
                                                                    },
                                                                    "schemas": {
                                                                        "type": "array",
                                                                        "items": {
                                                                            "type": "object",
                                                                            "required": [
                                                                                "allocationType",
                                                                                "usedSpace",
                                                                                "allocatedSpace"
                                                                            ],
                                                                            "properties": {
                                                                                "allocationType": {
                                                                                    "type": "string",
                                                                                    "enum": [
                                                                                        "Clustered Index Data",
                                                                                        "LOB data",
                                                                                        "Non Clustered Index Data",
                                                                                        "Table Data"
                                                                                    ]
                                                                                },
                                                                                "usedSpace": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0,
                                                                                    "$comment": "size in bytes"
                                                                                },
                                                                                "allocatedSpace": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0,
                                                                                    "$comment": "size in bytes"
                                                                                }
                                                                            }
                                                                        }
                                                                    },
                                                                    "tablespaces": {
                                                                        "type": "array",
                                                                        "items": {
                                                                            "type": "object",
                                                                            "required": [
                                                                                "filename",
                                                                                "alloc",
                                                                                "used",
                                                                                "growth",
                                                                                "growthUnit",
                                                                                "fileType",
                                                                                "status"
                                                                            ],
                                                                            "properties": {
                                                                                "filename": {
                                                                                    "type": "string",
                                                                                    "minLength": 1,
                                                                                    "maxLength": 32
                                                                                },
                                                                                "alloc": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0,
                                                                                    "$comment": "size in bytes"
                                                                                },
                                                                                "used": {
                                                                                    "type": "integer",
                                                                                    "minimum": 0,
                                                                                    "$comment": "size in bytes"
                                                                                },
                                                                                "growth": {
                                                                                    "type": "number"
                                                                                },
                                                                                "growthUnit": {
                                                                                    "type": "string",
                                                                                    "enum": [
                                                                                        "%",
                                                                                        "MB"
                                                                                    ]
                                                                                },
                                                                                "fileType": {
                                                                                    "type": "string",
                                                                                    "enum": [
                                                                                        "ROWS",
                                                                                        "LOG",
                                                                                        "FILESTREAM",
                                                                                        "FULLTEXT"
                                                                                    ]
                                                                                },
                                                                                "status": {
                                                                                    "type": "string",
                                                                                    "enum": [
                                                                                        "ONLINE",
                                                                                        "RESTORING",
                                                                                        "RECOVERING",
                                                                                        "RECOVERY_PENDING",
                                                                                        "SUSPECT",
                                                                                        "OFFLINE",
                                                                                        "DEFUNCT"
                                                                                    ]
                                                                                }
                                                                            }
                                                                        }
                                                                    }
                                                                }
                                                            }
                                                        }
                                                    }
                                                },
                                                "features": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "product",
                                                            "feature",
                                                            "language",
                                                            "edition",
                                                            "version",
                                                            "clustered",
                                                            "configured"
                                                        ],
                                                        "properties": {
                                                            "product": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "feature": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "language": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 16
                                                            },
                                                            "edition": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "version": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 32
                                                            },
                                                            "clustered": {
                                                                "type": "boolean"
                                                            },
                                                            "configured": {
                                                                "type": "boolean"
                                                            }
                                                        }
                                                    }
                                                },
                                                "patches": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "required": [
                                                            "displayName",
                                                            "displayVersion",
                                                            "installDate"
                                                        ],
                                                        "properties": {
                                                            "displayName": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "displayVersion": {
                                                                "type": "string",
                                                                "minLength": 1,
                                                                "maxLength": 64
                                                            },
                                                            "installDate": {
                                                                "type": "string",
                                                                "format": "date"
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
        "filesystems": {
            "type": "array",
            "items": {
                "type": "object",
                "required": [
                    "filesystem",
                    "type",
                    "size",
                    "usedSpace",
                    "availableSpace",
                    "mountedOn"
                ],
                "properties": {
                    "filesystem": {
                        "type": "string",
                        "minLength": 1
                    },
                    "type": {
                        "type": "string",
                        "minLength": 1
                    },
                    "size": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "size in bytes"
                    },
                    "usedSpace": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "used in bytes"
                    },
                    "availableSpace": {
                        "type": "integer",
                        "minimum": 0,
                        "$comment": "availableSpace in bytes"
                    },
                    "mountedOn": {
                        "type": "string",
                        "minLength": 1
                    }
                }
            }
        },
        "clusters": {
            "anyOf": [
                {
                    "type": "null"
                },
                {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "required": [
                            "fetchEndpoint",
                            "type",
                            "name",
                            "cpu",
                            "sockets",
                            "vms"
                        ],
                        "properties": {
                            "fetchEndpoint": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 64
                            },
                            "type": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 16
                            },
                            "name": {
                                "type": "string",
                                "minLength": 1,
                                "maxLength": 128
                            },
                            "cpu": {
                                "type": "integer",
                                "minimum": 0
                            },
                            "sockets": {
                                "type": "integer",
                                "minimum": 0
                            },
                            "vms": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "required": [
                                        "name",
                                        "hostname",
                                        "cappedCPU",
                                        "virtualizationNode"
                                    ],
                                    "properties": {
                                        "name": {
                                            "type": "string",
                                            "minLength": 1,
                                            "maxLength": 128
                                        },
                                        "hostname": {
                                            "type": "string",
                                            "minLength": 1,
                                            "maxLength": 253,
                                            "format": "idn-hostname"
                                        },
                                        "cappedCPU": {
                                            "type": "boolean"
                                        },
                                        "virtualizationNode": {
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
