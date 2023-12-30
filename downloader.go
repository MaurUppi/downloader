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
	"regexp"
	"strings"
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

	// 获取当前工作目录，并在此基础上创建tmp目录保存下载文件
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Working dir is : %s\n", wd)

	// 解析页面获取下载链接URL、文件名和SHA1SUM
	downloadLink, filename, webSHA1SUM := parseDownloadInfo("https://db-ip.com/db/download/ip-to-city-lite")
	// 在 Console 显示信息
	fmt.Printf("Download Link: %s\n", downloadLink)
	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("SHA1SUM: %s\n", webSHA1SUM)


	// 初始化 Chrome 选项
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoDefaultBrowserCheck,   // 防止检查 Chrome 是否为默认浏览器
		chromedp.Flag("headless", true), // 排除无头模式
		chromedp.ExecPath(browserPath),   // 设置 Chrome 的执行路径
		chromedp.UserDataDir(""),         // 使用临时用户配置文件，即ignoring any existing user profiles
		// 设置下载选项
		chromedp.Flag("download.default_directory", wd),
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

	// 创建上下文实例
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	fmt.Println("chromedp context created")

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second) // 120秒内无法完成下载操作则为失败
	defer cancel()
	fmt.Println("Starting chromedp run")

	// 设置下载完成的通知通道
	done := make(chan string, 1)
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}
	})


	// 使用chromedp模拟浏览器行为，会受到context.WithTimeout()的影响
	if err := chromedp.Run(ctx,
		// 设置下载行为
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(wd).
			WithEventsEnabled(true),
		chromedp.Navigate(`https://db-ip.com/db/download/ip-to-city-lite`),
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
		printAction("Clicked download link"),
	); err != nil {
		fmt.Printf("Failed to complete chromedp run: %v\n", err)
		//log.Fatalf("Failed to complete chromedp run: %v", err)
	}

	// 等待下载完成
	guid := <-done
	//log.Printf("下载完成，文件路径：%s", filepath.Join(wd, guid))
	fmt.Printf("chromedp run completed, dir is: %s", filepath.Join(wd, guid))

	// 完整的文件路径
	fullFilePath := filepath.Join(wd, guid)
	//fmt.Printf("Full file path: %s\n", fullFilePath)
	
	// 解压.gz文件
	err = decompressGzipFile(fullFilePath, strings.TrimSuffix(fullFilePath, ".gz"))
	if err != nil {
		log.Fatalf("Failed to decompress file: %v", err)
	}

	// 校验SHA1SUM
	decompressedFilePath := strings.TrimSuffix(fullFilePath, ".gz")
	if !verifySHA1(decompressedFilePath, webSHA1SUM) {
		log.Fatalf("SHA1SUM verification failed")
	}

	// 变更解压后的文件名
	re := regexp.MustCompile(`-\d{4}-\d{2}`)
	newFilename := filepath.Join(wd, re.ReplaceAllString(filepath.Base(decompressedFilePath), ""))
	err = os.Rename(decompressedFilePath, newFilename)
	if err != nil {
		log.Fatalf("Failed to rename file: %v", err)
	}
	fmt.Println("File renamed successfully")

	// 打开或创建日志文件，使用 os.O_TRUNC 来覆盖旧内容
	logFile, err := os.OpenFile("DW-dbip-city-lite.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// 使用 fmt.Fprintf 直接写入成功消息
	fmt.Fprintf(logFile, "Successed Download Link: %s\n", downloadLink)
	fmt.Fprintf(logFile, "Verified SHA1SUM: %s\n", webSHA1SUM)

}

// Function to parse HTML and extract download link and SHA1SUM
func parseDownloadInfo(url string) (string, string, string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	downloadLink := doc.Find("a[href$='.mmdb.gz']").AttrOr("href", "")

	// 提取下载链接中的文件名
	parts := strings.Split(downloadLink, "/")
	filename := parts[len(parts)-1] // 获取最后一个部分作为文件名

	webSHA1SUM := doc.Find("div.card dl:contains('MMDB')").Find("dt:contains('SHA1SUM') + dd").Text()

	return downloadLink, filename, webSHA1SUM
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

// printAction is a helper function to print a debug message for chromeedp.Run
func printAction(message string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		fmt.Println(message)
		return nil
	}
}
