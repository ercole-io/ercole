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
}
