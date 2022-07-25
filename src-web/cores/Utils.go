package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logMode"
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

func parseJsonFiles(infoForID gjson.Result) InformationStruct {
	var temp InformationStruct
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
		logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
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
func ZipCompress(srcDir string, zipFileName string) {
	zipFileName = zipFileName + ".zip"
	zipFile, _ := os.Create(zipFileName)
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
		if path == srcDir {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Name = resolveFileName(path)
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}

// resolveFileName \folder1\folder2\folder3\xxx.cpp -> xxx.cpp
func resolveFileName(fileName string) string {
	length := len(fileName) - 1
	for fileName[length] != '\\' {
		length--
	}
	return fileName[length+1:]
}
