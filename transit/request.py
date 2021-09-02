import subprocess

class TimedHttpRequest:
    requestString = (
        "time_namelookup=%{time_namelookup},"
        "time_connect=%{time_connect},"
        "time_appconnect=%{time_appconnect},"
        "time_pretransfer=%{time_pretransfer},"
        "time_starttransfer=%{time_starttransfer},"
        "time_total=%{time_total}"
    )

    def make(url):
        process = subprocess.run(["curl", "-w", TimedHttpRequest.requestString,
                                  "-o", "/dev/null",
                                  "-s", url],
                                 capture_output=True)
        # stdout is in bytes, convert to string
        return process.stdout.decode("ascii")


class HttpRequestStatistics:
    def __init__(self, statistics):
        stats = dict([stat[0], stat[1]]
                     for stat in HttpRequestStatistics.extractStatistics(statistics))
        self.time_namelookup = stats["time_namelookup"]
        self.time_connect = stats["time_connect"]
        self.time_appconnect = stats["time_appconnect"]
        self.time_pretransfer = stats["time_pretransfer"]
        self.time_starttransfer = stats["time_starttransfer"]
        self.time_total = stats["time_total"]

    # takes a string and creates a 2d array [[statName, statValue],...]
    # assumes format is "key=value,..."
    def extractStatistics(statistics):
        return [[pair[0], float(pair[1])] for pair in (stat.split("=") for stat in statistics.split(","))]
