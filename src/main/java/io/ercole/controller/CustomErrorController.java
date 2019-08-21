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

import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.web.servlet.error.ErrorAttributes;
import org.springframework.boot.web.servlet.error.ErrorController;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.context.request.WebRequest;

/**
 * Custom Controller for error handling.
 */
@RestController
public class CustomErrorController implements ErrorController {

	@Value("${server.error.path}")
	private static String errorPath;

	private final ErrorAttributes errorAttributes;

	/**
	 * @param errorAttr 
	 */
	@Autowired
	public CustomErrorController(final ErrorAttributes errorAttr) {
		this.errorAttributes = errorAttr;
	}

	/**
	 * @param request from front-end
	 * @return error attributes in json style
	 */
	@RequestMapping("${server.error.path}")
	public Map<String, Object> error(final WebRequest request) {
		return getErrorAttributes(request, false);
	}

	private Map<String, Object> getErrorAttributes(final WebRequest request, final boolean includeStacktrace) {
		return this.errorAttributes.getErrorAttributes(request, includeStacktrace);
	}

	@Override
	public final String getErrorPath() {
		return errorPath;
	}

}
