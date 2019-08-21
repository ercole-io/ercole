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

package io.ercole.test.controller;


import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;

import javax.servlet.http.HttpServletRequest;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.web.server.LocalServerPort;

import io.ercole.controller.AgentDataController;

/**
 * Test class for the HomeRestController, to be run only with the application up
 * and running.
 */

@RunWith(MockitoJUnitRunner.class)
//@RunWith(SpringRunner.class)
//@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
public class AgentDataControllerTest {
	
	@LocalServerPort
    private int port;
	
	@Autowired
	private TestRestTemplate restTemplate;

	@Mock
	HttpServletRequest request;

	@InjectMocks
	AgentDataController agentController;

	private String agentPost = "{\"Hostname\":\"TestingHost\", \"Environment\":\"TST\","
			+ " \"Databases\":\"MIAO\", \"Schemas\":\"SIP HR\", \"Extra\":{\"Databases\":"
			+ "[{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":"
			+ "[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},"
			+ "{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},"
			+ "{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},"
			+ "{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":"
			+ "\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},"
			+ "{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":"
			+ "\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},"
			+ "{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":"
			+ "\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},"
			+ "{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,"
			+ "\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,"
			+ "\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,"
			+ "\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},"
			+ "{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":"
			+ "\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},"
			+ "{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":"
			+ "\"Active Data Guard\"}],\"Name\":\"orcl\",\"Version\":\"12.2.0.1.0 Enterprise Edition\","
			+ "\"Work\":\"N/A\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management "
			+ "Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,"
			+ "\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,"
			+ "\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},"
			+ "{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":"
			+ "\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation "
			+ "Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},"
			+ "{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},"
			+ "{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},"
			+ "{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},"
			+ "{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},"
			+ "{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},"
			+ "{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},"
			+ "{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,"
			+ "\"Name\":\"Configuration Management Pack for Oracle Database\"},"
			+ "{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},"
			+ "{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},"
			+ "{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":"
			+ "[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\","
			+ "\"Version\":\"12.2.0.1.0\",\"Database\":\"orcl\",\"PatchID\":\" \",\"Date\":\" \"}],"
			+ "\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\","
			+ "\"UsedPerc\":\"1.38\",\"Used\":\"452.625\",\"Total\":\"480\",\"Database\":\"orcl\","
			+ "\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\","
			+ "\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},"
			+ "{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\","
			+ "\"Database\":\"orcl\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":"
			+ "\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"12.1875\",\"Total\":\"70\",\"Database\":\"orcl\","
			+ "\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":"
			+ "\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\","
			+ "\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\","
			+ "\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\","
			+ "\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\","
			+ "\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\","
			+ "\"Total\":0,\"Database\":\"orcl\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\","
			+ "\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,557\",\"Used\":\"2\",\"CPUCount\":\"4\","
			+ "\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"orcl\","
			+ "\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\","
			+ "\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},"
			+ "{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},"
			+ "{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},"
			+ "{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},"
			+ "{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},"
			+ "{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},"
			+ "{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},"
			+ "{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},"
			+ "{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},"
			+ "{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},"
			+ "{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},"
			+ "{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":"
			+ "\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},"
			+ "{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},"
			+ "{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],"
			+ "\"Name\":\"MIODB\",\"Version\":\"12.2.0.1.0 Enterprise Edition\",\"Work\":\"11\",\"Features\":"
			+ "[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},"
			+ "{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},"
			+ "{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},"
			+ "{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":"
			+ "\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},"
			+ "{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},"
			+ "{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{"
			+ "\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},"
			+ "{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},"
			+ "{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},"
			+ "{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},"
			+ "{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},"
			+ "{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},"
			+ "{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":"
			+ "\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},"
			+ "{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},"
			+ "{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],"
			+ "\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\","
			+ "\"Version\":\"12.2.0.1.0\",\"Database\":\"MIODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\","
			+ "\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.38\",\"Used\":\"450.5625\","
			+ "\"Total\":\"480\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\","
			+ "\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},"
			+ "{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.25\",\"Total\":\"800\",\"Database\":\"MIODB\",\"MaxSize\":"
			+ "\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.04\",\"Used\":\"14.1875\",\"Total\":\"70\","
			+ "\"Database\":\"MIODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"MIODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"MIODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"TUODB\",\"Version\":\"10.2.0.1.0 Enterprise Edition\",\"Work\":\"4\",\"Features\":[{\"Status\":false,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":true,\"Name\":\"Spatial and Graph\"},{\"Status\":false,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":true,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"TUODB\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"TUODB\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"TUODB\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"TUODB\",\"NCharset\":\"AL16UTF16\"},{\"Platform\":\"Linux x86 64-bit\",\"MemoryTarget\":\"0\",\"Licenses\":[{\"Count\":0,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Count\":0,\"Name\":\"Tuning Pack\"},{\"Count\":0,\"Name\":\"Spatial and Graph\"},{\"Count\":0,\"Name\":\"Secure Backup\"},{\"Count\":0,\"Name\":\"Real Application Testing\"},{\"Count\":0,\"Name\":\"Real Application Clusters\"},{\"Count\":0,\"Name\":\"Real Application Clusters One Node\"},{\"Count\":0,\"Name\":\"RAC or RAC One Node\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Count\":0,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Count\":0,\"Name\":\"Pillar Storage\"},{\"Count\":0,\"Name\":\"Partitioning\"},{\"Count\":0,\"Name\":\"OLAP\"},{\"Count\":0,\"Name\":\"Multitenant\"},{\"Count\":0,\"Name\":\"Label Security\"},{\"Count\":0,\"Name\":\"HW\"},{\"Count\":0,\"Name\":\"GoldenGate\"},{\"Count\":0,\"Name\":\"Exadata\"},{\"Count\":0,\"Name\":\"Diagnostics Pack\"},{\"Count\":0,\"Name\":\"Database Vault\"},{\"Count\":0,\"Name\":\"Database In-Memory\"},{\"Count\":0,\"Name\":\"Database Gateway\"},{\"Count\":0,\"Name\":\"Data Masking Pack\"},{\"Count\":0,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Count\":0,\"Name\":\"Change Management Pack\"},{\"Count\":0,\"Name\":\"Advanced Security\"},{\"Count\":0,\"Name\":\"Advanced Compression\"},{\"Count\":0,\"Name\":\"Advanced Analytics\"},{\"Count\":0,\"Name\":\"Active Data Guard\"}],\"Name\":\"AAAA\",\"Version\":\"10.2.0.1.0 Enterprise Edition\",\"Work\":\"2\",\"Features\":[{\"Status\":true,\"Name\":\"WebLogic Server Management Pack Enterprise Edition\"},{\"Status\":false,\"Name\":\"Tuning Pack\"},{\"Status\":false,\"Name\":\"Spatial and Graph\"},{\"Status\":true,\"Name\":\"Secure Backup\"},{\"Status\":false,\"Name\":\"Real Application Testing\"},{\"Status\":false,\"Name\":\"Real Application Clusters One Node\"},{\"Status\":false,\"Name\":\"Real Application Clusters\"},{\"Status\":false,\"Name\":\"RAC or RAC One Node\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack for Database\"},{\"Status\":false,\"Name\":\"Provisioning and Patch Automation Pack\"},{\"Status\":false,\"Name\":\"Pillar Storage\"},{\"Status\":false,\"Name\":\"Partitioning\"},{\"Status\":false,\"Name\":\"OLAP\"},{\"Status\":false,\"Name\":\"Multitenant\"},{\"Status\":false,\"Name\":\"Label Security\"},{\"Status\":false,\"Name\":\"HW\"},{\"Status\":false,\"Name\":\"GoldenGate\"},{\"Status\":false,\"Name\":\"Exadata\"},{\"Status\":false,\"Name\":\"Diagnostics Pack\"},{\"Status\":false,\"Name\":\"Database Vault\"},{\"Status\":false,\"Name\":\"Database In-Memory\"},{\"Status\":false,\"Name\":\"Database Gateway\"},{\"Status\":false,\"Name\":\"Data Masking Pack\"},{\"Status\":false,\"Name\":\"Configuration Management Pack for Oracle Database\"},{\"Status\":false,\"Name\":\"Change Management Pack\"},{\"Status\":false,\"Name\":\"Advanced Security\"},{\"Status\":false,\"Name\":\"Advanced Compression\"},{\"Status\":false,\"Name\":\"Advanced Analytics\"},{\"Status\":false,\"Name\":\"Active Data Guard\"}],\"ASM\":false,\"Patches\":[{\"Action\":\"BOOTSTRAP\",\"Description\":\"RDBMS_12.2.0.1.0_LINUX.X64_170125\",\"Version\":\"12.2.0.1.0\",\"Database\":\"AAAA\",\"PatchID\":\" \",\"Date\":\" \"}],\"Status\":\"OPEN\",\"Dataguard\":false,\"Tablespaces\":[{\"Status\":\"ONLINE\",\"UsedPerc\":\"1.36\",\"Used\":\"445.0625\",\"Total\":\"470\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSAUX\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.00\",\"Used\":\"0\",\"Total\":\"5\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"USERS\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"2.43\",\"Used\":\"795.125\",\"Total\":\"800\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"SYSTEM\"},{\"Status\":\"ONLINE\",\"UsedPerc\":\"0.15\",\"Used\":\"48.375\",\"Total\":\"70\",\"Database\":\"AAAA\",\"MaxSize\":\"32767.9844\",\"Name\":\"UNDOTBS1\"}],\"SGATarget\":\"4,672\",\"SGAMaxSize\":\"4,672\",\"Charset\":\"AL32UTF8\",\"Schemas\":[{\"User\":\"REMOTE_SCHEDULER_AGENT\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYS$UMF\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"GGSYS\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"DBSFWUSER\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0},{\"User\":\"SYSRAC\",\"Total\":0,\"Database\":\"AAAA\",\"Tables\":0,\"Indexes\":0,\"LOB\":0}],\"Allocated\":\"129\",\"Archivelog\":\"NOARCHIVELOG\",\"PGATarget\":\"1,556\",\"Used\":\"2\",\"CPUCount\":\"4\",\"Elapsed\":\"0\",\"DBTime\":\"0\",\"BlockSize\":\"8192\",\"UniqueName\":\"AAAA\",\"NCharset\":\"AL16UTF16\"}],\"Filesystems\":[{\"UsedPerc\":\"70%\",\"Used\":\"33G\",\"Size\":\"50G\",\"FsType\":\"ext4\",\"Available\":\"15G\",\"Filesystem\":\"/dev/mapper/fedora-root\",\"MountedOn\":\"/\"},{\"UsedPerc\":\"19%\",\"Used\":\"171M\",\"Size\":\"976M\",\"FsType\":\"ext4\",\"Available\":\"738M\",\"Filesystem\":\"/dev/sda2\",\"MountedOn\":\"/boot\"},{\"UsedPerc\":\"5%\",\"Used\":\"8.9M\",\"Size\":\"200M\",\"FsType\":\"vfat\",\"Available\":\"191M\",\"Filesystem\":\"/dev/sda1\",\"MountedOn\":\"/boot/efi\"},{\"UsedPerc\":\"0%\",\"Used\":\"0\",\"Size\":\"7.8G\",\"FsType\":\"devtmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"devtmpfs\",\"MountedOn\":\"/dev\"},{\"UsedPerc\":\"2%\",\"Used\":\"110M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.7G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/dev/shm\"},{\"UsedPerc\":\"91%\",\"Used\":\"147G\",\"Size\":\"171G\",\"FsType\":\"ext4\",\"Available\":\"17G\",\"Filesystem\":\"/dev/mapper/fedora-home\",\"MountedOn\":\"/home\"},{\"UsedPerc\":\"1%\",\"Used\":\"2.2M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run\"},{\"UsedPerc\":\"1%\",\"Used\":\"52K\",\"Size\":\"1.6G\",\"FsType\":\"tmpfs\",\"Available\":\"1.6G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run/user/1000\"},{\"UsedPerc\":\"1%\",\"Used\":\"12K\",\"Size\":\"1.6G\",\"FsType\":\"tmpfs\",\"Available\":\"1.6G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/run/user/42\"},{\"UsedPerc\":\"0%\",\"Used\":\"0\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/sys/fs/cgroup\"},{\"UsedPerc\":\"1%\",\"Used\":\"18M\",\"Size\":\"7.8G\",\"FsType\":\"tmpfs\",\"Available\":\"7.8G\",\"Filesystem\":\"tmpfs\",\"MountedOn\":\"/tmp\"}]} , \"Info\":{\"hostInfo\":{\"OS\":\"Fedora release 27 (Twenty Seven)\",\"SunCluster\":false,\"Hostname\":\"srota\",\"Virtual\":false,\"Type\":\"VMWARE\",\"VeritasCluster\":false,\"Environment\":\"AMBIENTE\",\"CPUThreads\":48,\"SwapTotal\":7,\"CPUCores\":24,\"Kernel\":\"4.16.5-200.fc27.x86_64\",\"Socket\":1,\"MemoryTotal\":15,\"OracleCluster\":false,\"CPUModel\":\"Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz\"} }, \"Location\":\"Italia\"}";
	
	String agentPostNoHostname = "{\"Hostname\":\"\"}";

	//TODO: da terminare
//	@Test
//    public void greetingShouldReturnDefaultMessage() throws Exception {
//        assertThat(this.restTemplate.getForObject("http://localhost:" + port + "/",
//                String.class)).contains("Hello World");
//    }

//	/**
//	 * Validation schema for the " / " JSON response.
//	 */
//	@Test
//	public void rootSchemaMatchesJSONSchema() {
//		get("/").then().assertThat().body(matchesJsonSchemaInClasspath("files/rootValidation.json"));
//	}


	/**
	 * "/macchina/update" should accept "application/json" media type.
	 */
	@Test
	public void updateAgentConsumingJSON() {
		post("/host/update").accept("application/json");
	}
//
//	/**
//	 * "/macchina/update" should require the right number of parameters.
//	 */
//	@Test
//	public void validateUpdateWithWrongPassword() {
//		Map<String, String> headerMap = new HashMap<>();
//		String username = "user";
//		String password = "Wrong_password";
//		String credentials = username + ":" + password;
//		byte[] src = credentials.getBytes(Charset.forName("UTF-8"));
//		String basicAuthString = "Basic " + new String(Base64.getEncoder().encodeToString(src));
//		
//		headerMap.put("Authorization", basicAuthString);
//		headerMap.put("Accept", "application/json");
//		headerMap.put("Content-type", "application/json");
//		
//		given().headers(headerMap).when().post("http://localhost:9080/host/update").then()
//				.statusCode(HttpStatus.UNAUTHORIZED.value());
//	}
//	
//	
//	@Test
//	public void validateUpdateWithNoJSON() {
//		Map<String, String> headerMap = new HashMap<>();
//		String username = "user";
//		String password = "password";
//		String credentials = username + ":" + password;
//		byte[] src = credentials.getBytes(Charset.forName("UTF-8"));
//		String basicAuthString = "Basic " + new String(Base64.getEncoder().encodeToString(src));
//		
//		headerMap.put("Authorization", basicAuthString);
//		headerMap.put("Accept", "application/json");
//		headerMap.put("Content-type", "application/json");
//		
//		given().headers(headerMap).when().post("http://localhost:9080/host/update").then()
//				.statusCode(HttpStatus.BAD_REQUEST.value());
//	}
//	
//	@Test
//	public void validateUpdateWithJSON() {
//		Map<String, String> headerMap = new HashMap<>();
//		String username = "user";
//		String password = "password";
//		String credentials = username + ":" + password;
//		byte[] src = credentials.getBytes(Charset.forName("UTF-8"));
//		String basicAuthString = "Basic " + new String(Base64.getEncoder().encodeToString(src));
//		
//		headerMap.put("Authorization", basicAuthString);
//		headerMap.put("Accept", "application/json");
//		headerMap.put("Content-type", "application/json");
//		
//		given().headers(headerMap).body(agentPost).when().post("http://localhost:9080/host/update").then()
//				.statusCode(HttpStatus.OK.value());
//	}
//	
//	@Test
//	public void validateUpdateWithNoHostname() {
//		Map<String, String> headerMap = new HashMap<>();
//		String username = "user";
//		String password = "password";
//		String credentials = username + ":" + password;
//		byte[] src = credentials.getBytes(Charset.forName("UTF-8"));
//		String basicAuthString = "Basic " + new String(Base64.getEncoder().encodeToString(src));
//		
//		headerMap.put("Authorization", basicAuthString);
//		headerMap.put("Accept", "application/json");
//		headerMap.put("Content-type", "application/json");
//		
//		given().headers(headerMap).body(agentPostNoHostname).when().post("http://localhost:9080/host/update").then()
//			.statusCode(HttpStatus.BAD_REQUEST.value());
//	} 
	

}
