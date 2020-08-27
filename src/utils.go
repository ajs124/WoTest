package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
)

type EnvEntry struct {
	key   string
	value string
}

type Runtime int8

// When modifying this, also update config.d
const (
	Node = iota
	Python
	Java
)

// Execute the command `name` in the folder `dir` while extending the environment with `env` and the arguments `args`
// returns the cmd struct and error from exec.Start()
func StartCmd(name, dir string, env []EnvEntry, args ...string) (*exec.Cmd, error) {
	var envStrings []string
	for _, v := range env {
		envStrings = append(envStrings, v.key+"="+v.value)
	}
	// create context
	cmd := exec.Command(name, args...)
	// set up more stuff
	cmd.Env = append(os.Environ(), envStrings...)
	cmd.Dir = dir

	err := cmd.Start()
	if err != nil {
		log.Debug().Err(err).Str("Name", name).Str("Directory", dir).Strs("Environment", envStrings).Strs("args", args)
	}
	return cmd, err
}

// Same as StartCmd, but runs node and sets NODE_PATH
func StartNode(nodePath, dir string, args ...string) (*exec.Cmd, error) {
	return StartCmd("node", dir, []EnvEntry{{"NODE_PATH", nodePath}}, args...)
}

// Same as StartCmd, but runs python and sets PYTHON_PATH
func StartPython(pythonPath, dir string, args ...string) (*exec.Cmd, error) {
	return StartCmd("python", dir, []EnvEntry{{"PYTHONPATH", pythonPath}}, args...)
}

// TODO: pass on classPath
func StartJava(classPath, dir string, args ...string) (*exec.Cmd, error) {
	return StartCmd("java", dir, []EnvEntry{}, args...)
}
