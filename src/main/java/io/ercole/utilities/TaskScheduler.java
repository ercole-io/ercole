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

package io.ercole.utilities;

import java.util.Date;

import javax.transaction.Transactional;

import org.apache.commons.lang3.time.DateUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import io.ercole.model.Alert;
import io.ercole.model.AlertCode;
import io.ercole.model.AlertSeverity;
import io.ercole.model.CurrentHost;
import io.ercole.model.HistoricalHost;
import io.ercole.repositories.AlertRepository;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.repositories.HistoricalHostRepository;

/**
 * Generic scheduler.
 */
@Component("ErcoleTaskScheduler")
public class TaskScheduler {
	private Logger logger = LoggerFactory.getLogger(TaskScheduler.class);
	
	@Value("${current.host.cleaning.hourRate}")
	private int currentHours;
	
	@Value("${historical.host.cleaning.hourRate}")
	private int historicalHours;
	
	@Value("${current.host.freshness.check.freshness.threshold}")
	private int freshnessThreshold;

	@Autowired
	private CurrentHostRepository currentRepository;
	
	@Autowired
	private HistoricalHostRepository historicalRepository;
	
	@Autowired
	private AlertRepository alertRepository;

	
	/**
	 * It cleans current_hosts moving the ones older than fixed "current.host.cleaning.hourRate"
	 * hours attribute to historical_hosts.
	 */
	@Scheduled(cron = "${current.host.cleaning.crontab}")
	@Transactional
	public void currentHostCleaning() {
		Date lowerLimit = DateUtils.addHours(new Date(), -currentHours);
		
		for (CurrentHost current : currentRepository.findAllNotUpdated(lowerLimit)) {
			HistoricalHost historical = new HistoricalHost(current.getId(), current.getHostname(),
					current.getEnvironment(), current.getLocation(), 
					current.getVersion(), current.getServerVersion(),
					current.getHostType(), current.getDatabases(),
					current.getSchemas(), current.getExtraInfo(), current.getAssociatedClusterName(), current.getAssociatedHypervisorHostname(), current.getHostInfo(),
					current.getUpdated());
			historicalRepository.save(historical);
			currentRepository.delete(current);
			
			if (!logger.isDebugEnabled()) {
				logger.info(String.format("%s has been moved because it have passed more than %s  "
					+ "hours from last update (%s)", current, currentHours, current.getUpdated()));
			}
		}
	}
	
	/**
	 * It cleans historical_hosts deleting the ones older than fixed 
	 * "historical.host.cleaning.hourRate" hours attribute.
	 */
	@Scheduled(cron = "${historical.host.cleaning.crontab}")
	public void historicalHostCleaning() {
		Date lowerLimit = DateUtils.addHours(new Date(), -historicalHours);
		
		for (HistoricalHost historical : historicalRepository.findAllOlderThan(lowerLimit)) {
			historicalRepository.delete(historical);
			
			if (!logger.isDebugEnabled()) {
				logger.info(String.format("%s has been deleted because it have passed more than %s  "
					+ "hours from last archive (%s)", historical, historicalHours, 
					historical.getArchived()));
			}
		}
	}
	
	/**
	 * It alert the user when current_hosts is older than fixed 
	 * "current.host.freshness.check.freshness.threshold" days attribute.
	 */
	@Scheduled(cron = "${current.host.freshness.check.crontab}")
	public void checkCurrentHostFreshness() {
		Date lowerLimit = DateUtils.addDays(new Date(), -freshnessThreshold);

		logger.info("Start checking freshness");
		for (CurrentHost c : currentRepository.findAllNotUpdated(lowerLimit)) {
			if (!alertRepository.existsByHostnameAndCode(c.getHostname(), AlertCode.NO_DATA)) {
				alertRepository.save(new Alert(c.getHostname(), AlertCode.NO_DATA, 
				"No data received from the host in the last " + freshnessThreshold + " days", AlertSeverity.MAJOR));
			}
		}
	}
	
}
