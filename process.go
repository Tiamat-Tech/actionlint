package actionlint

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/mattn/go-shellwords"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"golang.org/x/sys/execabs"
)

// cmdExecution represents a single command line execution.
type cmdExecution struct {
	cmd           string
	args          []string
	stdin         string
	combineOutput bool
}

func (e *cmdExecution) run() ([]byte, error) {
	cmd := exec.Command(e.cmd, e.args...)
	cmd.Stderr = nil

	p, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("could not make stdin pipe for %s process: %w", e.cmd, err)
	}
	if _, err := io.WriteString(p, e.stdin); err != nil {
		p.Close()
		return nil, fmt.Errorf("could not write to stdin of %s process: %w", e.cmd, err)
	}
	p.Close()

	var stdout []byte
	if e.combineOutput {
		stdout, err = cmd.CombinedOutput()
	} else {
		stdout, err = cmd.Output()
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()

			stderr := exitErr.Stderr
			if e.combineOutput {
				stderr = stdout
			}

			if code < 0 {
				return nil, fmt.Errorf("%s was terminated. stderr: %q", e.cmd, stderr)
			}

			if len(stdout) == 0 {
				return nil, fmt.Errorf("%s exited with status %d but stdout was empty. stderr: %q", e.cmd, code, stderr)
			}

			// Reaches here when exit status is non-zero and stdout is not empty, shellcheck successfully found some errors
		} else {
			return nil, err
		}
	}

	return stdout, nil
}

// concurrentProcess is a manager to run process concurrently. Since running process consumes OS
// resources, running too many processes concurrently causes some issues. On macOS, making too many
// process makes the parent process hang (see issue #3). And running processes which open files can
// cause the error "pipe: too many files to open". To avoid it, this type manages how many processes
// are run at once.
type concurrentProcess struct {
	ctx  context.Context
	sema *semaphore.Weighted
	wg   sync.WaitGroup
}

// newConcurrentProcess creates a new ConcurrentProcess instance. The `par` argument represents how
// many processes can be run in parallel. It is recommended to use the value returned from
// runtime.NumCPU() for the argument.
func newConcurrentProcess(par int) *concurrentProcess {
	return &concurrentProcess{
		ctx:  context.Background(),
		sema: semaphore.NewWeighted(int64(par)),
	}
}

func (proc *concurrentProcess) run(eg *errgroup.Group, exec *cmdExecution, callback func([]byte, error) error) {
	proc.wg.Add(1)
	eg.Go(func() error {
		defer proc.wg.Done()
		if err := proc.sema.Acquire(proc.ctx, 1); err != nil {
			return fmt.Errorf("could not acquire semaphore to run %q: %w", exec.cmd, err)
		}
		stdout, err := exec.run()
		proc.sema.Release(1)
		return callback(stdout, err)
	})
}

// wait waits all goroutines started by this concurrentProcess instance finish.
func (proc *concurrentProcess) wait() {
	proc.wg.Wait() // Wait for all goroutines completing to shutdown
}

// newCommandRunner creates new external command runner for given executable. The executable path
// is resolved in this function.
func (proc *concurrentProcess) newCommandRunner(exe string, combineOutput bool) (*externalCommand, error) {
	var args []string
	p, args, err := resolveExternalCommand(exe)
	if err != nil {
		return nil, err
	}
	cmd := &externalCommand{
		proc:          proc,
		exe:           p,
		args:          args,
		combineOutput: combineOutput,
	}
	return cmd, nil
}

func resolveExternalCommand(exe string) (string, []string, error) {
	c, err := execabs.LookPath(exe)
	if err == nil {
		return c, nil, nil
	}

	// Try to parse the string as a command line instead of a single executable file path.
	if a, err := shellwords.Parse(exe); err == nil && len(a) > 0 {
		if c, err := execabs.LookPath(a[0]); err == nil {
			return c, a[1:], nil
		}
	}

	return "", nil, err
}

// externalCommand is struct to run specific command concurrently with concurrentProcess bounding
// number of processes at the same time. This type manages fatal errors while running the command
// by using errgroup.Group. The wait() method must be called at the end for checking if some fatal
// error occurred.
type externalCommand struct {
	proc          *concurrentProcess
	eg            errgroup.Group
	exe           string
	args          []string
	combineOutput bool
}

// run runs the command with given arguments and stdin. The callback function is called after the
// process runs. First argument is stdout and the second argument is an error while running the
// process.
func (cmd *externalCommand) run(args []string, stdin string, callback func([]byte, error) error) {
	if len(cmd.args) > 0 {
		var allArgs []string
		allArgs = append(allArgs, cmd.args...)
		allArgs = append(allArgs, args...)
		args = allArgs
	}
	exec := &cmdExecution{cmd.exe, args, stdin, cmd.combineOutput}
	cmd.proc.run(&cmd.eg, exec, callback)
}

// wait waits until all goroutines for this command finish. Note that it does not wait for
// goroutines for other commands.
func (cmd *externalCommand) wait() error {
	return cmd.eg.Wait()
}
