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
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.mockito.Mockito.when;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

import io.ercole.repositories.AlertRepository;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.services.DashboardService;
import io.ercole.utilities.JsonFilter;

@RunWith(PowerMockRunner.class)
@PrepareForTest(JsonFilter.class)
public class DashboardServiceTest {

	@Mock
	private CurrentHostRepository currentRepo;
	
	@Mock
	private AlertRepository alertRepo;
	
	@InjectMocks
	private DashboardService dashService;
	
	@Test
	public void acknowledgeAlertTest() {
		Long[] array = {1l, 2l, 3l};
		
		for (Long id : array) {
			when(alertRepo.setFromNewToAck(id)).thenReturn(1);
		}
		assertEquals(true, dashService.acknowledgeAlert(array));
		
		
		for (Long id : array) {
			when(alertRepo.setFromNewToAck(id)).thenReturn(0);
		}
		assertEquals(false, dashService.acknowledgeAlert(array));
	}
	
	@Test
	public void getDbFeaturesCount0() {
		List<String> queryVal = new ArrayList<>();
		when(currentRepo.getDbFeaturesCount("*")).thenReturn(queryVal);
		List<String> emptyList = new ArrayList<>();
		
		List<Map<String, Object>> retVal = dashService.getDbFeaturesCount("*");
		assertEquals(emptyList, retVal);
	}
	
	
	
	
	@Test
	public void getDbFeaturesCountMoreThan0() {
		List<String> queryVal = new ArrayList<>();
		queryVal.add("Feature1");
		queryVal.add("Feature2");
		queryVal.add("Feature1");
				
		Map<String, Object> map1 = new HashMap<>();
		map1.put("feature", "Feature1");
		map1.put("value", Integer.valueOf(1));
		
		Map<String, Object> map2 = new HashMap<>();
		map2.put("feature", "Feature2");
		map2.put("value", Integer.valueOf(1));

		List<Map<String, Object>> retVal = new ArrayList<>();
		retVal.add(map1);
		retVal.add(map2);
		
		when(currentRepo.getDbFeaturesCount("*")).thenReturn(queryVal);
			
		List<Map<String, Object>> controlList = dashService.getDbFeaturesCount("*");
		assertEquals(retVal.size(), controlList.size());
		assertTrue(controlList.get(0).get("license").equals("Feature2")
			|| controlList.get(0).get("license").equals("Feature1"));
		assertTrue(controlList.get(1).get("license").equals("Feature2")
			|| controlList.get(1).get("license").equals("Feature1"));
		
	}
	
	
	@Test
	public void getServerEnvTest() {
		String location = "Italia";

		List<Map<String, Object>> lista = new ArrayList<>();
		when(currentRepo.getServerTypeCount(location)).thenReturn(lista);
		assertEquals(lista, dashService.getServerEnv(location));
	}
	
	@Test
	public void getLocationsTest() {
		List<String> lista = new ArrayList<>();
		
		when(currentRepo.getLocations()).thenReturn(lista);
		assertEquals(lista, dashService.getServerLocation());
		assertTrue(dashService.getServerLocation().isEmpty());
	}
	
	
	@Test
	public void getDbEnvTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getDbEnvs(location)).thenReturn(lista);
		assertEquals(lista, dashService.getDbEnv(location));
		assertTrue(dashService.getDbEnv(location).isEmpty());
	}
	
	@Test
	public void getDbVersionsCountTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getDbVersionsCount(location)).thenReturn(lista);
		assertEquals(lista, dashService.getDbVersionsCount(location));
		assertTrue(dashService.getDbVersionsCount(location).isEmpty());
	}
	
	@Test
	public void getHostTypeCountTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getHostTypeCount(location)).thenReturn(lista);
		assertEquals(lista, dashService.getHostTypeCount(location));
		assertTrue(dashService.getHostTypeCount(location).isEmpty());
	}
	
	
	@Test
	public void getOsTypeCountTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getOsTypeCount(location)).thenReturn(lista);
		assertNotNull(lista);
		assertNotNull(dashService);
		assertNotNull(location);
		assertNotNull(dashService.getOsTypeCount(location));
		//assertEquals(lista, dashService.getOsTypeCount(location));
		assertTrue(dashService.getOsTypeCount(location).isEmpty());
	}

	
	@Test
	public void getResizeTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getResizeByHosts(location)).thenReturn(lista);
		assertEquals(lista, dashService.getResize(location));
		assertTrue(dashService.getResize(location).isEmpty());
	}
	
	
	@Test
	public void getWorkByDbsTest() {
		List<Map<String, Object>> lista = new ArrayList<>();
		String location = "Italia";
		
		when(currentRepo.getWorkByDbs(location)).thenReturn(lista);
		assertEquals(lista, dashService.getWorkByDbs(location));
		assertTrue(dashService.getWorkByDbs(location).isEmpty());
	}
	
	
	@Test
	public void getLicensesCount() {
		String location = "Italia";

		List<Map<BigDecimal, String>> lista = new ArrayList<>();
		when(currentRepo.getLicensesCount(location)).thenReturn(lista);
		assertEquals(lista, dashService.getLicensesCount(location));
		
		Map<BigDecimal, String> mappa = new HashMap<>();
		lista.add(mappa);
		assertEquals(lista, dashService.getLicensesCount(location));
	}
}
