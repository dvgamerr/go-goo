package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/barasher/go-exiftool"
)

type Options struct {
	Grouping bool   `arg:"-g,--group"`
	Write    bool   `arg:"-w,--write"`
	Dir      string `arg:"positional" default:"."`
}

var flag Options

func main() {
	arg.MustParse(&flag)
	flagDirName, err := filepath.Abs(flag.Dir)
	if err != nil {
		log.Fatal(err)
	}

	finfo, err := os.Stat(flagDirName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("intializing...")

	et, err := exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer et.Close()

	if finfo.IsDir() {
		err := filepath.WalkDir(flagDirName, walk(et))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		fmt.Printf("done")
		if flag.Write {
			fmt.Printf(", rewrited")
		}
		fmt.Println(".")
	} else {
		printOriginalDate(et.ExtractMetadata(flagDirName))
	}

}

// var allowExt []string = []string{".jpg", ".jpeg", ".mp4", ".heic", ".png", ".mov"}

func walk(et *exiftool.Exiftool) func(s string, d fs.DirEntry, err error) error {
	return func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileInfos := et.ExtractMetadata(s)

			for i := range fileInfos {
				fileInfo := &fileInfos[i]
				if fileInfo.Err != nil {
					fmt.Printf(" - %v: %v\n", fileInfo.File, fileInfo.Err)
					continue
				}
				_, dateOriginal := getOriginalDate(fileInfo.Fields)

				fileName := strings.ReplaceAll(d.Name(), path.Ext(d.Name()), "")
				prefixDate := strings.ReplaceAll(dateOriginal[:7], ":", "-")

				if len(filepath.Base(filepath.Dir(s))) >= 7 {
					if prefixDate != filepath.Base(filepath.Dir(s))[:7] {
						dateLayout := "2006:01:02 15:04:05-07:00"
						location, _ := time.LoadLocation("Asia/Bangkok")

						i, err := strconv.ParseInt(fileName, 10, 64)
						if err == nil {
							dateOriginal = time.UnixMilli(i).In(location).Format(dateLayout)
						}
					}

					prefixDate := strings.ReplaceAll(dateOriginal[:7], ":", "-")
					if prefixDate != filepath.Base(filepath.Dir(s))[:7] {
						fmt.Printf(" - [%s>%s] >> %s\\%s\n", prefixDate, filepath.Base(filepath.Dir(s))[:7], filepath.Base(filepath.Dir(s)), d.Name())
					}
				}
				if flag.Write {
					for k := range fileInfo.Fields {
						if strings.Contains(k, "Date") {
							fileInfo.Fields[k] = dateOriginal
						}
					}
					fileInfo.Fields["DateTimeOriginal"] = dateOriginal
				}
			}
			if flag.Write {
				et.WriteMetadata(fileInfos)
			}
			// if slices.Contains(allowExt, strings.ToLower(path.Ext(s))) {
			// 	println(s)
			// } else {
			// 	println(s)
			// }
		} else if flag.Write {
			fmt.Printf(" > scan %s\n", filepath.Base(s))
		}
		return nil
	}
}

// func FileGroupping() {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	et, err := exiftool.NewExiftool()
// 	if err != nil {
// 		fmt.Printf("Error when intializing: %v\n", err)
// 		return
// 	}
// 	defer et.Close()

// 	if options.exif != "" {
// 		pwd = options.exif
// 	}

// 	entries, err := os.ReadDir(pwd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, e := range entries {
// 		log.Printf("dir: %s", e.Name())

// 		err = moveFileToDir(et, pwd, e)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// if e.IsDir() {
// 		// 	subDir := path.Join(pwd, e.Name())
// 		// 	subEntries, err := os.ReadDir(subDir)
// 		// 	if err != nil {
// 		// 		log.Fatal(err)
// 		// 	}
// 		// 	if skipMonthDir(et, subDir, subEntries) {
// 		// 		continue
// 		// 	}

// 		// 	for _, s := range subEntries {
// 		// 		err = moveFileToDir(et, subDir, s)
// 		// 		if err != nil {
// 		// 			log.Fatal(err)
// 		// 		}
// 		// 	}
// 		// }
// 	}
// }

func printOriginalDate(fileInfos []exiftool.FileMetadata) {
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("- %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			if strings.Contains(k, "Date") {
				fmt.Printf("  - %v = %s\n", v, k)
			}
		}
		keyDate, dateOriginal := getOriginalDate(fileInfo.Fields)
		fmt.Printf("- [%s]%v >> %s\n", keyDate, dateOriginal, path.Base(fileInfo.File))
	}
}

func checkValDate(val string) bool {
	return val != "" && val != "0000:00:00 00:00:00"
}

func getOriginalDate(Fields map[string]interface{}) (string, string) {
	k := "DateTimeOriginal"
	val, ok := Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "TrackCreateDate"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "SubSecDateTimeOriginal"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "MetadataDate"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "ModifyDate"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "DateTimeDigitized"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "SubSecCreateDate"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	k = "FileModifyDate"
	val, ok = Fields[k].(string)
	if ok && checkValDate(val) {
		return k, parseOriginalDate(val)
	}

	for k, v := range Fields {
		fmt.Printf("  - %v = %s\n", v, k)
	}
	log.Panic("unknow date")
	return "", ""
}

func parseOriginalDate(val string) string {
	dateLayout := "2006:01:02 15:04:05-07:00"
	location, _ := time.LoadLocation("Asia/Bangkok")

	t, err := time.Parse("2006:01:02 15:04:05", val)
	if err == nil {
		return t.Add(-7 * time.Hour).In(location).Format(dateLayout)
	}

	t, err = time.Parse("2006:01:02 15:04:05Z", val)
	if err == nil {
		return t.In(location).Format(dateLayout)
	}

	t, err = time.Parse("2006:01:02 15:04:05-07:00", val)
	if err == nil {
		return t.In(location).Format(dateLayout)
	}

	t, err = time.Parse("2006:01:02", val)
	if err == nil {
		return t.In(location).Format(dateLayout)
	}
	log.Panicf("%s --->> %v\n", val, err)
	return ""
}

// func moveFileToDir(et *exiftool.Exiftool, pwd string, e fs.DirEntry) error {
// 	currentFile := path.Join(pwd, e.Name())
// 	dirName := getDirNameCreateDate(et, currentFile)
// 	newFile := path.Join(pwd, dirName, e.Name())
// 	// log.Printf(" - %s >> %s", e.Name(), dirName)
// 	err := os.MkdirAll(path.Join(pwd, dirName), 0700)
// 	if err != nil {
// 		fmt.Printf("%s - %s\n", dirName, e.Name())
// 		return err
// 	}
// 	err = os.Rename(currentFile, newFile)
// 	if err != nil {
// 		fmt.Printf("%s - %s\n", dirName, e.Name())
// 		return err
// 	}
// 	return nil
// }

// // func skipMonthDir(et *exiftool.Exiftool, pwd string, entries []fs.DirEntry) bool {
// // 	var dirNew map[string]bool = make(map[string]bool)
// // 	for _, e := range entries {
// // 		currentFile := path.Join(pwd, e.Name())
// // 		dirName := getDirNameCreateDate(et, currentFile)
// // 		dirNew[dirName] = true
// // 	}
// // 	return len(dirNew) <= 1
// // }

// func getDirNameCreateDate(et *exiftool.Exiftool, filepath string) string {
// 	fileInfos := et.ExtractMetadata(filepath)

// 	dirName := ""
// 	for _, fileInfo := range fileInfos {
// 		if fileInfo.Err != nil {
// 			fmt.Printf("%s: %v\n", path.Base(fileInfo.File), fileInfo.Err)
// 			continue
// 		}

// 		t, err := getDate(fileInfo)
// 		if err != nil {
// 			fmt.Printf("%s: %v\n", path.Base(fileInfo.File), err)
// 			for k, v := range fileInfo.Fields {
// 				if strings.Contains(k, "date") {
// 					fmt.Printf("[%v] %v\n", k, v)
// 				}
// 			}
// 			continue
// 		}
// 		dirName = t.Format("2006-01")
// 	}
// 	return dirName
// }

// func getDate(fileInfo exiftool.FileMetadata) (time.Time, error) {
// 	date, ok := fileInfo.Fields["DateTimeOriginal"].(string)
// 	if !ok {
// 		date, ok = fileInfo.Fields["CreateDate"].(string)
// 		if !ok {
// 			date, ok = fileInfo.Fields["DateTimeDigitized"].(string)
// 			if !ok {
// 				date, _ = fileInfo.Fields["SubSecCreateDate"].(string)
// 			}
// 		}
// 	}
// 	return time.Parse("2006:01:02 15:04:05", date)
// }
