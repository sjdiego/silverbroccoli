package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"silverbroccoli/pkg/reddit"
)

// BaseDir for storing temporary files
const (
	BaseDir = "./tmp/"
)

func newVideo(post reddit.Post) *reddit.Video {
	videoSufix := fmt.Sprintf("DASH_%d.mp4", post.Media.Video.Height)
	audioSufix := "DASH_audio.mp4"

	return &reddit.Video{
		VideoURL:  fmt.Sprintf("%s/%s", post.URL, videoSufix),
		AudioURL:  fmt.Sprintf("%s/%s", post.URL, audioSufix),
		VideoPath: fmt.Sprintf("%s%s_%s", BaseDir, post.ID, videoSufix),
		AudioPath: fmt.Sprintf("%s%s_%s", BaseDir, post.ID, audioSufix),
		FilePath:  fmt.Sprintf("%svid_%s.mp4", BaseDir, post.ID),
	}
}

// CreateVideo func
func CreateVideo(post reddit.Post) reddit.Video {
	v := newVideo(post)

	if !fileExists(v.VideoPath) {
		err := downloadFile(v.VideoURL, v.VideoPath)
		if err != nil {
			panic(err)
		}
	}

	if !fileExists(v.AudioPath) {
		err := downloadFile(v.AudioURL, v.AudioPath)
		if err != nil {
			panic(err)
		}
	}

	if !fileExists(v.FilePath) {
		log.Println("[*] Trying to mount video for " + post.Title)
		cmd := exec.Command(
			"ffmpeg",
			"-i", v.VideoPath,
			"-i", v.AudioPath,
			"-c:v", "copy",
			"-c:a", "aac",
			v.FilePath,
		)
		err := cmd.Run()

		if err != nil {
			log.Println("[!] Unable to mount video:", err)
		} else {
			log.Println("[*] Video created successfully")
		}
	}

	return *v
}

// Delete func
func Delete(post reddit.Post) {
	v := newVideo(post)

	files := []string{v.VideoPath, v.AudioPath, v.FilePath}

	for _, file := range files {
		if fileExists(file) {
			err := os.Remove(file)

			if err != nil {
				log.Println("[!] Error deleting file", file)
			} else {
				log.Println("[-] Removed file from filesystem", file)
			}
		}
	}
}

func downloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[!] Error downloading file for '%s', status: %d", url, resp.StatusCode)
		return err
	}
	log.Printf("[+] Downloaded file for '%s', status: %d", url, resp.StatusCode)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
