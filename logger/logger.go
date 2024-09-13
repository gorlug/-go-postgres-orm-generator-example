package logger

import (
	"fmt"
	"strings"
)

func LogDebug(message string, a ...any) {
	println(fmt.Sprintf(message, a...))
}

func Debug(a ...any) {
	elementsToLog := []string{}

	for i := 0; i < len(a); i++ {
		elementsToLog = append(elementsToLog, strings.TrimSpace(fmt.Sprint(a[i])))
	}
	println(strings.Join(elementsToLog, " "))
}

func LogError(message string, a ...any) {
	println(fmt.Sprintf("Error: %v", fmt.Sprintf(message, a...)))
}

func Error(a ...any) {
	elementsToLog := []string{"Error:"}

	for i := 0; i < len(a); i++ {
		elementsToLog = append(elementsToLog, strings.TrimSpace(fmt.Sprint(a[i])))
	}
	println(strings.Join(elementsToLog, " "))
}
