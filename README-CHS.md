# Downloader for DB-IP.com Databases

## [English README.md](https://github.com/MaurUppi/downloader/edit/main/README.md)

## 简介
本项目提供了一个自动化工具，专门用于从 [DB-IP.com](https://db-ip.com) 下载IP地理位置数据库。该工具主要通过GitHub Actions自动运行，确保数据库的持续更新和数据完整性。

## 背景

- DB-IP免费的Lite数据库是完整数据库的一个子集，其覆盖范围和准确度有所降低。
- 根据我的搜索和经验，就覆盖范围和准确度而言，它比MaxMind的`GeoLite2 City`和`GeoLite2 Country`数据库质量要好得多。
- 缺点是Lite版的下载更新频率为每月一次，远低于MaxMind的更新频率。
- 然而，DB-IP没有为免费（Lite版本）的IP地理位置数据库下载提供API密钥，这意味着您必须亲自访问网页并勾选`我同意许可条款`。
- 因此，我编写了一个Golang程序来自动完成下载。

## 主要功能和操作流程
- **自动化下载**：程序自动访问DB-IP.com，找到最新的IP地理位置数据库文件。涵盖如下文件
  - IP to City Lite CSV format
  - IP to City Lite MMDB format
  - IP to Country Lite CSV format
  - IP to Country Lite MMDB format
  - IP to ASN CSV format
  - IP to ASN MMDB format 
- **文件解压和校验**：下载的`.gz`格式文件会被自动解压，解压后的文件将进行SHA1校验，以确保下载内容的准确性和完整性。
  - 校验解压后的文件SHA1SUM值 对比 网页上提供的SHA1SUM 值
  - 校验通过后，将六个`.gz`文件发布到release
- **日志记录**：所有操作步骤和结果都会被记录在日志文件中，以便于问题追踪和操作审计。
  - 目前日志包含`DownloadLink`，`webSHA1SUM`，`confirmation msg`
  - 日志样本信息
    <details>
      <summary>点击展开查看详情</summary>
      
        DownloadLink: https://download.db-ip.com/free/dbip-asn-lite-2023-12.csv.gz
        webSHA1SUM: 3ef88d64af8d52def008c57a91df32ba5e4fe38a
        DownloadLink: https://download.db-ip.com/free/dbip-asn-lite-2023-12.mmdb.gz
        webSHA1SUM: cb874eb996813d3ac911755e8ff5e6d138e56541
        dbip-asn-lite-2023-12.csv.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value
        dbip-asn-lite-2023-12.mmdb.gz had been decompressed and SHA1SUM matched with webpage's SHA1SUM value    
    </details>


- **GitHub Actions自动化**：所有上述过程均通过GitHub Actions自动执行，不需要手动干预。
  - 提供了详细的`fmt.printf`输出到终端的信息，有兴趣可以围观
  - 由于Lite版本数据库是`月度`更新频率，因此Action是每个月的第二天凌晨运行 `cron: "0 0 2 * *"`

## GitHub Actions
项目已配置GitHub Actions，自动执行以下任务：
- 定期检查DB-IP.com网站以获取最新数据库文件。
- 自动下载、解压和校验文件。
- 日志记录操作过程和结果。
- 运行日志样本
      <details>
      <summary>点击展开查看详情</summary>
      
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

## 依赖
- Go version 1.21.4
- [goquery](https://github.com/PuerkitoBio/goquery): 用于解析HTML文档。
  - 根据标签名定位到文件链接以及对应的SHA1SUM
  - 执行网页的JavaScript确保checkbox勾选逻辑运行
    
- [chromedp](https://github.com/chromedp/chromedp): 用于Chrome浏览器的自动化操作。
  - 使用无头（headless）方式运行，`Chrome 114.0.5735.133 LTS`
  - chromedp `v0.9.2` 好像不兼容 Chrome 最新的 `120`，表现是死活都将将文件下载到默认目录（Windows 11环境）

## 输出
- **下载文件**：自动保存在GitHub Actions运行环境的`output`目录，并最终推送到Release。
- **日志文件**：操作日志`Downloaded_files.log`详细记录了下载和校验过程。

## 下载方式:
[![Download DB-IP's file.](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml/badge.svg?branch=main)](https://github.com/MaurUppi/downloader/actions/workflows/downlaoder.yml)
- [目前仅提供Github release](https://github.com/MaurUppi/downloader/releases)
- 由于超过了jsDelivr对于整包50M的限制，所以无法提供CDN下载，除非有人赞助其它CDN存储。


## License
[CC-BY-SA-4.0](https://creativecommons.org/licenses/by-sa/4.0/)
- This product includes **IP to ASN/Country/City Lite Database** data created by DB-IP.com, available from [IP Geolocation by DB-IP](https://db-ip.com)
- Distributed under the Creative Commons Attribution License.
