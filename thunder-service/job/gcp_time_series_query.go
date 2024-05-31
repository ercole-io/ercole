// Copyright (c) 2024 Sorint.lab S.p.A.
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
package job

import (
	"context"
	"errors"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func (job *GcpDataRetrieveJob) GetTimeSeries(ctx context.Context, opt option.ClientOption, req *monitoringpb.ListTimeSeriesRequest) (*monitoringpb.TimeSeries, error) {
	c, err := monitoring.NewMetricClient(ctx, opt)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	if req == nil {
		return nil, errors.New("nil request, cannot retrieve monitoring time series")
	}

	var res *monitoringpb.TimeSeries

	it := c.ListTimeSeries(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		res = resp
	}

	return res, nil
}
