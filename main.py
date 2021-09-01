import subprocess

def main():
    url = "http://google.com/"

    transferStats = HttpRequestStatistics(TimedHttpRequest.make(url))

class TimedHttpRequest:
    requestString = "time_namelookup=%{time_namelookup} \
        time_connect=%{time_connect} \
        time_appconnect=%{time_appconnect} \
        time_pretransfer=%{time_pretransfer} \
        time_starttransfer=%{time_starttransfer} \
        time_total=%{time_total}"

    def make(url):
        process = subprocess.run(["curl", "-w", TimedHttpRequest.requestString, 
                                  "-o", "/dev/null", 
                                  "-s", url],
                                  capture_output=True)
        # stdout is in bytes, convert to string
        return process.stdout.decode("ascii")

class HttpRequestStatistics:
        def __init__(self, statistics):
            statsMap = dict([stat[0], stat[1]] for stat in HttpRequestStatistics.extractStatistics(statistics))
            self.time_namelookup = statsMap["time_namelookup"]       # DNS
            self.time_connect = statsMap["time_connect"]             # TCP Handshake
            self.time_appconnect = statsMap["time_appconnect"]       # SSL Handshake
            self.time_pretransfer = statsMap["time_pretransfer"]     # HTTP Request Start
            self.time_starttransfer = statsMap["time_starttransfer"] # TTFB
            self.time_total = statsMap["time_total"]                 # total time

        # takes a string and creates a 2d array [[statName, statValue],...]
        # assumes format is "key=value,..."
        def extractStatistics(statistic):
            return [[pair[0], pair[1]] for pair in (stat.split("=") for stat in statistic.split())]

if __name__ == "__main__":
    main()