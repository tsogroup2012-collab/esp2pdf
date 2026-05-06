package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func convertEPSToPDF(epsPath, pdfPath string) error {
	// Используем контекст или просто exec.Command
	cmd := exec.Command("gswin64c", "-dNOPAUSE", "-dBATCH", "-dPDFSETTINGS=/default", "-dEPSCrop", "-sDEVICE=pdfwrite",
		"-sOutputFile="+pdfPath, epsPath)
	return cmd.Run()
}

func main() {
	files, err := filepath.Glob("*.eps")
	if err != nil || len(files) == 0 {
		log.Fatal("EPS файлы не найдены")
	}

	outputPDF := "combined_output.pdf"

	// Используем WaitGroup для отслеживания горутин
	var wg sync.WaitGroup
	// Ограничиваем количество одновременных конвертаций (по числу ядер CPU)
	semaphore := make(chan struct{}, runtime.NumCPU())

	// Слайс для имен временных файлов (нужен mutex для безопасной записи из разных потоков)
	tempPDFs := make([]string, len(files))
	var mu sync.Mutex

	for i, epsFile := range files {
		wg.Add(1)
		go func(i int, epsFile string) {
			defer wg.Done()

			// Занимаем слот в семафоре
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			tempPDF := strings.TrimSuffix(epsFile, filepath.Ext(epsFile)) + ".tmp.pdf"
			fmt.Printf("[%d/%d] Конвертация %s...\n", i+1, len(files), epsFile)

			if err := convertEPSToPDF(epsFile, tempPDF); err != nil {
				log.Printf("Ошибка %s: %v", epsFile, err)
				return
			}

			// Безопасно добавляем файл в список для объединения
			mu.Lock()
			tempPDFs[i] = tempPDF
			mu.Unlock()
		}(i, epsFile)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// Фильтруем пустые строки (если были ошибки конвертации)
	var finalFiles []string
	for _, f := range tempPDFs {
		if f != "" {
			finalFiles = append(finalFiles, f)
		}
	}

	if len(finalFiles) == 0 {
		log.Fatal("Нет файлов для объединения")
	}

	fmt.Println("Объединение в", outputPDF)
	if err := api.MergeCreateFile(finalFiles, outputPDF, false, nil); err != nil {
		log.Fatal("Ошибка объединения:", err)
	}

	// Удаление временных файлов
	for _, f := range finalFiles {
		os.Remove(f)
	}

	fmt.Println("Готово!")
}
