package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
}

// Phương thức ghi log cấp độ INFO
func (l *Logger) Info(msg string) {
	str := fmt.Sprintf("[INFO]    %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	l.writeLog(str)
}

// Phương thức ghi log cấp độ ERROR
func (l *Logger) Error(msg string) {
	str := fmt.Sprintf("[ERROR]   %s: %s \n", time.Now().Format("2006-01-02 15:04:05"), msg)
	l.writeLog(str)
}
func (l *Logger) Fatal(err error) {
	str := fmt.Sprintf("[FATAL]   %s: %s \n", time.Now().Format("2006-01-02 15:04:05"), err)
	l.writeLog(str)
}

// Phương thức ghi log cấp độ WARNING
func (l *Logger) Warning(msg string) {
	str := fmt.Sprintf("[WARNING] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	l.writeLog(str)
}

// Phương thức ghi log vào file
func (l *Logger) writeLog(str string) {
	currentFileName := fmt.Sprintf("./logs/%s.log", time.Now().UTC().Format("2006-01-02"))
	// Nếu thư mục logs không tồn tại, tạo mới
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		err := os.Mkdir("./logs", os.ModePerm)
		if err != nil {
			log.Fatalf("Could not create logs directory: %v", err)
		}
	}

	// Mở file để ghi log (tạo file nếu chưa có)
	f, err := os.OpenFile(currentFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open or create log file: %v", err)
	}
	defer f.Close()

	// Ghi nội dung log vào file
	if _, err := f.WriteString(str); err != nil {
		log.Fatalf("Could not write to log file: %v", err)
	}
}
