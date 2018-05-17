package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/ami-GS/cuimeter"
)

type Channel struct {
	Reads  float64
	Writes float64
}

type Socket struct {
	Channels []Channel
	// Node?
}

type System struct {
	Sockets []Socket
	Read    float64
	Write   float64
	Memory  float64
}

type PCMMemoryHint struct {
	*cuimeter.BaseHint
	SystemData System
}

func NewPCMMemoryHint(unit string) *PCMMemoryHint {
	return &PCMMemoryHint{
		// interval is not needed for pipe
		BaseHint: cuimeter.NewBaseHint(unit, 0),
		// dynamically decided
		SystemData: System{nil, 0, 0, 0},
	}
}

var firstRead = true

func (s *PCMMemoryHint) lateInit(data *string) {
	numSocket := strings.Count(*data, "Socket")
	numChannel := strings.Count(*data, "Ch") / numSocket
	s.SystemData.Sockets = make([]Socket, numSocket)
	for i := 0; i < numSocket; i++ {
		s.SystemData.Sockets[i].Channels = make([]Channel, numChannel)
	}
}

// This should be changed based on platform configuration
// or automatically configured at firstRead?
const onePipeMessageSize = 1974

func (s *PCMMemoryHint) read() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	// read 2KB at a time
	buff := make([]byte, onePipeMessageSize)
	allLen := 0
	for allLen < onePipeMessageSize {
		length, err := reader.Read(buff[allLen:])
		if err != nil {
			return "", err
		}
		allLen += length
	}

	str := string(buff[:allLen])
	if firstRead {
		s.lateInit(&str)
		firstRead = false
	}

	return str, nil
}

func (s *PCMMemoryHint) parse(dat string) (interface{}, error) {
	lines := strings.Split(dat, "\n")
	lines = lines[5:]

	// from Mem Ch   0: ...
	memChParser := func(line string) float64 {
		content := strings.Fields(line)
		// for Mem Ch   0: Reads (MB/s):   XX.XX
		if content[0] == "Mem" {
			dat, _ := strconv.ParseFloat(content[5], 64)
			return dat
		}
		// for             Writes(MB/s):   XX.XX
		dat, _ := strconv.ParseFloat(content[1], 64)
		return dat
	}
	stripFunc := func(line string) []string {
		// from "|-- content A --||-- content B --||-- content C --|"
		// to ["content A", "content B", "content C"]
		return strings.Split(
			strings.Trim(strings.Trim(line, "|-- "), " --|"),
			" --||-- ")
	}

	chCounter := 0
	for i := 0; i < len(lines); {
		if strings.HasSuffix(lines[i], "-------|") || len(lines[i]) == 0 {
			i++
		} else if strings.HasPrefix(lines[i], "|-- Mem Ch") {
			for j := 0; j < 2; j++ { // for Reads and Writes
				sockets := stripFunc(lines[i+j])
				for k, socket := range sockets {
					if j == 0 {
						s.SystemData.Sockets[k].Channels[chCounter].Reads = memChParser(socket)
					} else {
						s.SystemData.Sockets[k].Channels[chCounter].Writes = memChParser(socket)
					}
				}
			}
			i += 2
			chCounter++
		} else if strings.HasPrefix(lines[i], "|-- NODE ") {
			for j := 0; j < 4; j++ {
				_ = stripFunc(lines[i+j])
				// ignore as of now
			}
			i += 4
		} else { // for system data
			dat := strings.Fields(lines[i])
			s.SystemData.Read, _ = strconv.ParseFloat(dat[4], 64)
			dat = strings.Fields(lines[i+1])
			s.SystemData.Read, _ = strconv.ParseFloat(dat[4], 64)
			dat = strings.Fields(lines[i+2])
			s.SystemData.Memory, _ = strconv.ParseFloat(dat[4], 64)
			i += 3
		}
	}
	return []float64{s.SystemData.Sockets[0].Channels[2].Reads, s.SystemData.Sockets[1].Channels[2].Reads}, nil
}

func memorystatus() {
	hint := NewPCMMemoryHint("MB")
	graph := cuimeter.NewGraph([]string{"Socket0_Ch1_Read", "Socket1_Ch1_Read"})
	graph.RunWithPipe(hint)
}

func main() {
	memorystatus()
}
