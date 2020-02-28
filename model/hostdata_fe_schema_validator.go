// Copyright (c) 2019 Sorint.lab S.p.A.
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
		"title": "FEHostData",
		"description": "A hostdata from FE",
		"type": "object",
		"required": [
			"Hostname",
			"Environment",
			"Location",
			"Version",
			"HostDataSchemaVersion",
			"Info",
			"Extra"
		],
		"properties": {
			"Hostname": { "type": "string" },
			"Environment": { "type": "string" },
			"Location": { "type": "string" },
			"HostType": {
				"enum": [ 
					"oracledb",
					"virtualization",
					"exadata"
				]
			},
			"Version": { "type": "string" },
			"HostDataSchemaVersion": { "type": "integer" }, 
			"Databases": { "type": "string" },
			"Schemas": { "type": "string" },
			"Info": {
				"type": "object",
				"required": [
					"Hostname",
					"Environment",
					"Location",
					"CPUModel",
					"CPUCores",
					"CPUThreads",
					"Socket",
					"Type",
					"Virtual",
					"Kernel",
					"OS",
					"MemoryTotal",
					"SwapTotal",
					"OracleCluster",
					"VeritasCluster",
					"SunCluster",
					"AixCluster"
				],
				"properties": {
					"Hostname": {
						"type": "string",
						"pattern": "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
					},
					"Environment": {
						"type": "string"
					},
					"Location": {
						"type": "string"
					},
					"CPUModel": {
						"type": "string"
					},
					"CPUCores": {
						"type": "integer"
					},
					"CPUThreads": {
						"type": "integer"
					},
					"Socket": {
						"type": "integer"
					},
					"Type": {
						"type": "string"
					},
					"Virtual": {
						"type": "boolean"
					},
					"Kernel": {
						"type": "string"
					},
					"OS": {
						"type": "string"
					},
					"MemoryTotal": {
						"type": "number"
					},
					"SwapTotal": {
						"type": "number"
					},
					"OracleCluster": {
						"type": "boolean"
					},
					"VeritasCluster": {
						"type": "boolean"
					},
					"SunCluster": {
						"type": "boolean"
					},
					"AixCluster": {
						"type": "boolean"
					}
				}
			},
			"Extra": {
				"type": "object",
				"required": [
					"Filesystems"
				],
				"properties": {
					"Filesystems": {
						"type": "array",
						"items": {
							"type": "object",
							"required": [
								"Filesystem",
								"FsType",
								"Size",
								"Used",
								"Available",
								"UsedPerc",
								"MountedOn"
							],
							"properties": {
								"Filesystem": {
									"type": "string"
								},
								"FsType": {
									"type": "string"
								},
								"Size": {
									"type": "string"
								},
								"Used": {
									"type": "string"
								},
								"UsedPerc": {
									"type": "string"
								},
								"MountedOn": {
									"type": "string"
								}
							}
						}
					},
					"Databases": {
						"anyOf": [
							{ "type": "null" },
							{
								"type": "array",
								"items": {
									"type": "object",
									"required": [
										"InstanceNumber",
										"Name",
										"UniqueName",
										"Status",
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
										"Used",
										"DailyCPUUsage",
										"Allocated",
										"Elapsed",
										"DBTime",
										"Work",
										"ASM",
										"Dataguard",
										"Patches",
										"Tablespaces",
										"Schemas",
										"Features",
										"Licenses",
										"ADDMs",
										"SegmentAdvisors",
										"LastPSUs",
										"Backups"
									],
									"properties": {
										"InstanceNumber": {
											"type": "string"
										},
										"Name": {
											"type": "string"
										},
										"UniqueName": {
											"type": "string"
										},
										"Status": {
											"type": "string"
										},
										"Version": {
											"type": "string"
										},
										"Platform": {
											"type": "string"
										},
										"Archivelog": {
											"type": "string"
										},
										"Charset": {
											"type": "string"
										},
										"NCharset": {
											"type": "string"
										},
										"BlockSize": {
											"type": "string"
										},
										"CPUCount": {
											"type": "string"
										},
										"SGATarget": {
											"type": "string"
										},
										"PGATarget": {
											"type": "string"
										},
										"MemoryTarget": {
											"type": "string"
										},
										"SGAMaxSize": {
											"type": "string"
										},
										"SegmentsSize": {
											"type": "string"
										},
										"Used": {
											"type": "string"
										},
										"Allocated": {
											"type": "string"
										},
										"Elapsed": {
											"type": "string"
										},
										"DBTime": {
											"type": "string"
										},
										"DailyCPUUsage": {
											"type": "string"
										},
										"Work": {
											"type": "string"
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
													"Database",
													"Version",
													"PatchID",
													"Action",
													"Description",
													"Date"
												],
												"properties": {
													"Database": {
														"type": "string"
													},
													"Version": {
														"type": "string"
													},
													"PatchID": {
														"type": "string"
													},
													"Action": {
														"type": "string"
													},
													"Description": {
														"type": "string"
													},
													"Date": {
														"type": "string"
													}
												}
											}
										},
										"Tablespaces": {
											"type": "array",
											"items": {
												"type": "object",
												"required": [
													"Database",
													"Name",
													"MaxSize",
													"Total",
													"Used",
													"UsedPerc",
													"Status"
												],
												"properties": {
													"Database": {
														"type": "string"
													},
													"Name": {
														"type": "string"
													},
													"MaxSize": {
														"type": "string"
													},
													"Total": {
														"type": "string"
													},
													"Used": {
														"type": "string"
													},
													"UsedPerc": {
														"type": "string"
													},
													"Status": {
														"type": "string"
													}
												}
											}
										},
										"Schemas": {
											"type": "array",
											"items": {
												"type": "object",
												"required": [
													"Database",
													"User",
													"Total",
													"Tables",
													"Indexes",
													"LOB"
												],
												"properties": {
													"Database": {
														"type": "string"
													},
													"User": {
														"type": "string"
													},
													"Total": {
														"type": "integer"
													},
													"Tables": {
														"type": "integer"
													},
													"Indexes": {
														"type": "integer"
													},
													"LOB": {
														"type": "integer"
													}
												}
											}
										},
										"Features": {
											"type": "array",
											"items": {
												"type": "object",
												"required": [
													"Name",
													"Status"										
												],
												"properties": {
													"Name": {
														"type": "string"
													},
													"Status": {
														"type": "boolean"
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
														"type": "string"
													},
													"Count": {
														"type": "number"
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
														"type": "string"
													},
													"Recommendation": {
														"type": "string"
													},
													"Action": {
														"type": "string"
													},
													"Benefit": {
														"type": "string"
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
														"type": "string"
													},
													"SegmentName": {
														"type": "string"
													},
													"SegmentType": {
														"type": "string"
													},
													"PartitionName": {
														"type": "string"
													},
													"Reclaimable": {
														"type": "string"
													},
													"Recommendation": {
														"type": "string"
													}
												}
											}
										},
										"LastPSUs": {
											"type": "array",
											"items": {
												"type": "object",
												"required": [
													"Date",
													"Description"
												],
												"properties": {
													"Date": {
														"type": "string"
													},
													"Description": {
														"type": "string"
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
														"type": "string"
													},
													"Hour": {
														"type": "string"
													},
													"WeekDays": {
														"type": "string"
													},
													"AvgBckSize": {
														"type": "string"
													},
													"Retention": {
														"type": "string"
													}
												}
											}
										}
									}
								}
							}
						]
					},
					"Clusters": {
						"anyOf": [
							{ "type": "null" },
							{ 
								"type": "array",
								"items": {
									"type": "object",
									"required": [
										"Name",
										"Type",
										"CPU",
										"Sockets",
										"VMs"
									],
									"properties": {
										"Name": {
											"type": "string"
										},
										"Type": {
											"type": "string"
										},
										"CPU": {
											"type": "integer"
										},
										"Sockets": {
											"type": "integer"
										},
										"VMs": { 
											"type": "array",
											"items": {
												"type": "object",
												"required": [
													"Name",
													"ClusterName",
													"Hostname",
													"CappedCPU",
													"PhysicalHost"
												],
												"properties": {
													"Name": {
														"type": "string"
													},
													"ClusterName": {
														"type": "string"
													},
													"Hostname": {
														"type": "string"
													},
													"CappedCPU": {
														"type": "boolean"
													},
													"PhysicalHost": {
														"type": "string"
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
						"anyOf": [
							{ "type": "null" },
							{ 
								"type": "object",
								"required": [
									"Devices"
								],
								"properties": {
									"Devices": {
										"type": "array",
										"items": {
											"type": "object",
											"required": [
												"Hostname",
												"ServerType",
												"Model",
												"ExaSwVersion",
												"CPUEnabled",
												"Memory",
												"Status",
												"PowerCount",
												"PowerStatus",
												"FanCount",
												"FanStatus",
												"TempActual",
												"TempStatus",
												"CellsrvService",
												"MsService",
												"RsService",
												"FlashcacheMode",
												"CellDisks"
											],
											"properties": {
												"Hostname": {
													"type": "string"
												},
												"ServerType": {
													"type": "string"
												},
												"Model": {
													"type": "string"
												},
												"ExaSwVersion": {
													"type": "string"
												},
												"CPUEnabled": {
													"type": "string"
												},
												"Memory": {
													"type": "string"
												},
												"Status": {
													"type": "string"
												},
												"PowerCount": {
													"type": "string"
												},
												"PowerStatus": {
													"type": "string"
												},
												"FanCount": {
													"type": "string"
												},
												"FanStatus": {
													"type": "string"
												},
												"TempActual": {
													"type": "string"
												},
												"TempStatus": {
													"type": "string"
												},
												"CellsrvService": {
													"type": "string"
												},
												"MsService": {
													"type": "string"
												},
												"RsService": {
													"type": "string"
												},
												"FlashcacheMode": {
													"type": "string"
												},
												"CellDisks": {
													"anyOf": [
														{ "type": "null" },
														{
															"type": "array",
															"items": {
																"type": "object",
																"required": [
																	"Name",
																	"Status",
																	"ErrCount",
																	"UsedPerc"
																],
																"properties": {
																	"Name": {
																		"type": "string"
																	},
																	"Status": {
																		"type": "string"
																	},
																	"ErrCount": {
																		"type": "string"
																	},
																	"UsedPerc": {
																		"type": "string"
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
		}
	}
`
