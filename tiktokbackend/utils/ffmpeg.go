package utils

import (
	"fmt"
	"os/exec"
)

func SaveFaceImage(fileName string) error {
	inputFile := "./public/videos/" + fileName + ".mp4"
	// 设置转码后文件路径
	outputFile := "./public/images/" + fileName + ".mp4"

	// 设置 ffmpeg 命令行参数
	// ffmpeg -i input_file -y -f image2 -t 0.001 -s 352x240 output.jpg
	args := []string{"-i", inputFile, "-y", "-f", "image2", "-t", "0.001", outputFile}

	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// 运行 ffmpeg 命令
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
