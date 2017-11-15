package command

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gobs/args"
	shellquote "github.com/kballard/go-shellquote"
)

var errorTimeout = errors.New("error: execution timeout")

type Result struct {
	RealTime time.Duration
	UserTime time.Duration
	SysTime  time.Duration
	Stdout   *bytes.Buffer
	Stderr   *bytes.Buffer
	Pid      int
	ExitCode int
	Failed   bool
	User     string
	Rusage   *syscall.Rusage
}

type Command struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Pid    int

	result    *Result
	cmd       *exec.Cmd
	command   string
	startTime time.Time
	endTime   time.Time
	failed    bool
	exitCode  int
	params    struct {
		user        string
		timeout     time.Duration
		workingDir  string
		environment []string
	}
}

func Escape(command string, args ...string) string {
	for _, arg := range args {
		command = shellquote.Join(command, arg)
	}
	return command
}

// New returns the Command struct to execute the named program with
// the given arguments.
func New(command string) *Command {
	return &Command{
		Stdout:  bytes.NewBuffer(nil),
		Stderr:  bytes.NewBuffer(nil),
		command: command,
	}
}

func (c *Command) SetTimeout(timeout time.Duration) {
	c.params.timeout = timeout
}

func (c *Command) SetUser(username string) {
	c.params.user = username
}

func (c *Command) SetWorkingDir(workingDir string) {
	c.params.workingDir = workingDir
}

func (c *Command) SetEnvironment(environment []string) {
	c.params.environment = environment
}

func (c *Command) Start() error {
	c.buildExecCmd()
	c.setOutput()
	if err := c.setCredentials(); err != nil {
		return err
	}

	if err := c.cmd.Start(); err != nil {
		return err
	}

	c.startTime = time.Now()
	c.Pid = c.cmd.Process.Pid

	return nil
}

func (c *Command) buildExecCmd() {
	arguments := args.GetArgs(c.command)
	aname, err := exec.LookPath(arguments[0])
	if err != nil {
		aname = arguments[0]
	}

	c.cmd = &exec.Cmd{
		Path: aname,
		Args: arguments,
	}

	if c.params.workingDir != "" {
		c.cmd.Dir = c.params.workingDir
	}

	if c.params.environment != nil {
		c.cmd.Env = c.params.environment
	}
}

func (c *Command) setOutput() {
	c.cmd.Stdout = c.Stdout
	c.cmd.Stderr = c.Stderr
}

func (c *Command) setCredentials() error {
	if c.params.user == "" {
		return nil
	}

	uid, gid, err := c.getUIDAndGIDInfo(c.params.user)
	if err != nil {
		return err
	}

	c.cmd.SysProcAttr = &syscall.SysProcAttr{}
	c.cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}

	return nil
}

func (c *Command) getUIDAndGIDInfo(username string) (uint32, uint32, error) {
	user, err := user.Lookup(username)
	if err != nil {
		return 0, 0, err
	}

	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	return uint32(uid), uint32(gid), nil
}

func (c *Command) Wait() error {
	if err := c.doWait(); err != nil {
		c.failed = true

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				c.exitCode = status.ExitStatus()
			}
		} else {
			if err != errorTimeout {
				return err
			}

			c.Kill()
		}
	}

	c.endTime = time.Now()
	c.buildResponse()

	return nil
}

func (c *Command) Kill() error {
	c.failed = true
	c.exitCode = -1

	return c.cmd.Process.Kill()
}

func (c *Command) doWait() error {
	if c.params.timeout != 0 {
		return c.doWaitWithTimeout()
	}

	return c.doWaitWithoutTimeout()
}

func (c *Command) doWaitWithoutTimeout() error {
	return c.cmd.Wait()
}

func (c *Command) doWaitWithTimeout() error {
	go func() {
		time.Sleep(c.params.timeout)
		c.Kill()
	}()

	return c.cmd.Wait()
}

func (c *Command) buildResponse() {
	result := &Result{
		RealTime: c.endTime.Sub(c.startTime),
		UserTime: c.cmd.ProcessState.UserTime(),
		SysTime:  c.cmd.ProcessState.UserTime(),
		Rusage:   c.cmd.ProcessState.SysUsage().(*syscall.Rusage),
		Stdout:   c.Stdout,
		Stderr:   c.Stderr,
		Pid:      c.cmd.Process.Pid,
		Failed:   c.failed,
		ExitCode: c.exitCode,
		User:     c.params.user,
	}

	c.result = result
}

func (c *Command) Result() *Result {
	return c.result
}

func (c *Command) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (r *Result) StdoutString() string {
	return strings.TrimSuffix(string(r.Stdout.Bytes()), "\n")
}

func (r *Result) StderrString() string {
	return strings.TrimSuffix(string(r.Stderr.Bytes()), "\n")
}

func (c *Command) RunWithTTY() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", c.command)
	} else {
		cmd = exec.Command("sh", "-c", c.command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Join(c string, args []string) string {
	for _, arg := range args {
		c = c + " " + arg
	}
	return c
}
