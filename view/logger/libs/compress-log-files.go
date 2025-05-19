package logger

import (
	"os"
)

func CompressFile(f *os.File, compressFolderPath string, level string){
	// Comprime el archivo a GZIP
		// Crea un archivo "level-compressed-timestamp.txt.gz"
		// Copia el contenido del archivo original al archivo comprimido
		// Elimina el archivo original
		// Mueve el archivo GZIP a una carpeta donde se encuentran los archivos comprimidos
}