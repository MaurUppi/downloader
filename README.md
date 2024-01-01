## [中文版介绍](https://github.com/MaurUppi/downloader/blob/main/README-CHS.md)

# Downloader for DB-IP.com Databases

## Introduction
This project provides an automated tool specifically for downloading IP geolocation databases from [DB-IP.com](https://db-ip.com). This tool primarily runs automatically through GitHub Actions, ensuring continuous updates and data integrity of the database.

## Background

- The DB-IP Lite free database is a subset of the full database with reduced coverage and accuracy.
- Based on my search and experience, its way better than MaxMind's `GeoLite2 City` and `GeoLite2 Country` database quality in terms of coverage and accuracy.
- Cons is the Lite downloads are updated monthly, way lower than MaxMind's update frequency.
- However, DP-IP didn't provide API key for the Free (lite verison) IP geolocation database download, means you must visit the webpage and click the checkbox of `I agree with the licensing terms`
- So, Wrote a Golang program to download it automaticly.

## Main Features and Operating Process
- **Automated Download**: The program automatically accesses DB-IP.com to find the latest IP geolocation database files, including:
  - IP to City Lite CSV format
  - IP to City Lite MMDB format
  - IP to Country Lite CSV format
  - IP to Country Lite MMDB format
  - IP to ASN CSV format
  - IP to ASN MMDB format 
- **File Decompression and Verification**: The downloaded `.gz` files are automatically decompressed, and the decompressed files undergo SHA1 verification to ensure accuracy and completeness of the downloaded content.
  - Verifies the SHA1SUM values of the decompressed files against the SHA1SUM values provided on the webpage.
  - After verification, the six `.gz` files are published to the release.
- **Logging**: All operation steps and results are recorded in a log file for troubleshooting and operational audit.
  - The current log includes `DownloadLink`, `webSHA1SUM`, `confirmation msg`.
  - Sample log information
    <details>
      <summary>Click to expand for details</summary>
      
        DownloadLink: https://download.db-ip.com/free/dbip-asn-lite-2023-12.csv.gz
        webSHA1SUM: 3ef88d64af8d52def008c57a91df32ba5e4fe38a
        DownloadLink: https://download.db-ip.com/free/dbip-asn-lite-2023-12.mmdb.gz
        webSHA1SUM: cb874eb996813d3ac911755e8ff5e6d138e56541
        dbip-asn-lite-2023-12.csv.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
        dbip-asn-lite-2023-12.mmdb.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value   
      
    </details>

- **GitHub Actions Automation**: All the above processes are automatically executed through GitHub Actions, without the need for manual intervention.
  - Detailed `fmt.printf` output to the terminal is provided, and you can observe if interested.
  - Since the Lite version of the database is updated monthly, the Action runs on the second day of each month at midnight `cron: "0 0 2 * *"`.

## GitHub Actions
The project is configured with GitHub Actions to automatically perform the following tasks:
- Regularly checks the DB-IP.com website for the latest database files.
- Automatically downloads, decompresses, and verifies files.
- Records the operation process and results.
- Utilizec GitHub LFS to push the File larger than 50M.
- Running log sample
      <details>
      <summary>Click to expand for details</summary>
      
      Chrome path is : /opt/hostedtoolcache/chromium/114.0.5735.133/x64/chrome
      Working dir is : /home/runner/work/downloader/downloader
      ouput dir create : /home/runner/work/downloader/downloader/output
      chromedp allocator context created
      URL: https://db-ip.com/db/download/ip-to-asn-lite
      File Type: .csv.gz
      Download Link: https://download.db-ip.com/free/dbip-asn-lite-2023-12.csv.gz
      SHA1SUM: 3ef88d64af8d52def008c57a91df32ba5e4fe38a
      URL: https://db-ip.com/db/download/ip-to-asn-lite
      File Type: .mmdb.gz
      Download Link: https://download.db-ip.com/free/dbip-asn-lite-2023-12.mmdb.gz
      SHA1SUM: cb874eb996813d3ac911755e8ff5e6d138e56541
      License agreement visible
      Checked checkbox
      Download link visible
      Clicked mmdb file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      CSV Download link visible
      Clicked CSV file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      Processing file: /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.csv.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.csv.gz to /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.csv
      dbip-asn-lite-2023-12.csv.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
      Processing file: /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.mmdb.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.mmdb.gz to /home/runner/work/downloader/downloader/output/dbip-asn-lite-2023-12.mmdb
      dbip-asn-lite-2023-12.mmdb.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
      URL: https://db-ip.com/db/download/ip-to-country-lite
      File Type: .mmdb.gz
      Download Link: https://download.db-ip.com/free/dbip-country-lite-2023-12.mmdb.gz
      SHA1SUM: a14ed000e7eea06b409dc34a2a6572babf3ef921
      URL: https://db-ip.com/db/download/ip-to-country-lite
      File Type: .csv.gz
      Download Link: https://download.db-ip.com/free/dbip-country-lite-2023-12.csv.gz
      SHA1SUM: fc5b4422ac7a8a52b336509d4f344c5052fe1825
      License agreement visible
      Checked checkbox
      Download link visible
      Clicked mmdb file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      CSV Download link visible
      Clicked CSV file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      Processing file: /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.csv.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.csv.gz to /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.csv
      dbip-country-lite-2023-12.csv.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
      Processing file: /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.mmdb.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.mmdb.gz to /home/runner/work/downloader/downloader/output/dbip-country-lite-2023-12.mmdb
      dbip-country-lite-2023-12.mmdb.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
      URL: https://db-ip.com/db/download/ip-to-city-lite
      File Type: .csv.gz
      Download Link: https://download.db-ip.com/free/dbip-city-lite-2023-12.csv.gz
      SHA1SUM: e93d44a611ee181c04cdec360432d6c196a3bc0b
      URL: https://db-ip.com/db/download/ip-to-city-lite
      File Type: .mmdb.gz
      Download Link: https://download.db-ip.com/free/dbip-city-lite-2023-12.mmdb.gz
      SHA1SUM: e1a6ab58d7858b5e8cec9c6722c5f52d0db99092
      License agreement visible
      Checked checkbox
      Download link visible
      Clicked mmdb file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      CSV Download link visible
      Clicked CSV file download link
      下载进度：0.00%
      下载进度：0.00%
      下载进度：60.76%
      下载进度：100.00%
      下载进度：100.00%
      下载进度：100.00%
      Processing file: /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.csv.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.csv.gz to /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.csv
      dbip-city-lite-2023-12.csv.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
      Processing file: /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.mmdb.gz
      Decompressing file: /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.mmdb.gz to /home/runner/work/downloader/downloader/output/dbip-city-lite-2023-12.mmdb
      dbip-city-lite-2023-12.mmdb.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value    
      
    </details>

## Dependencies
- Go version 1.21.4
- [goquery](https://github.com/PuerkitoBio/goquery): Used for parsing HTML documents.
  - Locates file links and corresponding SHA1SUM values based on tag names.
  - Executes webpage JavaScript to ensure checkbox ticking logic runs.
    
- [chromedp](https://github.com/chromedp/chromedp): Used for automating operations in the Chrome browser.
  - Runs in headless mode, `Chrome 114.0.5735.133 LTS`.
  - chromedp `v0.9.2` seems incompatible with the latest Chrome `120`, as observed in Windows 11 environment where files are always downloaded to the default directory.

## Output
- **Downloaded Files**: Automatically saved in the `output` directory of the GitHub Actions running environment and eventually pushed to Release.
- **Log File**: The `Downloaded_files.log` log file records the download and verification process in detail.



## Downloadable assets:
[![Download DB-IP's file.](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml/badge.svg?branch=main)](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml)
- [Currently only available via Github release](https://github.com/MaurUppi/downloader/releases)
- Unable to provide CDN download due to exceeding jsDelivr's 50M limit for the entire package, unless someone sponsors alternative CDN storage.


## License

[CC-BY-SA-4.0](https://creativecommons.org/licenses/by-sa/4.0/)

- This product includes **IP to Country Lite Database** data created by DB-IP.com, available from [IP Geolocation by DB-IP](https://db-ip.com)
- Distributed under the Creative Commons Attribution License.
