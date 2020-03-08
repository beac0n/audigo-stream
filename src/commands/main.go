package commands

import (
	"bytes"
	"os/exec"
	"strconv"
)

var pactl = "pactl"
var parec = "parec"
var lame = "lame"
var pacmd = "pacmd"

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var outBytes bytes.Buffer

	cmd.Stdout = &outBytes

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return outBytes.String(), nil
}

func CreateNullSinkStreamCommand(nullSinkName string) *exec.Cmd {
	return exec.Command(parec, "--format=s16le", "-d", nullSinkName+".monitor")
}

func CreateAudioStreamToMp3Command() *exec.Cmd {
	return exec.Command(lame, "-r", "--quiet", "-q", "3", "--lowpass", "17", "--abr", "192", "-", "-")
}

func CreateNullSink(nullSinkName string) (string, error) {
	return RunCommand(pactl, "load-module", "module-null-sink", "sink_name="+nullSinkName)
}

func MoveNullSink(sinkInputIndex int64, nullSinkName string) (string, error) {
	return RunCommand(pactl, "move-sink-input", strconv.FormatInt(sinkInputIndex, 10), nullSinkName)
}

func ListSinks() (string, error) {
	return RunCommand(pactl, "list", "sinks")
}

func UnloadModule(ownerModule string) (string, error) {
	return RunCommand(pactl, "unload-module", ownerModule)
}

func ListSinkInputs() (string, error) {
	return RunCommand(pacmd, "list-sink-inputs")
}
