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

package io.ercole.test.services;

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.when;

import java.text.ParseException;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang.time.DateUtils;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import io.ercole.model.CurrentHost;
import io.ercole.model.HistoricalHost;
import io.ercole.repositories.AlertRepository;
import io.ercole.repositories.ClusterRepository;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.repositories.HistoricalHostRepository;
import io.ercole.services.HostService;
import io.ercole.utilities.JsonFilter;
import org.json.JSONArray;
import org.json.JSONObject;

/**
 * Tests for HostService
 */

@RunWith(PowerMockRunner.class)
@PrepareForTest(JsonFilter.class)
public class HostServiceTest {
	
	@Mock
	private HistoricalHostRepository historicalRepo;

	@Mock
	private CurrentHostRepository currentRepo;
	
	@Mock
	private ClusterRepository clusterRepo;
	
	@Mock
	private AlertRepository alertRepo;
		
	@InjectMocks
	private HostService hostService;
	
	String oldDbArray = "[{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"orcl\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"orcl\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.38\",\"Used\":\"452.625\",\"Total\":\"480\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"12.1875\",\"Total\":\"70\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,557\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"orcl\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"MIODB\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"MIODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.38\",\"Used\":\"450.5625\",\"Total\":\"480\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"14.1875\",\"Total\":\"70\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"MIODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"TUODB\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":false,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":true,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"TUODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"TUODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"AAAA\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":false,\"Name\":\"Spatial and Graph\"},{\"Status\":true,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"AAAA\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"AAAA\",\"NCharset\":\"AL16UTF16\"}]";
	String extraInfo = "{\"extraInfo\":{\"Databases\":[{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"orcl\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"orcl\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.38\",\"Used\":\"452.625\",\"Total\":\"480\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"12.1875\",\"Total\":\"70\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,557\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"orcl\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"MIODB\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"MIODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.38\",\"Used\":\"450.5625\",\"Total\":\"480\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"14.1875\",\"Total\":\"70\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"MIODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"TUODB\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":false,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"TUODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"TUODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"AAAA\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":false,\"Name\":\"Spatial and Graph\"},{\"Status\":true,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"AAAA\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"AAAA\",\"NCharset\":\"AL16UTF16\"}],\"Filesystems\":[{\"UsedPerc\":\"70%\",\"Used\":\"33G\",\"Size\":\"50G\",\"FsType\":\"ext4\",\"Available\":\"15G\",\"Filesystem\":\"/dev/mapper/fedora-root\",\"MountedOn\":\"/\"},{\"UsedPerc\":\"19%\",\"Used\":\"171M\",\"Size\":\"976M\",\"FsType\":\"ext4\",\"Available\":\"738M\",\"Filesystem\":\"/dev/sda2\",\"MountedOn\":\"/boot\"},{\"UsedPerc\":\"5%\",\"Used\":\"8.9M\",\"Size\":\"200M\",\"FsType\":\"vfat\",\"Available\":\"191M\",\"Filesystem\":\"/dev/sda1\",\"MountedOn\":\"/boot/efi\"},{\"UsedPerc\":\"0%\",\"Used\":\"0\",\"Size\":\"7.8G\",\"FsType\":\"devtmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"devtmpfs\",\"MountedOn\":\"/dev\"},{\"UsedPerc\":\"2%\",\"Used\":\"110M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.7G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/dev/shm\"},{\"UsedPerc\":\"91%\",\"Used\":\"147G\",\"Size\":\"171G\",\"FsType\":\"ext4\",\"Available\":\"17G\",\"Filesystem\":\"/dev/mapper/fedora-home\",\"MountedOn\":\"/home\"},{\"UsedPerc\":\"1%\",\"Used\":\"2.2M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run\"},{\"UsedPerc\":\"1%\",\"Used\":\"52K\",\"Size\":\"1.6G\",\"FsType\":\"tmpfs\",\"Available\":\"1.6G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run/user/1000\"},{\"UsedPerc\":\"1%\",\"Used\":\"12K\",\"Size\":\"1.6G\",\"FsType\":\"tmpfs\",\"Available\":\"1.6G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run/user/42\"},{\"UsedPerc\":\"0%\",\"Used\":\"0\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/sys/fs/cgroup\"},{\"UsedPerc\":\"1%\",\"Used\":\"18M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/tmp\"}]}}";
	
	
	@Test
	public void updateWithAgentExistingHostnameInTimeRangeResultsUPDATED() throws ParseException {
		PowerMockito.mockStatic(JsonFilter.class);
		
		JSONObject json1 = new JSONObject("{\"Hostname\":\"host2\",\"Databases\":\"DB1\",\"Info\":{},\"Extra\":{\"Databases\":["
				+ "{\"prova\":\"prova\","
				+ "\"Name\":\"DB1\","
				+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
				+ "}]},"
				+ "\"Environment\":\"\",\"Location\":\"\",\"Schemas\":\"\"}");
		
		Date firstEntry = new Date(15000000000l);
		CurrentHost currentHost1 = new CurrentHost(Long.valueOf("1"), "host2", "PRD", "Italia", "oracledb", "", "", "DB1", 
				"ADMIN", "{\"Databases\":["
						+ "{\"prova\":\"prova\","
						+ "\"Name\":\"DB1\","
						+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
						+ "}]}", null,"info", firstEntry);

		when(currentRepo.findByHostname(json1.getString("Hostname"))).thenReturn(currentHost1);
		when(JsonFilter.buildCurrentHostFromJSON(json1)).thenReturn(currentHost1);
		when(clusterRepo.findOneVMInfoByHostname(json1.getString("Hostname"))).thenReturn(null);
		String status1 = hostService.updateWithAgent(json1, "oracledb");
		assertEquals("updated", status1);
	}
	
	@Test
	public void updateWithNoDatabaseInNewUpdate() throws ParseException {
		PowerMockito.mockStatic(JsonFilter.class);
		
		JSONObject json1 = new JSONObject("{\"Hostname\":\"host2\",\"Databases\":\"\",\"Info\":{},\"Extra\":{\"Databases\":["
				+ "]},"
				+ "\"Environment\":\"\",\"Location\":\"\",\"Schemas\":\"\"}");
		
		
		Date firstEntry = new Date(15000000000l);
		CurrentHost currentHost1 = new CurrentHost(Long.valueOf("1"), "host2", "PRD", "Italia", "oracledb", "", "", "DB1", 
				"ADMIN", "{\"Databases\":["
						+ "{\"prova\":\"prova\","
						+ "\"Name\":\"DB1\","
						+ "\"Features\":[{\"Status\":false,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
						+ "}]}", null, "info", firstEntry);
		
		CurrentHost currentHost2 = new CurrentHost(Long.valueOf("1"), "host2", "PRD", "Italia", "oracledb", "", "", "DB2",
				"ADMIN", "{\"Databases\":["
						+ "{\"prova\":\"prova\","
						+ "\"Name\":\"DB2\","
						+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
						+ "}]}", null, "info", DateUtils.addSeconds(firstEntry, 30));

		when(currentRepo.findByHostname(json1.getString("Hostname"))).thenReturn(currentHost1);
		when(clusterRepo.findOneVMInfoByHostname(json1.getString("Hostname"))).thenReturn(null);
		when(JsonFilter.buildCurrentHostFromJSON(json1)).thenReturn(currentHost2);
		
		List<String> newDbs = new ArrayList<>();
		newDbs.add("DB2");
		when(JsonFilter.getNewDatabases(currentHost2, currentHost1)).thenReturn(newDbs);
		when(JsonFilter.hasMoreCPUCores(currentHost1, currentHost2)).thenReturn(true);
		
		List<String> oldDbs = new ArrayList<>();
		newDbs.add("DB1");
			
		String status1 = hostService.updateWithAgent(json1, "oracledb");
		assertEquals("updated", status1);
	}
	

	
	@Test
	public void updateWithNewDatabaseInNewUpdate() throws ParseException {
		PowerMockito.mockStatic(JsonFilter.class);
		
		JSONObject json1 = new JSONObject("{\"Hostname\":\"host2\",\"Databases\":\"\",\"Info\":{},\"Extra\":{\"Databases\":["
				+ "{\"prova\":\"prova\","
				+ "\"Name\":\"DB2\","
				+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
				+ "}]},"
				+ "\"Environment\":\"\",\"Location\":\"\",\"Schemas\":\"\"}");
		
		
		Date firstEntry = new Date(15000000000l);
		CurrentHost currentHost1 = new CurrentHost(Long.valueOf("1"), "host2", "PRD", "Italia", "oracledb", "", "", "DB1", 
				"ADMIN", "{\"Databases\":["
						+ "{\"prova\":\"prova\","
						+ "\"Name\":\"DB1\","
						+ "\"Features\":[{\"Status\":false,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
						+ "}]}", null,  "info", firstEntry);
		
		CurrentHost currentHost2 = new CurrentHost(Long.valueOf("1"), "host2", "PRD", "Italia", "oracledb", "", "", "DB2", 
				"ADMIN", "{\"Databases\":["
						+ "{\"prova\":\"prova\","
						+ "\"Name\":\"DB2\","
						+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}]"
						+ "}]}", null, "info", DateUtils.addSeconds(firstEntry, 30));

		when(currentRepo.findByHostname(json1.getString("Hostname"))).thenReturn(currentHost1);
		when(clusterRepo.findOneVMInfoByHostname(json1.getString("Hostname"))).thenReturn(null);
		when(JsonFilter.buildCurrentHostFromJSON(json1)).thenReturn(currentHost2);
		
		List<String> newDbs = new ArrayList<>();
		newDbs.add("DB2");
		when(JsonFilter.getNewDatabases(currentHost2, currentHost1)).thenReturn(newDbs);
		when(JsonFilter.hasMoreCPUCores(currentHost1, currentHost2)).thenReturn(true);
		
		List<String> oldDbs = new ArrayList<>();
		oldDbs.add("DB1");
		
		
		Map<String, Map<String, Boolean>> newDbArrayWithFeatures = new HashMap<>();
		Map<String, Boolean> newFeatureMap = new HashMap<>();
		newFeatureMap.put("WebLogic Server Management Pack Enterprise Edition", true);
		newDbArrayWithFeatures.put("DB2", newFeatureMap);
		
		Map<String, Map<String, Boolean>> oldDbArrayWithFeatures = new HashMap<>();
		Map<String, Boolean> oldFeatureMap = new HashMap<>();
		oldFeatureMap.put("WebLogic Server Management Pack Enterprise Edition", false);
		newDbArrayWithFeatures.put("DB1", oldFeatureMap);
		
		when(JsonFilter.getFeaturesMapping(new JSONObject(currentHost2.getExtraInfo()).getJSONArray("Databases")))
		.thenReturn(newDbArrayWithFeatures);
		when(JsonFilter.getFeaturesMapping(new JSONObject(currentHost1.getExtraInfo()).getJSONArray("Databases")))
		.thenReturn(oldDbArrayWithFeatures);
			
		String status1 = hostService.updateWithAgent(json1, "oracledb");
		assertEquals("updated", status1);
	}
	
	
	

	
	
	
	// @Test
	// public void updateWithAgentNewHostnameResultsINSERTED() throws ParseException {
	// 	JSONObject json = new JSONObject("{\"Hostname\":\"host2\",\"Databases\":\"DB1 DB2\",\"Info\":{},\"Extra\":{\"Databases\":["
	// 			+ "{\"prova\":\"prova\","
	// 			+ "\"Name\":\"DB1\","
	// 			+ "\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"}],"
	// 			+ "\"Licenses\":["
	// 			+ "{\"Name\":\"Oracle EXE\",\"Count\":0},"
	// 			+ "{\"Name\":\"Oracle ENT\",\"Count\":1},"
	// 			+ "{\"Name\":\"Oracle STD\",\"Count\":0},"
	// 			+ "{\"Name\":\"WebLogic Server Management Pack Enterprise Edition\",\"Count\":0}"
	// 			+ "]}]},"
	// 			+ "\"Environment\":\"\",\"Location\":\"\",\"Schemas\":\"\"}");
	// 	// assertNotNull(hostService);
	// 	// assertNotNull(json);
	// 	assertEquals("inserted", 
	// 	hostService.
	// 	updateWithAgent(
	// 		json, 
	// 		"oracledb"));
	// }
	
	
	
	
	// @Test
	// public void updateWithAgentFloodingResultsERROR() throws ParseException {
	// 	JSONObject json3 = new JSONObject("{\"Hostname\":\"host3\",\"Info\":{},\"Extra\":{},\"Databases\":\"DB1 DB2\","
	// 			+ "\"Environment\":\"\",\"Location\":\"\",\"Schemas\":\"\"}");
		
	// 	Date wrongUpdate = new Date();
	// 	CurrentHost currentHost3 = new CurrentHost(Long.valueOf("1"), "host3", "PRD", "Italia", "oracledb", "", "", "BRB CCC", 
	// 			"ADMIN",  "info", null, "info", wrongUpdate);

	// 	when(currentRepo.findByHostname((String) json3.get("Hostname"))).thenReturn(currentHost3);
	// 	when(clusterRepo.findOneVMInfoByHostname(json3.getString("Hostname"))).thenReturn(null);

	// 	hostService.setUpdateRate(8000);
	// 	String status3 = hostService.updateWithAgent(json3, "oracledb");

	// 	assertEquals("error", status3);
	// }
	
	@Test
	public void getHistoricalLogsFromNoHistory() {	
		Date date = new Date(150000000000l);
		when(historicalRepo.findFirstHostnameByArchivedDesc("hostname", date)).thenReturn(new ArrayList<>());
		HistoricalHost historical = hostService.getHistoricalLogs("hostname", date);
		assertEquals(null, historical);	
	}
	
	
	@Test
	public void getHistoricalLogsFromYesHistory() {	
		Date date = new Date(150000000000l);
		HistoricalHost historical1 = new HistoricalHost(1l, "hostname1", "PRD", "Italy", "oracledb", "", "", "DB1", "ADMIN", "extraInfo", null, "hostInfo", new Date());
		List<HistoricalHost> list = new LinkedList<>();
		list.add(historical1);
		
		when(historicalRepo.findFirstHostnameByArchivedDesc("hostname1", date)).thenReturn(list);
		HistoricalHost retVal = hostService.getHistoricalLogs("hostname1", date);
		assertEquals(historical1, retVal);
	}
	

}
