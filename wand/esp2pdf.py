from wand.image import Image
import os
import glob

def convert_eps_folder_to_pdf(folder_path, output_pdf='merged_output.pdf', resolution=300):
    """
    Конвертирует все EPS-файлы в указанной папке в один многостраничный PDF.

    Args:
        folder_path (str): путь к папке с EPS-файлами.
        output_pdf (str): имя выходного PDF-файла.
        resolution (int): разрешение для конвертации (DPI).
    """
    # Находим все EPS-файлы в папке (без учёта регистра)
    eps_patterns = [
        os.path.join(folder_path, '*.eps'),
        os.path.join(folder_path, '*.EPS')
    ]
    eps_files = []
    for pattern in eps_patterns:
        eps_files.extend(glob.glob(pattern))

    if not eps_files:
        print("В указанной папке не найдено EPS-файлов.")
        return

    print(f"Найдено EPS-файлов: {len(eps_files)}")
    print("Начинаем конвертацию...")

    with Image() as pdf:
        for eps_file in sorted(eps_files):  # сортируем для предсказуемого порядка страниц
            print(f"Обрабатываем: {os.path.basename(eps_file)}")
            try:
                # Загружаем EPS с заданным разрешением
                with Image(filename=eps_file, resolution=resolution) as img:
                    # Конвертируем в RGB для корректного отображения цветов
                    img.format = 'pdf'
                    # Добавляем страницу в общий PDF
                    pdf.sequence.append(img)
            except Exception as e:
                print(f"Ошибка при обработке {eps_file}: {e}")

        # Сохраняем многостраничный PDF
        pdf.save(filename=output_pdf)

    print(f"Конвертация завершена! Результат сохранён как: {output_pdf}")

# Пример использования
if __name__ == "__main__":
    # Укажите путь к вашей папке с EPS-файлами
    folder_with_eps = "C://Users//useroks//Downloads//a7b3d83a-d2ff-47b4-8801-818a54615c04_begin_offset_5000_number_of_codes_2000"  # замените на реальный путь
    output_filename = "result.pdf"  # имя выходного файла

    convert_eps_folder_to_pdf(
        folder_path=folder_with_eps,
        output_pdf=output_filename,
        resolution=300  # можно изменить на 150 или 600 для другого качества
    )
