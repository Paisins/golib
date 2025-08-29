package utils

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type WorkerFunc func() error

type StressTestResult struct {
	TotalRequests  int64
	SuccessCount   int64
	ErrorCount     int64
	ActualDuration time.Duration
}

func RunStressTest(ctx context.Context, qps int, duration time.Duration, worker WorkerFunc, progress *atomic.Int64) *StressTestResult {
	if qps <= 0 {
		panic("qps must be greater than 0")
	}

	var (
		totalRequests atomic.Int64
		successCount  atomic.Int64
		errorCount    atomic.Int64
		wg            sync.WaitGroup
	)

	// 控制 QPS 的 ticker
	interval := time.Second / time.Duration(qps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 设置测试超时
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	startTime := time.Now()

	// 启动工作 goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if progress != nil {
						progress.Add(1)
					}
					totalRequests.Add(1)
					if err := worker(); err != nil {
						errorCount.Add(1)
					} else {
						successCount.Add(1)
					}
				}()

			case <-ctx.Done():
				return
			}
		}
	}()

	// 等待测试完成
	<-ctx.Done()
	wg.Wait()
	actualDuration := time.Since(startTime)

	return &StressTestResult{
		TotalRequests:  totalRequests.Load(),
		SuccessCount:   successCount.Load(),
		ErrorCount:     errorCount.Load(),
		ActualDuration: actualDuration,
	}
}

func StressTesting(targetQPS int, testDuration time.Duration, worker WorkerFunc) {
	// 设置中断信号
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()

	// 捕获Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("\nTest interrupted by user")
		cancel()
	}()

	// 打印系统信息
	fmt.Println("======= System Information =======")
	fmt.Printf("Go Version:      %s\n", runtime.Version())
	fmt.Printf("OS/Arch:         %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU Cores:       %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS:      %d\n", runtime.GOMAXPROCS(0))
	fmt.Println("=================================")
	fmt.Printf("\nStarting stress test: QPS=%d, Duration=%s\n\n", targetQPS, testDuration)

	// 用于实时进度的原子计数器
	var progress atomic.Int64
	startTime := time.Now()

	// 实时进度显示
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				current := progress.Load()
				elapsed := time.Since(startTime).Seconds()
				if elapsed > 0 {
					currentQPS := float64(current) / elapsed
					fmt.Printf("\rCurrent QPS: %6.0f, Requests: %d", currentQPS, current)
					os.Stdout.Sync()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// 执行压测
	result := RunStressTest(ctx, targetQPS, testDuration, worker, &progress)
	duration := time.Since(startTime)

	// 打印结果
	fmt.Println("\n\n======= Stress Test Result =======")
	fmt.Printf("Target QPS:         %d\n", targetQPS)
	fmt.Printf("Actual QPS:         %.2f\n", float64(result.TotalRequests)/duration.Seconds())
	fmt.Printf("Total Requests:     %d\n", result.TotalRequests)
	fmt.Printf("Successful:         %d\n", result.SuccessCount)
	fmt.Printf("Errors:             %d\n", result.ErrorCount)
	fmt.Printf("Error Rate:         %.2f%%\n", float64(result.ErrorCount)*100/float64(result.TotalRequests))
	fmt.Printf("Test Duration:      %s\n", duration.Round(time.Millisecond))
	fmt.Printf("Theoretical Max:    %d (based on worker latency)\n", int(math.Ceil(1000.0/5.0)))
	fmt.Println("=================================")
}
