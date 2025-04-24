package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// ExecCommand ejecuta un comando del sistema operativo y devuelve su salida.
func ExecCommand(command string) {

	parts := strings.Fields(command)

	if len(parts) == 0 {
		fmt.Println("Comando vacío.")
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error ejecutando el comando:", err)
		return
	}

	fmt.Println(string(output))
}
