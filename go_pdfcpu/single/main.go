package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// convertEPSToPDF использует Ghostscript для конвертации eps в pdf
func convertEPSToPDF(epsPath, pdfPath string) error {
	// Команда "gs" для Linux/macOS. Для Windows используйте "gswin64c"
	cmd := exec.Command("gswin64c", "-dNOPAUSE", "-dBATCH", "-dPDFSETTINGS=/default", "-dEPSCrop", "-sDEVICE=pdfwrite",
		"-sOutputFile="+pdfPath, epsPath)
	return cmd.Run()
}

func main() {
	// 1. Найти все EPS файлы
	files, err := filepath.Glob("*.eps")
	if err != nil || len(files) == 0 {
		log.Fatal("EPS файлы не найдены")
	}

	var tempPDFs []string
	outputPDF := "combined_output.pdf"

	// 2. Конвертировать каждый EPS в временный PDF
	for _, epsFile := range files {
		tempPDF := strings.TrimSuffix(epsFile, filepath.Ext(epsFile)) + ".tmp.pdf"
		fmt.Printf("Конвертация %s -> %s...\n", epsFile, tempPDF)

		err := convertEPSToPDF(epsFile, tempPDF)
		if err != nil {
			log.Printf("Ошибка конвертации %s: %v", epsFile, err)
			continue
		}
		tempPDFs = append(tempPDFs, tempPDF)
	}

	// 3. Объединить PDF файлы с помощью pdfcpu
	fmt.Println("Объединение файлов в", outputPDF)
	err = api.MergeCreateFile(tempPDFs, outputPDF, false, nil)
	if err != nil {
		log.Fatal("Ошибка объединения:", err)
	}

	// 4. Удалить временные файлы
	for _, f := range tempPDFs {
		os.Remove(f)
	}

	fmt.Println("Готово!")
}
