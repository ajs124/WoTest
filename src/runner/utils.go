package runner

import (
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
	"time"
)

type EnvEntry struct {
	key   string
	value string
}

// poll pipe once every ms and write into byte buffer
func readPipe(r io.ReadCloser, ob *[]byte, ctx context.Context) {
	var err error
	ticker := time.NewTicker(1 * time.Millisecond)
	for err != io.EOF {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n := 1024
			for n == 1024 {
				buf := make([]byte, 1024)
				n, err = r.Read(buf)
				if n > 0 {
					i := 0
					for i < n {
						*ob = append(*ob, buf[i])
						i++
					}
					// log.Debug().Bytes("stdout", *ob).Int("n", n).Msg("")
				}
			}
		}
	}
}

// Execute the command `name` in the folder `dir` while extending the environment with `env` and the arguments `args`
// returns the cmd struct and error from exec.Start()
func StartCmd(name, dir string, ctx context.Context, env []EnvEntry, stdout, stderr *[]byte, args ...string) (*exec.Cmd, error) {
	var envStrings []string
	env = append(env, EnvEntry{"LC_ALL", "C"}) // localization breaks everything
	for _, v := range env {
		envStrings = append(envStrings, v.key+"="+v.value)
	}
	// create context
	cmd := exec.CommandContext(ctx, name, args...)
	// set up more stuff
	cmd.Env = append(os.Environ(), envStrings...)
	cmd.Dir = dir

	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	go readPipe(outPipe, stdout, ctx)
	go readPipe(errPipe, stderr, ctx)

	err := cmd.Start()
	if err != nil {
		log.Debug().Err(err).Str("Name", name).Str("Directory", dir).Strs("Environment", envStrings).Strs("args", args).Msg("failed to run")
	} else {
		log.Debug().Str("Name", name).Str("Directory", dir).Strs("Environment", envStrings).Strs("args", args).Msg("ran")
	}
	return cmd, err
}

// Same as StartCmd, but runs node and sets NODE_PATH
func StartNode(nodePath, dir string, ctx context.Context, stdout, stderr *[]byte, args ...string) (*exec.Cmd, error) {
	cwd, _ := os.Getwd()
	return StartCmd("node", dir, ctx, []EnvEntry{{"NODE_PATH", cwd + "/" + nodePath}}, stdout, stderr, args...)
}

// Same as StartCmd, but runs python and sets PYTHON_PATH
func StartPython(pythonPath, dir string, ctx context.Context, stdout, stderr *[]byte, args ...string) (*exec.Cmd, error) {
	cwd, _ := os.Getwd()
	return StartCmd("python", dir, ctx, []EnvEntry{{"PYTHONPATH", cwd + "/" + pythonPath}}, stdout, stderr, args...)
}

// TODO: pass on classPath
func StartJava(classPath, dir string, ctx context.Context, stdout, stderr *[]byte, args ...string) (*exec.Cmd, error) {
	return StartCmd("java", dir, ctx, []EnvEntry{}, stdout, stderr, args...)
}
