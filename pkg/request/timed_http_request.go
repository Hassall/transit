package request

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// URLRequest contains url to request statistics on
type URLRequest struct {
	URL string
}

var requestProperties string = `%{time_namelookup}
%{time_connect}
%{time_appconnect}
%{time_pretransfer}
%{time_starttransfer}
%{time_total}`

var requestPropertiesLength int = 6

// RequestStatistic represents the internal stages of a http request
type RequestStatistic struct {
	namelookup    float64
	connect       float64
	appconnect    float64
	pretransfer   float64
	starttransfer float64
	total         float64
}

func (requestStatistic RequestStatistic) String() string {
	return fmt.Sprintf(`Namelookup=%f
connect=%f
appconnect=%f
pretransfer=%f
starttransfer=%f
total=%f`,
		requestStatistic.namelookup,
		requestStatistic.connect,
		requestStatistic.appconnect,
		requestStatistic.pretransfer,
		requestStatistic.starttransfer,
		requestStatistic.total)
}

// TimedHTTPRequest executes http request via curl and returns breakdown of
// http stages.
func TimedHTTPRequest(url string) RequestStatistic {
	cmd, err := exec.Command("curl",
		"-w", requestProperties,
		"-o", "/dev/null",
		"-s", url).Output()

	if err != nil {
		fmt.Println(err)
		return RequestStatistic{}
	}

	return stringToRequestStatistic(string(cmd))
}

func stringToRequestStatistic(stats string) RequestStatistic {
	data := make([]float64, requestPropertiesLength)
	for i, stat := range strings.Split(stats, "\n") {
		if f, err := strconv.ParseFloat(stat, 64); err == nil {
			data[i] = f
		}
	}

	return RequestStatistic{data[0], data[1], data[2], data[3], data[4], data[5]}
}
