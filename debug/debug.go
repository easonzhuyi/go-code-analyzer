package debug

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

var (
    // Debug 模式开关
    DebugMode = os.Getenv("GO_ANALYZER_DEBUG") == "1"
    startTime = time.Now()
)

// Log 输出调试信息
func Log(format string, args ...interface{}) {
    if !DebugMode {
        return
    }
    
    // 获取调用者信息
    _, file, line, _ := runtime.Caller(1)
    elapsed := time.Since(startTime).Milliseconds()
    prefix := fmt.Sprintf("[DEBUG %dms] %s:%d", elapsed, filepath.Base(file), line)
    
    fmt.Printf(prefix+" "+format+"\n", args...)
}

// DumpEnv 打印环境信息
func DumpEnv() {
    if !DebugMode {
        return
    }
    
    fmt.Println("\n=== Environment Info ===")
    fmt.Printf("GOPATH: %s\n", os.Getenv("GOPATH"))
    fmt.Printf("PWD: %s\n", getCurrentDir())
    fmt.Printf("GO_ANALYZER_DEBUG: %s\n", os.Getenv("GO_ANALYZER_DEBUG"))
    fmt.Println("=====================\n")
}

func getCurrentDir() string {
    dir, err := os.Getwd()
    if err != nil {
        return "unknown"
    }
    return dir
}

// Timer 用于性能分析
type Timer struct {
    name      string
    startTime time.Time
}

func NewTimer(name string) *Timer {
    if !DebugMode {
        return nil
    }
    return &Timer{
        name:      name,
        startTime: time.Now(),
    }
}

func (t *Timer) Stop() {
    if t == nil {
        return
    }
    elapsed := time.Since(t.startTime)
    Log("%s took %v", t.name, elapsed)
}
