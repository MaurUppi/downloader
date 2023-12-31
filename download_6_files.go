package main

import (
	"compress/gzip"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

func main() {
	// Edge的可执行文件路径
	//browserPath := `C:\Program Files\Google\Chrome\Application\chrome.exe`
	// 使用 Action browser-actions/setup-chrome指定chrome版本以及 Outputs 变量
	browserPath := os.Getenv("CHROME_PATH") // 由 Github Action "Set Chrome Path" 步骤赋值
	fmt.Printf("Chrome path is : %s\n", browserPath)

	// 获取当前工作目录，并在此基础上创建output目录保存下载文件
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Working dir is : %s\n", wd)

	// 创建 output 子目录
	outputDir := filepath.Join(wd, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	fmt.Printf("ouput dir create : %s\n", outputDir)

	// 打开或创建日志文件，使用 os.O_TRUNC 来覆盖旧内容
	logFile, err := os.OpenFile(filepath.Join(outputDir, "Downloaded_files.log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 定义要处理的 URL 列表，解析页面获取下载链接URL、文件名和SHA1SUM
	urls := []string{
		"https://db-ip.com/db/download/ip-to-asn-lite",
		"https://db-ip.com/db/download/ip-to-country-lite",
		"https://db-ip.com/db/download/ip-to-city-lite",
	}

	// 初始化 Chrome 选项
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoDefaultBrowserCheck,  // 防止检查 Chrome 是否为默认浏览器
		chromedp.Flag("headless", true), // 排除无头模式
		chromedp.ExecPath(browserPath),  // 设置 Chrome 的执行路径
		chromedp.UserDataDir(""),        // 使用临时用户配置文件，即ignoring any existing user profiles
		// 设置下载选项
		chromedp.Flag("download.default_directory", outputDir),
		chromedp.Flag("download.prompt_for_download", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("download.directory_upgrade", true),
		chromedp.Flag("safebrowsing.enabled", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
	}

	// 创建 Chrome 实例
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	fmt.Println("chromedp allocator context created")

	// 循环处理页面访问及下载动作
	for _, url := range urls {
		// 先获取下载链接和 SHA1SUM
		fileInfo, sha1Info := parseDownloadInfo(url)

		// 记录到日志文件
		for fileType, downloadLink := range fileInfo {
			webSHA1SUM := sha1Info[fileType]

			// 在控制台显示信息
			fmt.Printf("URL: %s\n", url)
			fmt.Printf("File Type: %s\n", fileType)
			fmt.Printf("Download Link: %s\n", downloadLink)
			fmt.Printf("SHA1SUM: %s\n", webSHA1SUM)

			_, err = fmt.Fprintf(logFile, "DownloadLink: %s\n", downloadLink)
			if err != nil {
				log.Fatal("Failed to write to log file:", err)
			}
			_, err = fmt.Fprintf(logFile, "webSHA1SUM: %s\n", webSHA1SUM)
			if err != nil {
				log.Fatal("Failed to write to log file:", err)
			}
		}
		// 为每个 URL 创建一个新的上下文
		ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancelCtx()
		ctx, cancelCtx = context.WithTimeout(ctx, 60*time.Second)
		defer cancelCtx()

		// 为每个 URL 初始化下载计数器和完成通道
		done := make(chan bool)
		downloadCounter := 0
		totalDownloads := 2 // 每个 URL 预期下载的文件总数

		// 运行下载文件的函数，开始下载文件
		err := downloadFile(ctx, url, browserPath, outputDir, &downloadCounter, totalDownloads, done)
		if err != nil {
			log.Fatal(err)
		}

		// 开始异步处理下载文件的解压与SHA1SUM比较

		// 使用 os.ReadDir 获取 outputDir 目录中的所有文件
		dirEntries, err := os.ReadDir(outputDir)
		if err != nil {
			log.Fatalf("Failed to read output directory: %v", err)
		}

		// 转换为文件名列表
		var downloadedFileNames []string
		for _, entry := range dirEntries {
			if !entry.IsDir() {
				// 获取完整的文件路径
				downloadedFilePath := filepath.Join(outputDir, entry.Name())
				downloadedFileNames = append(downloadedFileNames, downloadedFilePath)
			}
		}
		fmt.Printf("Downloaded files in output dir: %v\n", downloadedFileNames)

		var wg sync.WaitGroup
		for _, file := range downloadedFileNames {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				// 从文件名提取文件类型
				var fileType string
				if strings.HasSuffix(f, ".mmdb.gz") {
					fileType = ".mmdb.gz"
				} else if strings.HasSuffix(f, ".csv.gz") {
					fileType = ".csv.gz"
				}
				expectedSHA1 := sha1Info[fileType]
				if err := processAndVerifyFile(f, expectedSHA1, outputDir); err != nil {
					// 错误处理
					fmt.Printf("Error processing file %s: %v\n", f, err)
				}
			}(file)
		}
		wg.Wait()

	}
}

// 下载文件的函数
func downloadFile(ctx context.Context, url string, browserPath string, outputDir string, downloadCounter *int, totalDownloads int, done chan bool) error {
	// 设置下载完成的通知通道
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				*downloadCounter++
				if *downloadCounter == totalDownloads {
					done <- true
					close(done)
				}
			} else if ev.TotalBytes > 0 {
				// 计算并打印下载进度
				progress := float64(ev.ReceivedBytes) / float64(ev.TotalBytes) * 100
				fmt.Printf("下载进度：%.2f%%\n", progress)
			}
		}
	})

	// 执行 chromedp 任务
	if err := chromedp.Run(ctx,
		// 设置下载行为
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow). // 保留其原始文件名
											WithDownloadPath(outputDir).
											WithEventsEnabled(true),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#license_agree`, chromedp.ByID),
		printAction("License agreement visible"),

		// 执行页面上的Javascript代码
		chromedp.Evaluate(`document.querySelector('#license_agree').checked = true;`, nil),
		printAction("Checked checkbox"),
		chromedp.Sleep(2*time.Second), // 等待2秒以确保JavaScript逻辑执行完成

		// 现在假设license_agree已经被选中，点击MMDB的下载链接
		chromedp.WaitVisible(`a.free_download_link[href$=".mmdb.gz"]`, chromedp.ByQuery),
		printAction("Download link visible"),
		chromedp.Click(`a.free_download_link[href$=".mmdb.gz"]`, chromedp.ByQuery),
		printAction("Clicked mmdb file download link"),

		// 等待一段时间，以确保第一个文件开始下载
		chromedp.Sleep(5*time.Second),

		// 点击 CSV 的下载链接
		chromedp.WaitVisible(`a.free_download_link[href$=".csv.gz"]`, chromedp.ByQuery),
		printAction("CSV Download link visible"),
		chromedp.Click(`a.free_download_link[href$=".csv.gz"]`, chromedp.ByQuery),
		printAction("Clicked CSV file download link"),
	); err != nil {
		fmt.Printf("Failed to complete chromedp run: %v\n", err)
	}

	// 等待下载完成
	<-done

	return nil
}

// 解析 HTML 并提取两种文件的下载链接和 SHA1SUM
func parseDownloadInfo(url string) (map[string]string, map[string]string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo := make(map[string]string)
	sha1Info := make(map[string]string)

	types := []string{".mmdb.gz", ".csv.gz"}
	for _, t := range types {
		downloadLink := doc.Find(fmt.Sprintf("a[href$='%s']", t)).AttrOr("href", "")
		fileInfo[t] = downloadLink

		// 提取 SHA1SUM
		webSHA1SUM := doc.Find(fmt.Sprintf("div.card dl:contains('%s')", t)).Find("dt:contains('SHA1SUM') + dd").Text()
		sha1Info[t] = webSHA1SUM
	}

	return fileInfo, sha1Info
}

// processAndVerifyFile 解压文件并验证 SHA1SUM
func processAndVerifyFile(gzipFilePath, expectedSHA1, outputDir string) error {
	// 构造解压后的文件路径
	outputPath := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(gzipFilePath), ".gz"))

	// 解压 .gz 文件
	if err := decompressGzipFile(gzipFilePath, outputPath); err != nil {
		return fmt.Errorf("failed to decompress file: %v", err)
	}

	// 对解压后的文件执行 SHA1SUM 校验
	if verified := verifySHA1(outputPath, expectedSHA1); !verified {
		return fmt.Errorf("SHA1SUM mismatch for file: %s", outputPath)
	}

	// 如果需要，此处可以添加日志记录或其他处理

	return nil
}

// 解压.gz文件的函数
func decompressGzipFile(gzipPath, outputPath string) error {
	gzFile, err := os.Open(gzipPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzipReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, gzipReader)
	return err
}

// 对文件执行SHA1SUM校验
func verifySHA1(filePath, expectedSHA1 string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatal(err)
		return false
	}

	computedSHA1 := fmt.Sprintf("%x", hasher.Sum(nil))
	return computedSHA1 == expectedSHA1
}

// printAction is a helper function to print a debug message for chromeedp.Run
func printAction(message string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		fmt.Println(message)
		return nil
	}
}
