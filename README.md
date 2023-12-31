## Intro

The DB-IP Lite free database is a subset of the full database with reduced coverage and accuracy.

Based on my search and experience, its way better than MaxMind's `GeoLite2 City` and `GeoLite2 Country` database quality in terms of coverage and accuracy.

Cons is the Lite downloads are updated monthly, way lower than MaxMind's update frequency.

However, DP-IP didn't provide API key for the Free (lite verison) IP geolocation database download, means you must visit the webpage and click the checkbox of `I agree with the licensing terms`

So, Wrote a Golang program to download it automaticly. 

## Downloadable assets:
[![Download DB-IP's file.](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml/badge.svg?branch=main)](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml)
- [Only on GitHub](https://github.com/MaurUppi/downloader/releases)
- Not avaliable on jsDelivr due to Package size exceeded the configured limit of 50 MB.


## License

[CC-BY-SA-4.0](https://creativecommons.org/licenses/by-sa/4.0/)

- This product includes **IP to Country Lite Database** data created by DB-IP.com, available from [IP Geolocation by DB-IP](https://db-ip.com)
- Distributed under the Creative Commons Attribution License.
