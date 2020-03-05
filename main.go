package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func main() {
	var files []string

	var (
		inDir         string
		outDir        string
		startPlaylist int
		endPlaylist   int
		verbose       bool
	)
	flag.StringVar(&inDir, "in", "./", "input folder to scan recursively for mp3s")
	flag.StringVar(&outDir, "out", "./hoerbert", "destination folder to put hoerbert content to")
	flag.IntVar(&startPlaylist, "start", 0, "hoerbert button to start with (0-8)")
	flag.IntVar(&endPlaylist, "end", 8, "hoerbert button to end with (0-8)")
	flag.BoolVar(&verbose, "v", false, "verbose output")

	flag.Parse()

	if startPlaylist > 8 || startPlaylist < 0 {
		panic("Invalid start value: " + strconv.Itoa(startPlaylist))
	}

	if endPlaylist > 8 || endPlaylist < 0 {
		panic("Invalid start value: " + strconv.Itoa(startPlaylist))
	}

	if endPlaylist < startPlaylist {
		panic("end playlist value bigger than start playlist value.")
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.MkdirAll(outDir, 0777)
	}

	fmt.Println("Reading files from " + inDir)

	err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".mp3" || filepath.Ext(path) == ".MP3" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	playlistTrackCount := len(files) / (endPlaylist + 1 - startPlaylist)
	remainder := len(files) % (endPlaylist + 1 - startPlaylist)

	currentPlaylist := startPlaylist
	currentTrack := 0
	fmt.Printf("Writing %d files to %s, %d tracks per button. Remaining %d files distributed amongst the first buttons.\n", len(files), outDir, playlistTrackCount, remainder)
	for _, srcFile := range files {
		trackCountCurrentPlaylist := playlistTrackCount
		if currentPlaylist < startPlaylist+remainder {
			trackCountCurrentPlaylist++
		}
		dstFile := filepath.Join(outDir, strconv.Itoa(currentPlaylist), strconv.Itoa(currentTrack)+".WAV")

		if _, err := os.Stat(filepath.Join(outDir, strconv.Itoa(currentPlaylist))); os.IsNotExist(err) {
			os.MkdirAll(filepath.Join(outDir, strconv.Itoa(currentPlaylist)), 0777)
		}

		fmt.Printf("converting %s to %s\n", srcFile, dstFile)
		cmd := exec.Command("sox", "--buffer", "131072", "--multi-threaded", "--no-glob", srcFile, "--clobber", "-r", "32000", "-b", "16", "-e", "signed-integer", "--no-glob", dstFile,
			"remix", "-",
			"gain", "-n", "-1.5",
			"bass", "+1", "loudness", "-1",
			"pad", "0", "0", "dither",
		)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		if verbose {
			slurpstderr, _ := ioutil.ReadAll(stderr)
			fmt.Printf("STDERR:\n %s\n", slurpstderr)
			slurpstdout, _ := ioutil.ReadAll(stdout)
			fmt.Printf("STDOUT:\n %s\n", slurpstdout)
		}
		if err := cmd.Wait(); err != nil {
			log.Printf("Convertion finished with error: %v", err)
			panicoutput := fmt.Errorf("Command faild: %v", cmd)
			panic(panicoutput)
		}
		currentTrack++
		if currentTrack == trackCountCurrentPlaylist {
			currentTrack = 0
			currentPlaylist++
		}
	}
}
