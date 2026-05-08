package utils

import (
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAccents remove acentos e marcas de uma string.
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SanitizeFilename limpa strings para serem usadas com segurança como arquivos ou pastas.
func SanitizeFilename(name string, isFile bool) string {
	if name == "" {
		if isFile {
			return "sem_titulo"
		}
		return "pasta_sem_nome"
	}

	temp := RemoveAccents(name)

	if isFile {
		temp = strings.ReplaceAll(temp, " ", "_")
	}

	re := regexp.MustCompile(`[\\/*?:"<>|]`)
	temp = re.ReplaceAllString(temp, "")

	if isFile {
		ext := filepath.Ext(temp)
		base := strings.TrimSuffix(temp, ext)

		reAlpha := regexp.MustCompile(`[^\w_-]`)
		base = reAlpha.ReplaceAllString(base, "")

		reUnder := regexp.MustCompile(`_+`)
		base = reUnder.ReplaceAllString(base, "_")

		base = strings.Trim(base, "_-")
		if base == "" {
			base = "arquivo_sem_nome"
		}
		temp = base + ext
	} else {
		reFolder := regexp.MustCompile(`[^\w\s_-]`)
		temp = reFolder.ReplaceAllString(temp, "")

		reSpace := regexp.MustCompile(`[\s_-]+`)
		temp = reSpace.ReplaceAllString(temp, " ")

		temp = strings.TrimSpace(temp)
	}

	if temp == "" {
		if isFile {
			return "sem_titulo"
		}
		return "pasta_sem_nome"
	}
	return temp
}

// ToWSLPath converte um caminho Windows (ex: D:\...) para o formato WSL (/mnt/d/...)
// se o caminho parecer ser um caminho Windows e estivermos em um ambiente Linux.
func ToWSLPath(path string) string {
	// Verifica se parece um caminho Windows (ex: C:\ ou D:\)
	if len(path) >= 3 && path[1] == ':' && path[2] == '\\' {
		drive := strings.ToLower(string(path[0]))
		remaining := strings.ReplaceAll(path[3:], "\\", "/")
		return "/mnt/" + drive + "/" + remaining
	}
	// Também lida com caminhos que usam barras normais mas tem letra de unidade (ex: D:/...)
	if len(path) >= 3 && path[1] == ':' && path[2] == '/' {
		drive := strings.ToLower(string(path[0]))
		return "/mnt/" + drive + "/" + path[3:]
	}
	return path
}

// extractNumbers é usado pela NaturalSort para dividir string e números.
func extractNumbers(s string) []string {
	re := regexp.MustCompile(`\d+|\D+`)
	return re.FindAllString(s, -1)
}

// NaturalSort ordena uma lista de strings mantendo a progressão numérica correta (ex: 1 antes de 2).
func NaturalSort(files []string) {
	sort.Slice(files, func(i, j int) bool {
		a := extractNumbers(files[i])
		b := extractNumbers(files[j])

		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				numA, errA := strconv.Atoi(a[k])
				numB, errB := strconv.Atoi(b[k])

				if errA == nil && errB == nil {
					return numA < numB
				}
				return strings.ToLower(a[k]) < strings.ToLower(b[k])
			}
		}
		return len(a) < len(b)
	})
}
