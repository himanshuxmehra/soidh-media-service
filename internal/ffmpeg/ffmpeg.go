package ffmpeg

import (
	"fmt"
	"os"
	"go.uber.org/zap"
	"os/exec"
	"strconv"
	"strings"
)

func VideoConversion(filePath string){
	logger := zap.L()
	if (filePath == ""){
		logger.Error("File path is empty")
		return
	}
	inputFile := filePath
	height, err := getVideoHeight(inputFile)
	if err != nil {
		logger.Error("Error getting video height")
		// logger.Error("Error getting video height: %v", err)
		return
	}

	if height > 1080 {
		convertVideo(inputFile, []int{1080, 720, 360})
		// convertVideo(inputFile, 720)
		// convertVideo(inputFile, 360)
	} else if height > 720 {
		convertVideo(inputFile, []int{720, 360})
		// convertVideo(inputFile, 720)
		// convertVideo(inputFile, 360)
	} else if height > 360 {
		convertVideo(inputFile, []int{360})
		// convertVideo(inputFile, 360)
	} else {
		fmt.Println("Input video is already 360p or lower. No conversion needed.")
	}

	return
}


func getVideoHeight(inputFile string) (int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", 
		"-count_packets", "-show_entries", "stream=height", "-of", "csv=p=0", inputFile)
	
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	height, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return height, nil
}

func convertVideo(inputFile string, resolutions []int) error {
	outputDir := "output_hls"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	var filterComplex strings.Builder
	var streamMap strings.Builder

	for i, res := range resolutions {
		filterComplex.WriteString(fmt.Sprintf("[0:v]scale=-2:%d:flags=lanczos[v%d];", res, i))
		streamMap.WriteString(fmt.Sprintf(" -map [v%d] -map 0:a", i))
		streamMap.WriteString(fmt.Sprintf(" -c:v:%d libx264 -crf 23 -preset slow", i))
		streamMap.WriteString(fmt.Sprintf(" -c:a:%d aac -b:a:%d 128k", i, i))
		streamMap.WriteString(fmt.Sprintf(" -var_stream_map \"v:%d,a:%d\"", i, i))
	}

	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-filter_complex", filterComplex.String(),
		"-c:a", "aac",
		"-b:a", "128k",
		streamMap.String(),
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", fmt.Sprintf("%s/%%v_%%03d.ts", outputDir),
		"-master_pl_name", "master.m3u8",
		fmt.Sprintf("%s/%%v_stream.m3u8", outputDir))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Converting and segmenting video...")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting and segmenting video: %v", err)
	}

	fmt.Println("Conversion and segmentation complete.")
	return nil
}