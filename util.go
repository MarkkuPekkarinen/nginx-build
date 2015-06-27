package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func runCommand(cmd *exec.Cmd) error {
	checkVerboseEnabled(cmd)
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func saveCurrentDir() string {
	prevDir, _ := filepath.Abs(".")
	return prevDir
}

func restoreCurrentDir(prevDir string) {
	os.Chdir(prevDir)
}

func printFirstMsg() {
	fmt.Printf(`nginx-build: %s
Compiler: %s %s
`,
		NGINX_BUILD_VERSION,
		runtime.Compiler,
		runtime.Version())
}

func printLastMsg(workDir, srcDir string, openResty, configureOnly bool) {
	log.Println("Complete building nginx!")

	if !openResty {
		if !configureOnly {
			fmt.Println()
			err := printConfigureOptions()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	fmt.Println()

	lastMsgFormat := `Enter the following command for install nginx.

   $ cd %s/%s%s
   $ sudo make install
`
	if configureOnly {
		log.Printf(lastMsgFormat, workDir, srcDir, "\n   $ make")
	} else {
		log.Printf(lastMsgFormat, workDir, srcDir, "")
	}
}

func versionCheck(version string) {
	if len(version) == 0 {
		log.Println("[warn]nginx version is not set.")
		log.Printf("[warn]nginx-build use %s.\n", NGINX_VERSION)
	}
}

func fileGetContents(path string) (string, error) {
	conf := ""
	if len(path) > 0 {
		confb, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("confPath(%s) does not exist.", path)
		}
		conf = string(confb)
	}
	return conf, nil
}

func configureNginx() error {
	if VerboseEnabled {
		return runCommand(exec.Command("sh", "./nginx-configure"))
	}

	f, err := os.Create("nginx-configure.log")
	if err != nil {
		return runCommand(exec.Command("sh", "./nginx-configure"))
	}
	defer f.Close()

	cmd := exec.Command("sh", "./nginx-configure")
	writer := bufio.NewWriter(f)
	cmd.Stdout = writer
	defer writer.Flush()

	return cmd.Run()
}

func buildNginx(jobs int) error {
	if VerboseEnabled {
		return runCommand(exec.Command("make", "-j", strconv.Itoa(jobs)))
	}

	f, err := os.Create("nginx-build.log")
	if err != nil {
		return runCommand(exec.Command("make", "-j", strconv.Itoa(jobs)))
	}
	defer f.Close()

	cmd := exec.Command("make", "-j", strconv.Itoa(jobs))
	writer := bufio.NewWriter(f)
	cmd.Stderr = writer
	defer writer.Flush()

	return cmd.Run()
}

func extractArchive(path string) error {
	return runCommand(exec.Command("tar", "zxvf", path))
}

func switchRev(rev string) error {
	return runCommand(exec.Command("git", "checkout", rev))
}

func printConfigureOptions() error {
	cmd := exec.Command("objs/nginx", "-V")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func provideShell(sh string) error {
	if len(sh) == 0 {
		return nil
	}
	args := strings.Split(strings.Trim(sh, " "), " ")
	var err error
	if len(args) == 1 {
		err = runCommand(exec.Command(args[0]))
	} else {
		err = runCommand(exec.Command(args[0], args[1:]...))
	}
	return err
}

func normalizeConfigure(configure string) string {
	configure = strings.TrimRight(configure, "\n")
	configure = strings.TrimRight(configure, " ")
	configure = strings.TrimRight(configure, "\\")
	if configure != "" {
		configure += " "
	}
	return configure
}

func clearWorkDir(workDir string) error {
	err := os.RemoveAll(workDir)
	if err != nil {
		// workaround for a restriction of os.RemoveAll
		// os.RemoveAll call fd.Readdirnames(100).
		// So os.RemoveAll does not always remove all entries.
		// Some 3rd-party module(e.g. lua-nginx-module) tumbles this restriction.
		if fileExists(workDir) {
			err = os.RemoveAll(workDir)
		}
	}
	return err
}

func normalizeAddModulePaths(path, rootDir string) string {
	var result string
	if len(path) == 0 {
		return path
	}

	module_paths := strings.Split(path, ",")

	for _, module_path := range module_paths {
		if strings.HasPrefix(module_path, "/") {
			result += fmt.Sprintf("--add-module=%s \\\n", module_path)
		} else {
			result += fmt.Sprintf("--add-module=%s/%s \\\n", rootDir, module_path)
		}
	}

	return result
}

func fatalLog(err error, path string) {
	if VerboseEnabled {
		log.Fatal(err)
	}

	f, err2 := os.Open(path)
	if err2 != nil {
		log.Printf("error-log: %s is not found\n", path)
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		os.Stderr.Write(scanner.Bytes())
		os.Stderr.Write([]byte("\n"))
	}

	log.Fatal(err)
}
