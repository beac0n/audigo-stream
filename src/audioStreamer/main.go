package audioStreamer

import (
	"audigo-stream/src/commands"
	"audigo-stream/src/util"
	"errors"
	"github.com/sclevine/agouti"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mux = &sync.Mutex{}

type AudioStreamer struct {
	driver         *agouti.WebDriver
	page           *agouti.Page
	nullSinkName   string
	sinkInputIndex int64
}

func (audioStreamer *AudioStreamer) Cleanup() {
	if audioStreamer.driver == nil {
		return
	}

	if err := audioStreamer.driver.Stop(); err != nil {
		log.Println("ERROR - cleanup audio stream:", err)
	}

	audioStreamer.removeSinkInput()
}

func (audioStreamer *AudioStreamer) Create(url string) (*io.PipeReader, error) {
	// see https://askubuntu.com/questions/60837/record-a-programs-output-with-pulseaudio

	if err := audioStreamer.initDriver(); err != nil {
		return nil, err
	}

	if err := audioStreamer.setSinkInputIndex(url); err != nil {
		return nil, err
	}

	audioStreamReader, err := audioStreamer.getAudioStreamReader()
	if err != nil {
		return nil, err
	}

	return audioStreamReader, nil
}

func (audioStreamer *AudioStreamer) removeSinkInput() {
	if audioStreamer.nullSinkName == "" {
		return
	}

	output, err := commands.ListSinks()
	if err != nil {
		log.Println("ERROR - removeSinkInput list sinks:", err)
		return
	}

	split := strings.Split(output, "\n")

	var ownerModule string
	for index, value := range split {
		if strings.Contains(value, audioStreamer.nullSinkName) {
			ownerModule = strings.Split(split[index+5], ":")[1]
			ownerModule = strings.Trim(ownerModule, " ")
			break
		}
	}

	if _, err := commands.UnloadModule(ownerModule); err != nil {
		log.Println("ERROR - removeSinkInput unload-module ", ownerModule, ": ", err)
	}
}

func (audioStreamer *AudioStreamer) initDriver() error {
	options := []string{"--headless", "--disable-gpu", "--no-first-run", "--no-default-browser-check", "--no-sandbox"}
	chromeOptions := agouti.ChromeOptions("args", options)
	chromeDriver := agouti.ChromeDriver(chromeOptions, agouti.Debug)

	if err := chromeDriver.Start(); err != nil {
		return err
	}

	audioStreamer.driver = chromeDriver

	return nil
}

func (audioStreamer *AudioStreamer) getAudioStreamReader() (*io.PipeReader, error) {
	nullSinkName := "audiostreamer" + strconv.FormatInt(audioStreamer.sinkInputIndex, 10)

	if _, err := commands.CreateNullSink(nullSinkName); err != nil {
		return nil, err
	}

	if _, err := commands.MoveNullSink(audioStreamer.sinkInputIndex, nullSinkName); err != nil {
		return nil, err
	}

	parecCmd := commands.CreateNullSinkStreamCommand(nullSinkName)
	audioCmd := commands.CreateAudioStreamToMp3Command()

	readerAudio, writerParec := io.Pipe()
	readerWeb, writerAudio := io.Pipe()

	parecCmd.Stdout = writerParec
	audioCmd.Stdin = readerAudio
	audioCmd.Stdout = writerAudio

	if err := parecCmd.Start(); err != nil {
		return nil, err
	}

	if err := audioCmd.Start(); err != nil {
		return nil, err
	}

	audioStreamer.nullSinkName = nullSinkName

	return readerWeb, nil
}

func (audioStreamer *AudioStreamer) setSinkInputIndex(url string) error {
	defer mux.Unlock()
	mux.Lock()

	indexesBefore, errBefore := util.GetPulseAudioInputSinkIndexes()
	if errBefore != nil {
		return errBefore
	}

	if err := audioStreamer.goToPage(url); err != nil {
		return err
	}

	if err := audioStreamer.sendKeys(" "); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	indexesAfter, errAfter := util.GetPulseAudioInputSinkIndexes()
	if errAfter != nil {
		return errAfter
	}

	indexesDiff := util.Difference(indexesBefore, indexesAfter)
	if len(indexesDiff) == 0 || len(indexesDiff) > 1 {
		return errors.New("wrong sinks diff")
	}

	audioStreamer.sinkInputIndex = indexesDiff[0]

	return nil
}

func (audioStreamer *AudioStreamer) sendKeys(keys string) error {
	keystrokes := map[string][]string{"value": {keys}}
	if err := audioStreamer.page.Session().Send("POST", "keys", keystrokes, nil); err != nil {
		return err
	}

	return nil
}

func (audioStreamer *AudioStreamer) goToPage(url string) error {
	page, err := audioStreamer.driver.NewPage(agouti.Browser("chromium"))
	if err != nil {
		return err
	}

	if err = page.Navigate(url); err != nil {
		return err
	}

	if err = page.SetImplicitWait(1000); err != nil {
		return err
	}

	audioStreamer.page = page

	time.Sleep(500 * time.Millisecond)

	return nil
}
