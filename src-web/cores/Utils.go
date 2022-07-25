package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logserver"
	"Codeforces-ContestCodeDownload/src-web/model"
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func getRandomStringHex(strLen int) string {
	if strLen <= 0 {
		return string([]byte{})
	}
	var need int
	if strLen&1 == 0 {
		need = strLen
	} else {
		need = strLen + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return string([]byte{})
	}
	hex.Encode(dst, src)
	return string(dst[:strLen])
}

func parseJsonFiles(infoForID gjson.Result) model.InformationStruct {
	var temp model.InformationStruct
	temp.ID = infoForID.Get(`0.id`).Int()
	temp.CID = infoForID.Get(`0.contestId`).Int()
	temp.PID = infoForID.Get(`0.problem.index`).String()
	temp.PNAME = infoForID.Get(`0.problem.name`).String()
	if infoForID.Get(`0.author.members`).Int() == 1 {
		temp.CNAME = infoForID.Get(`0.author.members.0.handle`).String()
	} else {
		temp.CNAME = infoForID.Get(`0.author.members.0.name`).String()
	}
	temp.LANG = infoForID.Get(`0.programmingLanguage`).String()
	return temp
}

func getAllAcceptSubmissionID(apiJsonString string) []string {
	allContestantResult := gjson.Get(apiJsonString, `result.#(author.participantType="CONTESTANT")#.id`)
	allVerdictOKResult := gjson.Get(apiJsonString, `result.#(verdict="OK")#.id`)
	return intersectGjsonResult(allContestantResult, allVerdictOKResult)
}

func getAPIJsonString(signedURL, randomUID string) string {
	apiData, err := http.Get(signedURL)
	if err != nil {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("An error occurred while getting API URL.")
		return ""
	}
	apiBytes, _ := ioutil.ReadAll(apiData.Body)
	apiJsonString := string(apiBytes)
	return apiJsonString
}

// ZipCompress
// Copy from https://studygolang.com/articles/34943
func ZipCompress(srcDir, zipFileName, randomUID string) {
	logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
		"srcDir":      srcDir,
		"zipFileName": zipFileName,
		"randomUID":   randomUID,
	}).Infoln("Start ZipCompress.")
	zipFileName = zipFileName + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		logserver.GetLogMap(randomUID).Errorln("err: " + err.Error())
		return
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			logserver.GetLogMap(randomUID).Errorln(err.Error())
		}
	}(zipFile)
	archive := zip.NewWriter(zipFile)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			logserver.GetLogMap(randomUID).Errorln(err.Error())
		}
	}(archive)
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
		if path == srcDir {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		if info.IsDir() {
			var specialChar string
			if runtime.GOOS == "windows" {
				specialChar = "\\"
			} else {
				specialChar = "/"
			}
			header.Name += specialChar
		} else {
			header.Name = removeRedundantPartOfTheFileName(path)
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					logserver.GetLogMap(randomUID).Errorln(err.Error())
				}
			}(file)
			_, err := io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logserver.GetLogMap(randomUID).Errorln(err.Error())
		return
	}
}

// removeRedundantPartOfTheFileName \folder1\folder2\folder3\xxx.cpp -> xxx.cpp
func removeRedundantPartOfTheFileName(fileName string) string {
	length := len(fileName) - 1
	// TODO E: 最後在Linux平臺上會出現錯誤，疑似路徑表示方式不同
	var specialChar uint8
	if runtime.GOOS == "windows" {
		specialChar = '\\'
	} else {
		specialChar = '/'
	}
	for fileName[length] != specialChar {
		length--
	}
	return fileName[length+1:]
}
