package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile 读取文件内容
func ReadFile(path string) string {
	inputFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return ""
	}
	defer inputFile.Close()

	// Read all content at once
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	return string(content)
}

// WriteToFile 创建文件并将内容写入到文件
func WriteToFile(path, data string) {
	// 创建输出文件
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	// 创建一个 bufio.Writer 来写入输出文件
	writer := bufio.NewWriter(outputFile)

	// 将处理后的行写入输出文件
	_, err = writer.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	// 确保所有缓冲数据都写入输出文件
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
	}
	outputFile.Close()
}

// CreateFolder 创建目录
func CreateFolder(path string) {
	// Check if the directory already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// The directory does not exist, so create it using os.MkdirAll()
		err := os.MkdirAll(path, 0755) // 0755 sets the permissions for the directory
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}
}

// RemoveFolder 删除目录
func RemoveFolder(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
		return
	}
	fmt.Println("Directory deleted successfully!")
}

// RemoveFile 删除文件
func RemoveFile(file string) {
	err := os.Remove(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist.")
		} else {
			fmt.Println("Error deleting file:", err)
		}
	} else {
		fmt.Println("File deleted successfully!")
	}
}

// RemoveAllFilesInDir 删除目录下所有文件
func RemoveAllFilesInDir(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Printf("Deleted file: %s\n", path)
		}
		return nil
	})
}
