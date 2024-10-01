package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caiguanhao/readqr"
)

func main() {
	dir := "valid/"

	// Read all files from the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// Loop through each file in the directory
	for _, file := range files {
		if !file.IsDir() {
			// Construct full file path
			filePath := filepath.Join(dir, file.Name())

			// Open the image file
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", filePath, err)
				continue
			}
			defer f.Close()

			// Decode QR code
			result, err := readqr.Decode(f)
			if err != nil {
				fmt.Printf("Error decoding file %s: %v\n", filePath, err)
				continue
			}

			var newResult string
			var bankCode string
			checkfront := "10103"
			substringLength := 3
			checkBack := "5102"

			// Process the checkfront
			if strings.Contains(result, checkfront) {
				index := strings.Index(result, checkfront)
				if index != -1 {
					startIdx := index + len(checkfront)
					if startIdx+substringLength <= len(result) {
						bankCode = result[startIdx : startIdx+substringLength]
					} else {
						fmt.Printf("Not enough characters following '%s' to extract %d characters\n", checkfront, substringLength)
					}
					newResult = result[startIdx+substringLength:]
				}
			}

			// Process the checkBack
			if strings.Contains(newResult, checkBack) {
				index := strings.Index(newResult, checkBack)
				if index != -1 {
					newResult = newResult[:index]
				}
			} else {
				fmt.Printf("Check string '%s' not found in result\n", checkBack)
			}

			// Ensure newResult is long enough before cutting 4 characters
			if len(newResult) >= 4 {
				newResult = newResult[4:] // Cut the first 4 characters
			} else {
				fmt.Println("Not enough characters to cut 4 from the front of SlipContext")
			}

			// Remove the extension from the file name
			fileNameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			// Print the formatted output with the file name without extension
			fmt.Printf("Name:      \"TestValid_%s\",\n", fileNameWithoutExt) // Print file name without extension
			fmt.Printf("ImagePath: \"../../../../public/assets/%s\",\n", filePath)
			fmt.Printf("ExptRes: slipCtxWithBkRespStx{\n")
			fmt.Printf("	Slip: _slipDto.OCRAndQRTextResp{\n")
			fmt.Printf("		QRdata: \"%s\",\n", result)
			fmt.Printf("		SlipContext: \"%s\",\n", newResult) // SlipContext without the first 4 chars
			fmt.Printf("		BankCodeSender: \"%s\",\n", bankCode)
		}
	}
}
