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

import java.io.IOException;
import java.util.HashSet;
import java.util.List;
import org.apache.commons.io.output.ByteArrayOutputStream;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.xssf.usermodel.XSSFRow;
import org.apache.poi.xssf.usermodel.XSSFSheet;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

import io.ercole.model.ClusterInfo;
import io.ercole.model.VMInfo;
import io.ercole.repositories.ClusterRepository;

/**
 * The type Generate excel service.
 */
@Service
public class GenerateExcelClustersService {

    @Autowired
    private ClusterRepository clusterRepo;

    /**
     * Init excel response entity.
     *
     * @param search the search
     * @return the response entity
     * @throws IOException the io exception
     */
    public ResponseEntity<byte[]> initExcel(final String search) throws IOException {

        List<ClusterInfo> clustersList = clusterRepo.getClusters(search);

        try (Workbook workbook = new XSSFWorkbook(new ClassPathResource("template_clusters.xlsx").getInputStream())) {

            XSSFSheet xssfSheet = ((XSSFWorkbook) workbook).getSheet("Hypervisor");
            //number of row where we will write (row 0,1,2 contains the heading of the table)
            int rowNumber = 1;


            for (ClusterInfo cl : clustersList) {
                XSSFRow row = xssfSheet.createRow(rowNumber);
                HashSet<String> physicalHosts = new HashSet<>();

                cl.getVms().forEach((VMInfo vm) -> {
                    physicalHosts.add(vm.getPhysicalHost());
                });

                String[] dataOfHost = new String[5];

                dataOfHost[0]  = String.valueOf(cl.getName());

                dataOfHost[1]  = String.valueOf(cl.getType());

                dataOfHost[2]  = String.valueOf(cl.getCpu());

                dataOfHost[3]  = String.valueOf(cl.getSockets());

                dataOfHost[4]  = String.valueOf(String.join(" ", physicalHosts));

                //insert data of host into a new cell
                int cellid = 0;
                for (int i = 0; i < dataOfHost.length; i++) {
                    //don't delete. see templateVuoto ABCD?F

                    Cell cell = row.createCell(cellid++);
                    cell.setCellValue(dataOfHost[i]);

                }
                rowNumber++;
            }
            //writing changes in the open file (templateVuoto)

            try (ByteArrayOutputStream outputStream2 = new ByteArrayOutputStream()) {

                HttpHeaders headers = new HttpHeaders();
                workbook.write(outputStream2);
                headers.add(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=ADDM.xlsx");
                return ResponseEntity.ok()
                        .headers(headers)
                        .contentType(MediaType.parseMediaType("application/vnd.ms-excel"))
                        .body(outputStream2.toByteArray());
            }
        }
    }
}
