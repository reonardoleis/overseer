package functions

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Run(filename string, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	command := "node functions/" + filename + ".js"

	for _, arg := range args {
		command += " \"" + arg + "\""
	}

	output, err := exec.CommandContext(ctx, "bash", "-c", command).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Code(name, code string) string {
	code = "async function " + name + "(args) {\n" + code + "\n}\n\n"
	code += fmt.Sprintf(`process.argv.shift(); process.argv.shift();
	
	(async () => {
		console.log(await %s(process.argv));
	})();
	`, name)

	return code
}

var (
	forbiddenTokens = []string{
		"child_process",
		"fs",
		"require",
		"exec",
		"spawn",
		"execFile",
		"execFileSync",
		"execSync",
		"spawnSync",
		"readFileSync",
		"writeFileSync",
		"unlinkSync",
		"rmSync",
		"rmdirSync",
		"mkdirSync",
		"appendFileSync",
		"accessSync",
		"chmodSync",
		"chownSync",
		"lchmodSync",
		"lchownSync",
		"linkSync",
	}
)

func Validate(code string) bool {
	for _, token := range forbiddenTokens {
		if strings.Contains(code, token) {
			return false
		}
	}

	return true
}
