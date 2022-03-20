package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Benyam-S/onemembership/entity"
	"github.com/Benyam-S/onemembership/log"
)

func main() {

	var sysConfig entity.SystemConfig
	configFilesDir := "C:/Users/Administrator/go/src/github.com/Benyam-S/onemembership/config"
	maxSize := 2097152                // 2MB
	interval := int64(time.Hour) * 24 // 1 day

	// Reading data from config.server.json file and creating the systemconfig  object
	sysConfigDir := filepath.Join(configFilesDir, "/config.server.json")
	sysConfigData, _ := ioutil.ReadFile(sysConfigDir)

	err := json.Unmarshal(sysConfigData, &sysConfig)
	if err != nil {
		panic(err)
	}

	logs := &log.LogContainer{
		UserLogFile:             filepath.Join(sysConfig.LogsPath, sysConfig.Logs["user_log_file"]),
		ServiceProviderLogFile:  filepath.Join(sysConfig.LogsPath, sysConfig.Logs["service_provider_log_file"]),
		ProjectLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["project_log_file"]),
		SubscriptionPlanLogFile: filepath.Join(sysConfig.LogsPath, sysConfig.Logs["subscription_plan_log_file"]),
		SubscriptionLogFile:     filepath.Join(sysConfig.LogsPath, sysConfig.Logs["subscription_log_file"]),
		TransactionLogFile:      filepath.Join(sysConfig.LogsPath, sysConfig.Logs["transaction_log_file"]),
		DeletedLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["deleted_log_file"]),
		ServerLogFile:           filepath.Join(sysConfig.LogsPath, sysConfig.Logs["server_log_file"]),
		BotLogFile:              filepath.Join(sysConfig.LogsPath, sysConfig.Logs["bot_log_file"]),
		ErrorLogFile:            filepath.Join(sysConfig.LogsPath, sysConfig.Logs["error_log_file"]),
		ArchiveLogFile:          filepath.Join(sysConfig.LogsPath, sysConfig.Logs["archive_log_file"]),
	}
	logger := log.NewLogger(logs, log.Normal)

	// Checking the validity of the given log file
	logFiles := []string{logs.UserLogFile, logs.ServiceProviderLogFile, logs.SubscriptionLogFile,
		logs.DeletedLogFile, logs.ServerLogFile, logs.BotLogFile, logs.ErrorLogFile}

	Archive(logger, logFiles, int64(maxSize), interval, sysConfig.ArchivesPath)
}

// Archive is a function or a process that archives log files when they reaches a certain size
func Archive(logger log.ILogger, logFiles []string, maxSize, interval int64, archiveLocation string) {

	for {

		time.Sleep(time.Duration(interval))

		for _, logFile := range logFiles {

			file, err := os.OpenFile(logFile, os.O_RDWR, 0644)
			if err != nil {
				logger.LogToArchiveFile(err.Error())
				continue
			}

			status, err := file.Stat()
			if err != nil {
				logger.LogToArchiveFile(err.Error())
				file.Close()
				continue
			}

			if maxSize <= status.Size() {

				// Logging entry point
				logger.LogToArchiveFile(fmt.Sprintf("--------------- Start archiving file %s", logFile))

				current_time := time.Now()
				timeStamp := fmt.Sprintf("%d%02d%02d%02d%02d%02d",
					current_time.Year(), current_time.Month(), current_time.Day(),
					current_time.Hour(), current_time.Minute(), current_time.Second())

				archivedFileName := fmt.Sprintf("%s_%s.%s", filepath.Base(file.Name()), timeStamp, "tar")
				archivedFilePath := filepath.Join(archiveLocation, archivedFileName)
				archivedFile, err := os.Create(archivedFilePath)
				if err != nil {
					logger.LogToArchiveFile(err.Error())
					file.Close()
					continue
				}

				archiver := tar.NewWriter(archivedFile)
				hdr := &tar.Header{
					Name: filepath.Base(file.Name()),
					Mode: 0600,
					Size: status.Size(),
				}

				if err := archiver.WriteHeader(hdr); err != nil {
					logger.LogToArchiveFile(err.Error())
					archivedFile.Close()
					file.Close()
					continue
				}

				output, err := ioutil.ReadAll(file)
				if err != nil {
					logger.LogToArchiveFile(err.Error())
					archivedFile.Close()
					file.Close()
					continue
				}

				if _, err := archiver.Write(output); err != nil {
					logger.LogToArchiveFile(err.Error())
					archivedFile.Close()
					file.Close()
					continue
				}

				if err := archiver.Close(); err != nil {
					logger.LogToArchiveFile(err.Error())
					archivedFile.Close()
					file.Close()
					continue
				}

				err = archivedFile.Close()
				if err != nil {
					logger.LogToArchiveFile(err.Error())
					file.Close()
					continue
				}

				// cleaning the file
				err = file.Truncate(0)
				if err != nil {
					logger.LogToArchiveFile(err.Error())
					file.Close()
					continue
				}

				// Logging finishing point
				logger.LogToArchiveFile(fmt.Sprintf("--------------- Finished archiving file %s", logFile))
			}

			err = file.Close()
			if err != nil {
				logger.LogToArchiveFile(err.Error())
				continue
			}

		}
	}
}
