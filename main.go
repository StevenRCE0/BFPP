package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/alexmullins/zip"
)

type Configuration struct {
	Storage string `json:"storage"`
	PSK     string `json:"psk"`
	FragmentSize float64 `json:"fragment"`
}

type BackupEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func createZip(origin, zipPath, psk string, zipWriter *zip.Writer) {
	content, err := ioutil.ReadFile(origin)
	if err != nil {
		log.Fatal(err)
	}
	written, err := zipWriter.Encrypt(zipPath, psk)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(written, bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
}

func createZipRecurrsive(origin, zipPath, psk string, zipWriter *zip.Writer) {
	files, err := ioutil.ReadDir(origin)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			createZipRecurrsive(origin+"/"+file.Name(), zipPath+"/"+file.Name(), psk, zipWriter)
		} else {
			createZip(origin+"/"+file.Name(), zipPath+"/"+file.Name(), psk, zipWriter)
		}
	}
}

func readConfiguration(file string) Configuration {
	var config Configuration
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func syncToPan(configuration Configuration) {
	os.Chdir(configuration.Storage)
	fmt.Println("Uploading to Pan")
	command := exec.Command("bypy", "syncup")
	out, err := command.CombinedOutput()
	if err != nil {
		fmt.Printf(string(out))
		log.Fatalf(err.Error())
	}
	fmt.Printf(string(out))
}

func main() {
	configuration := readConfiguration(os.Args[1])
	backupEntryFile, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	var backupEntries []BackupEntry
	err = json.Unmarshal(backupEntryFile, &backupEntries)
	if err != nil {
		log.Fatal(err)
	}
	timeFormatString := strconv.Itoa(time.Now().Year()) + "-" + strconv.Itoa(int(time.Now().Month())) + "-" + strconv.Itoa(time.Now().Day())
	for i := 0; i < len(backupEntries); i++ {
		fmt.Println(backupEntries[i].Path)
		origin, targetPath, targetZipName := backupEntries[i].Path, configuration.Storage+"/"+backupEntries[i].Name, timeFormatString+".zip"
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		zipFile, err := os.Create(targetPath + "/" + targetZipName)
		if err != nil {
			log.Fatal(err)
		}
		zipWriter := zip.NewWriter(zipFile)
		createZipRecurrsive(origin, "", configuration.PSK, zipWriter)
		zipWriter.Close()
		zipWriter.Flush()
	}
	syncToPan(configuration)
}
