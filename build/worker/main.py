from transit import request

# WIP
response = request.TimedHttpRequest.make("google.com")
stats = request.HttpRequestStatistics(response)
print(stats)