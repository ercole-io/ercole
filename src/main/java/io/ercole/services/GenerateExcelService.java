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

package io.ercole.services;

import io.ercole.model.CurrentHost;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.utilities.JsonFilter;

import org.apache.commons.io.output.ByteArrayOutputStream;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.xssf.usermodel.XSSFRow;
import org.apache.poi.xssf.usermodel.XSSFSheet;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.json.JSONArray;
import org.json.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

import java.io.IOException;

/**
 * The type Generate excel service.
 */
@Service
public class GenerateExcelService {

    private static String cpuM = "CPUModel";
    private static String vrs = "Version";

    @Autowired
    private CurrentHostRepository currentRepo;


    /**
     * Init excel response entity.
     *
     * @return the response entity
     * @throws IOException      the io exception
     * @throws RuntimeException the runtime exception
     */
    public ResponseEntity<byte[]> initExcel() throws IOException {
        Iterable<CurrentHost> iterable = currentRepo.findAllByOrderByHostnameAsc();


        try (Workbook workbook = new XSSFWorkbook(new ClassPathResource("template_hosts.xlsm").getInputStream())) {

            XSSFSheet xssfSheet = ((XSSFWorkbook) workbook).getSheet("Database_&_EBS");
            //number of row where we will write (row 0,1,2 contains the heading of the table)
            int rowNumber = 3;

            for (CurrentHost host : iterable) {
                XSSFRow row = xssfSheet.createRow(rowNumber);
                //get json HostInfo (it contains another 16 field)
                JSONObject root = new JSONObject(host.getHostInfo());
                //get json ExtraInfo (it contains another 2 json array - databases and features)
                JSONObject root2 = new JSONObject(host.getExtraInfo());
                //get json databases
                JSONArray arrayExtraInfo = root2.getJSONArray("Databases");
                //In the db there are some hosts that don't serve databases
                if (arrayExtraInfo.length() == 0) {
                    continue;
                }
                JSONObject database = arrayExtraInfo.getJSONObject(0);
                //get json array features (after we will create jsonObject features from this array -getFeatures())
                JSONArray features = database.getJSONArray("Features");

                //index 
                int indexChiocciola = root.getString(cpuM).lastIndexOf(' ');
                int indexVersion = database.getString(vrs).indexOf(' ');

                //save data into array
                String[] dataOfHost = new String[20];
                dataOfHost[0]  =
                        host.getHostname();                     //physical server name
                dataOfHost[1]  =
                        host.getAssociatedClusterName();        //virtual server name
                dataOfHost[2]  =
                        root.getString("Type");           //virtualization technol.
                dataOfHost[3]  =
                        host.getDatabases();                   //db instace name
                dataOfHost[4]  =
                        " ";                                   //pluggable db name
                dataOfHost[5]  =
                        " ";                                   //connect string
                dataOfHost[6]  =
                        String.valueOf(database.get(vrs)).substring(0, 2); //product version
                dataOfHost[7]  =
                        String.valueOf(database.get(vrs)).substring(indexVersion);	  //product edition
                dataOfHost[8]  =
                        host.getEnvironment();                    //environment
                dataOfHost[9]  =
                        JsonFilter.getTrueFeatures(features);       // options and oem packs
                dataOfHost[10] =
                        " ";         //rac node names
                dataOfHost[11] =
                        String.valueOf(root.get(cpuM));             //processor model
                dataOfHost[12] =
                        String.valueOf(root.get("Socket"));               //processor socket
                dataOfHost[13] =
                        String.valueOf(root.get("CPUCores")).replace("'", ""); //core per processor
                Integer s      =
                        (Integer.parseInt(dataOfHost[12]) * Integer.parseInt(dataOfHost[13]));
                dataOfHost[14] =
                        s.toString();         //physical core
                if (String.valueOf(root.get(cpuM)).contains("SPARC")) {
                    dataOfHost[15] = "8";
                } else {
                    dataOfHost[15] = "2";
                }
                dataOfHost[16] =
                        String.valueOf(root.get(cpuM)).substring(indexChiocciola); //processor speed
                dataOfHost[17] =
                        " ";                             //server purchases date
                dataOfHost[18] =
                        String.valueOf(root.get("OS"));  //operating system
                dataOfHost[19] =
                        " ";                             //note

                //insert data of host into a new cell
                int cellid = 0;
                for (int i = 0; i < 20; i++) {
                    //don't delete. see templateVuoto ABCD?F
                    if (cellid == 5) {
                        cellid++;
                    }
                    Cell cell = row.createCell(cellid++);
                    cell.setCellValue(dataOfHost[i]);
                }
                rowNumber++;
            }
            //writing changes in the open file (templateVuoto)

            try (ByteArrayOutputStream outputStream = new ByteArrayOutputStream()) {

                workbook.write(outputStream);
                HttpHeaders headers = new HttpHeaders();
                headers.add(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=Hosts.xlsm");
                return ResponseEntity.ok()
                        .headers(headers)
                        .contentType(MediaType.parseMediaType("application/vnd.ms-excel"))
                        .body(outputStream.toByteArray());
            }
        }
    }


    /**
     * Init excel response entity.
     *
     * @return the response entity
     * @throws IOException      the io exception
     * @throws RuntimeException the runtime exception
     */
    public ResponseEntity<byte[]> initExcelWithoutTemplate() throws IOException {
        Iterable<CurrentHost> iterable = currentRepo.findAllByOrderByHostnameAsc();


        try (Workbook workbook = new XSSFWorkbook()) {

            XSSFSheet xssfSheet = ((XSSFWorkbook) workbook).createSheet();
            //header row
            XSSFRow rowHeader = xssfSheet.createRow(0);
            rowHeader.createCell(0).setCellValue("hostname");
            rowHeader.createCell(1).setCellValue("env");
            rowHeader.createCell(2).setCellValue("host type");
            rowHeader.createCell(3).setCellValue("cluster");
            rowHeader.createCell(4).setCellValue("physical host");
            rowHeader.createCell(5).setCellValue("last update");
            rowHeader.createCell(6).setCellValue("databases");
            rowHeader.createCell(7).setCellValue("OS");
            rowHeader.createCell(8).setCellValue("kernel");
            rowHeader.createCell(9).setCellValue("oracle cluster");
            rowHeader.createCell(10).setCellValue("sun cluster");
            rowHeader.createCell(11).setCellValue("veritas cluster");
            rowHeader.createCell(12).setCellValue("virtual");
            rowHeader.createCell(13).setCellValue("host type");
            rowHeader.createCell(14).setCellValue("cpu threads");
            rowHeader.createCell(15).setCellValue("cpu cores");
            rowHeader.createCell(16).setCellValue("sockets");
            rowHeader.createCell(17).setCellValue("mem total");
            rowHeader.createCell(18).setCellValue("swap total");

            //data rows
            int rowNumber = 1;
            for (CurrentHost host : iterable) {
                XSSFRow row = xssfSheet.createRow(rowNumber);
                //get json HostInfo (it contains another 16 field)
                JSONObject root = new JSONObject(host.getHostInfo());
                // //get json ExtraInfo (it contains another 2 json array - databases and features)
                // JSONObject root2 = new JSONObject(host.getExtraInfo());
                // //get json databases
                // JSONArray arrayExtraInfo = root2.getJSONArray("Databases");
                //In the db there are some hosts that don't serve databases
                // if (arrayExtraInfo.length() == 0) {
                //     continue;
                // }
                // JSONObject database = arrayExtraInfo.getJSONObject(0);
                // //get json array features (after we will create jsonObject features from this array -getFeatures())
                // JSONArray features = database.getJSONArray("Features");

                //save data into array
                String[] dataOfHost = new String[20];
                dataOfHost[0]  = host.getHostname();                     
                dataOfHost[1]  = host.getEnvironment();
                dataOfHost[2]  = host.getHostType();
                dataOfHost[3]  = host.getAssociatedClusterName();                   
                dataOfHost[4]  = host.getAssociatedHypervisorHostname();
                dataOfHost[5]  = host.getUpdated().toString();
                dataOfHost[6]  = host.getDatabases();
                dataOfHost[7] =  root.getString("OS");
                dataOfHost[8]  = root.getString("Kernel");
                dataOfHost[9]  = "" + root.getBoolean("OracleCluster");
                dataOfHost[10]  = "" + root.getBoolean("SunCluster"); 
                dataOfHost[11] = "" + root.getBoolean("VeritasCluster");
                dataOfHost[12] = "" + root.getBoolean("Virtual");
                dataOfHost[13] = root.getString("Type"); 
                dataOfHost[14] = "" + root.getInt("CPUThreads");
                dataOfHost[15] = "" + root.getInt("CPUCores");
                dataOfHost[16] = "" + root.getInt("Socket"); 
                dataOfHost[17] = "" + root.getInt("MemoryTotal"); 
                dataOfHost[18] = "" + root.getInt("SwapTotal"); 
                dataOfHost[19] = root.getString("CPUModel");

                //insert data of host into a new cell
                for (int i = 0; i < dataOfHost.length; i++) {
                    Cell cell = row.createCell(i);
                    cell.setCellValue(dataOfHost[i]);
                }
                rowNumber++;
            }
            //writing changes in the open file (templateVuoto)

            try (ByteArrayOutputStream outputStream = new ByteArrayOutputStream()) {

                workbook.write(outputStream);
                HttpHeaders headers = new HttpHeaders();
                headers.add(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=HostsRaw.xlsx");
                return ResponseEntity.ok()
                        .headers(headers)
                        .contentType(MediaType.parseMediaType("application/vnd.ms-excel"))
                        .body(outputStream.toByteArray());
            }
        }
    }
}
