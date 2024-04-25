package tests

import (
	"crypto/md5"
	"fmt"
	"io"
	bpfScript "linux-monitoring-utility/internal/bpfScript"
	"os"
	"testing"
)

func TestBpfScriptLayer(t *testing.T) {
	testCases := []struct {
		inputSyscalls   []string
		inputPath       string
		expectedFile    string
		expectedMessage string
	}{
		{[]string{}, "", "./defaultScriptCheck.bt", ""},
		{[]string{"readlink", "readlinkat"}, "", "./generateScriptCheck.bt", ""},
		{[]string{"name_to_handle_at"}, "", "", "The system call 'name_to_handle_at' is not valid."},
		{[]string{}, "./res/", "./defaultScriptCheck.bt", ""},
		{[]string{}, "./res/\t//", "", "The path './res/	//' could not be created."},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			bpfScriptFile, err := bpfScript.GenerateBpfScript(tc.inputSyscalls, tc.inputPath)
			if err != nil && err.Error() != tc.expectedMessage {
				t.Error("Incorrectly generated file\n")
			} else if !Equal(bpfScriptFile, tc.expectedFile) {
				t.Error("The generated file does not match the expected one\n")
			}
		})
	}
	os.RemoveAll("./res/")
}

func Equal(fileOld *os.File, fileCheck string) bool {
	if fileCheck == "" {
		return true
	}
	files := []*os.File{}

	fileOld, err := os.OpenFile(fileOld.Name(), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return false
	}
	defer fileOld.Close()
	files = append(files, fileOld)

	fileNew, err := os.OpenFile(fileCheck, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return false
	}
	defer fileNew.Close()
	files = append(files, fileNew)

	checksums := []string{}
	for _, f := range files {
		_, err := f.Seek(0, 0)
		if err != nil {
			return false
		}
		sum, err := getMD5SumString(f)
		if err != nil {
			return false
		}
		checksums = append(checksums, sum)
	}
	//fmt.Println("### Сравнение по контрольной сумме ###")
	return compareCheckSum(checksums[0], checksums[1])
}

func getMD5SumString(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func compareCheckSum(sum1, sum2 string) bool {
	return sum1 == sum2
	//fmt.Printf("MD5: %s и MD5: %s %s совпадают\n", sum1, sum2)
}
