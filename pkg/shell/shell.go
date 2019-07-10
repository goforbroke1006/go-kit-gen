package shell

import (
	"io/ioutil"
	"os/exec"
)

func Execute(bin string, args ...string) (output string, err error) {
	command := exec.Command(bin, args...)
	stdout, err := command.StdoutPipe()
	err = command.Start()
	if nil != err {
		return "", err
	}

	//log.Printf("Waiting for command to finish...")
	//err = command.Wait()
	//if nil != err {
	//	return "", err
	//}

	out, err := ioutil.ReadAll(stdout)
	return string(out), err

	//if cmd, e := exec.Run("/bin/ls", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
	//	b, _ := ioutil.ReadAll(cmd.Stdout)
	//	println("output: " + string(b))
	//}
}
