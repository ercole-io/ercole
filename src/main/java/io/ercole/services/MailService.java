package io.ercole.services;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.MailException;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

import io.ercole.model.Alert;
import io.ercole.model.AlertSeverity;

@Service
public class MailService {

	@Autowired
	private JavaMailSender sender;

	@Value("${alert.mail.from}")
	private String mailFrom;

	@Value("${alert.mail.to}")
	private String mailTo;

	@Value("${alert.mail.enabled}")
	private String mailEnabled;

	@Value("${alert.mail.severity}")
	private String mailSeverity;

	
	private Logger logger = LoggerFactory.getLogger(HostService.class);

	public static final int THREADS = 20;

	private ScheduledExecutorService executorService = Executors.newScheduledThreadPool(THREADS);

	public MailService() {
	}

	public void send(final Alert alert) throws MailException, RuntimeException {
		
		logger.debug("Sending email...");
		
		if (!Boolean.valueOf(mailEnabled)) { // mail alerts disabled
			logger.debug("Mail is disabled, skipping.");
			return;
		}
		
		if (!isWithinSeverity(alert)) { // lower severity than the one specified in settings
			logger.debug("Severity below threshold, skipping.");
			return;
		}
		
		SimpleMailMessage message = new SimpleMailMessage();
		message.setFrom(mailFrom);
		message.setTo(mailTo.split(","));
		message.setSubject(alert.getSeverity().name() + " "
				+ alert.getDescription() + " on " + alert.getHostname());
		message.setText("Date: " + alert.getDate() + "\nSeverity: " + alert.getSeverity() + "\nHost: "
				+ alert.getHostname() + "\nCode: " + alert.getCode() + "\n" + alert.getDescription());
		executorService.submit(new Runnable() {
			@Override
			public void run() {
				try {
					sender.send(message);
				} catch (Exception e) {
					logger.error("Error sending email", e);
				}
			}
		});
	}
	
	// Return true if the alert is within the severity threshold.
	// Do not rely on ordinal()
	private boolean isWithinSeverity(final Alert alert) {
		
		Map<AlertSeverity, Integer> severities = new HashMap<>();
		severities.put(AlertSeverity.NOTICE, 1);
		severities.put(AlertSeverity.MINOR, 2);
		severities.put(AlertSeverity.WARNING, 3);
		severities.put(AlertSeverity.MAJOR, 4);
		severities.put(AlertSeverity.CRITICAL, 5);
		
		int severity = severities.get(alert.getSeverity());
		int threshold = severities.get(AlertSeverity.valueOf(mailSeverity));
		
		return severity >= threshold;
	}
	
	
}
