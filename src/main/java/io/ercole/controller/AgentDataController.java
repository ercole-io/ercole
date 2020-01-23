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

package io.ercole.controller;

import java.io.IOException;
import java.nio.charset.Charset;
import java.text.ParseException;
import java.util.Base64;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.time.DateUtils;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import io.ercole.exceptions.AgentFloodException;
import io.ercole.exceptions.AgentLoginException;
import io.ercole.exceptions.HostNotFoundException;
import io.ercole.exceptions.NoHistoryFoundException;
import io.ercole.model.HistoricalHost;
import io.ercole.services.HostService;
import io.ercole.utilities.DateUtility;

/**
 * RestController for CRUD actions.
 */
@RestController
public class AgentDataController {
	private Logger logger = LoggerFactory.getLogger(AgentDataController.class);
	private static final String HT = "HostType";

	@Autowired
	private HostService hostService;

	@Value("${agent.user}")
	private String agentUser;

	@Value("${agent.password}")
	private String agentPassword;

	/**
	 * @param request
	 *            used to parse incoming JSON
	 * @throws HostNotFoundException
	 *             if you're updating a non-existent Host
	 * @throws AgentLoginException
	 *             if "user" and "password" header attributes aren't set in the
	 *             right manner
	 * @throws IOException 
	 * @throws AgentFloodException
	 *             if the agent tryes to update a hostname during the wrong time
	 *             range
	 */
	@PostMapping(value = "${agent.api.update}", consumes = "application/json")
	public void updateFromAgent(final HttpServletRequest request)
			throws IOException, HostNotFoundException, AgentLoginException, AgentFloodException {
		if (areAgentCredentialsValid(request.getHeader("Authorization"))) {
			try {
				String body = IOUtils.toString(request.getReader());
				String hostType;
				if (request.getParameterMap().containsKey(HT)) {
					hostType = request.getParameter(HT);
				} else {
					hostType = "oracledb";
				}

				JSONObject object = new JSONObject(body);

				if (object.isNull(HT)) {
					object.put(HT, "oracledb");
				}
				if (object.getString("Hostname").equals("")) {
					throw new HostNotFoundException("E' necessario un hostname per eseguire " 
							+ "l'update.");
				}

				String statusCode = hostService.updateWithAgent(object, hostType);

				if (statusCode.equals("ERROR")) {
					throw new AgentFloodException(object.getString("Hostname") 
							+ " non si può aggiornare attualmente.");
				}
			} catch (IOException e) {
				logger.error("Error receiving data", e);
			} catch (ParseException e) {
				logger.error(e.getMessage() + " feature presenta dei problemi: "
						+ "controllare se in tutti i DB del JSON "
						+ "sia scritta nello stesso modo!");
			}
		} else {
			throw new AgentLoginException("Parametri user o password dell'agente errati.");
		}
	}

	private boolean areAgentCredentialsValid(final String authorization) {
		if (authorization != null && authorization.startsWith("Basic")) {
			// Authorization: Basic base64credentials
			String base64Credentials = authorization.substring("Basic".length()).trim();
			String credentials = new String(Base64.getDecoder().decode(base64Credentials),
					Charset.forName("UTF-8"));
			// credentials = username:password
			final String[] values = credentials.split(":", 2);
			return (values[0].equals(this.agentUser) && values[1].equals(this.agentPassword));
		} else {
			return false;
		}
	}

	/**
	 * @param hostname
	 *            to search from HistoricalHost (history logs)
	 * @param data
	 *            date parameter for the search query
	 * @return last hostname log on the specified date
	 * @throws ParseException 
	 * @throws NoHistoryFoundException  
	 */
	@GetMapping(value = "/historical")
	public HistoricalHost getHostnameHistory(@RequestParam("hostname") final String hostname,
			@RequestParam("data") final String data) throws ParseException, NoHistoryFoundException {
		HistoricalHost historical = hostService.getHistoricalLogs(hostname, 
				DateUtils.addDays(DateUtility.fromStringToDate(data), 1));
		if (historical == null) {
			throw new NoHistoryFoundException("Non è presente alcun storico del hostname " 
					+ hostname + " e data " + data);
		} else {
			return historical;
		}
	}

	/**
	 * @param exception
	 *            thrown
	 * @param response
	 *            to whom we send the exception message
	 * @throws IOException 
	 */
	@ExceptionHandler(HostNotFoundException.class)
	public void handleHostNotFound(final HostNotFoundException exception, final HttpServletResponse response)
			throws IOException {
		response.sendError(HttpStatus.BAD_REQUEST.value(), exception.getMessage());
		logger.error(exception.getMessage());
	}

	/**
	 * @param exception
	 *            thrown
	 * @param response
	 *            to whom we send the exception message
	 * @throws IOException 
	 */
	@ExceptionHandler(JSONException.class)
	public void handleJsonParsing(final JSONException exception, final HttpServletResponse response)
			throws IOException {
		response.sendError(HttpStatus.BAD_REQUEST.value(), exception.getMessage());
		logger.error("Error parsing data", exception);
	}

	/**
	 * @param exception
	 *            thrown
	 * @param response
	 *            to whom we send the exception message
	 * @throws IOException 
	 */
	@ExceptionHandler(AgentLoginException.class)
	public void handleAgentLogin(final AgentLoginException exception, final HttpServletResponse response)
			throws IOException {
		response.sendError(HttpStatus.UNAUTHORIZED.value(), exception.getMessage());
		logger.error(exception.getMessage());
	}

	/**
	 * @param exception
	 *            thrown
	 * @param response
	 *            to whom we send the exception message
	 * @throws IOException 
	 */
	@ExceptionHandler(AgentFloodException.class)
	public void handleAgentFlood(final AgentFloodException exception, final HttpServletResponse response)
			throws IOException {
		response.sendError(HttpStatus.TOO_MANY_REQUESTS.value(), exception.getMessage());
		logger.error(exception.getMessage());
	}
	
	/**
	 * @param exception
	 *            thrown
	 * @param response
	 *            to whom we send the exception message
	 * @throws IOException 
	 */
	@ExceptionHandler(NoHistoryFoundException.class)
	public void handleHistoryNotFound(final NoHistoryFoundException exception, final HttpServletResponse response)
			throws IOException {
		response.sendError(HttpStatus.NO_CONTENT.value(), exception.getMessage());
		logger.info(exception.getMessage());
	}

	@PostMapping("/alerts/missing-host/{hostname}")
	public void checkHostAbsence(final HttpServletRequest request, final @PathVariable String hostname) throws AgentLoginException {
		if (areAgentCredentialsValid(request.getHeader("Authorization"))) {
			hostService.checkHostAbsence(hostname);
		} else {
			throw new AgentLoginException("Parametri user o password dell'agente errati.");
		}
	}
}
