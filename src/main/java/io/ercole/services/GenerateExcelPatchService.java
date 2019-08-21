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
import java.util.Calendar;
import java.util.Date;
import java.util.List;
import java.util.Map;
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

import io.ercole.repositories.CurrentHostRepository;

/**
 * The type Generate excel service.
 */
@Service
public class GenerateExcelPatchService {

    @Autowired
    private CurrentHostRepository currentRepo;

    /**
     * Init excel response entity.
     *
     * @param windowTime the window time
     * @param status     the status
     * @return the response entity
     * @throws IOException the io exception
     */
    public ResponseEntity<byte[]> initExcel(final int windowTime, final String status) throws IOException {

        Calendar calendar = Calendar.getInstance();
        calendar.set(Calendar.HOUR_OF_DAY, 0);
        calendar.set(Calendar.MINUTE, 0);
        calendar.set(Calendar.DATE, 1);
        calendar.add(Calendar.MONTH, -windowTime);
        Date time = calendar.getTime();

        List<Map<String, Object>> patchList = currentRepo.getAllHostPSUStatus(time, status);

        try (Workbook workbook = new XSSFWorkbook(new ClassPathResource("template_patch_advisor.xlsm").getInputStream())) {

            XSSFSheet xssfSheet = ((XSSFWorkbook) workbook).getSheet("Patch_Advisor");
            //number of row where we will write (row 0,1,2 contains the heading of the table)
            int rowNumber = 3;




            for (Map<String, Object> patchMap : patchList) {
                XSSFRow row = xssfSheet.createRow(rowNumber);

                String[] dataOfHost = new String[6];

                dataOfHost[0]  = String.valueOf(patchMap.get("psudescription"));

                dataOfHost[1]  = String.valueOf(patchMap.get("hostname"));

                dataOfHost[2]  = String.valueOf(patchMap.get("dbname"));

                dataOfHost[3]  = String.valueOf(patchMap.get("dbver"));

                dataOfHost[4]  = String.valueOf(patchMap.get("psudate"));

                dataOfHost[5]  = String.valueOf(patchMap.get("status"));

                //insert data of host into a new cell
                int cellid = 0;
                for (int i = 0; i < 6; i++) {
                    //don't delete. see templateVuoto ABCD?F

                    Cell cell = row.createCell(cellid++);
                    cell.setCellValue(dataOfHost[i]);
                }
                rowNumber++;
            }
            //writing changes in the open file (templateVuoto)

            try (ByteArrayOutputStream outputStream = new ByteArrayOutputStream()) {

                HttpHeaders headers = new HttpHeaders();
                workbook.write(outputStream);
                headers.add(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=PatchAdvisor.xlsm");
                return ResponseEntity.ok()
                        .headers(headers)
                        .contentType(MediaType.parseMediaType("application/vnd.ms-excel"))
                        .body(outputStream.toByteArray());
            }
        }
    }
}
