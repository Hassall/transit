from transit.request import TimedHttpRequest
from transit.request import HttpRequestStatistics


def test_request_string_contents():
    assert 'time_namelookup' in TimedHttpRequest.requestString
    assert 'time_connect' in TimedHttpRequest.requestString
    assert 'time_appconnect' in TimedHttpRequest.requestString
    assert 'time_pretransfer' in TimedHttpRequest.requestString
    assert 'time_starttransfer' in TimedHttpRequest.requestString
    assert 'time_total' in TimedHttpRequest.requestString


def test_http_request_stats():
    response = TimedHttpRequest.make("google.com")
    stats = HttpRequestStatistics(response)
    assert stats.time_namelookup >= 0
    assert stats.time_connect >= 0
    assert stats.time_appconnect >= 0
    assert stats.time_pretransfer >= 0
    assert stats.time_starttransfer >= 0
    assert stats.time_total >= 0
