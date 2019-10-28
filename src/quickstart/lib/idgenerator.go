package lib

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var id *IDGenerator

func GetIDInstance() *IDGenerator {
	if id == nil {
		id = &IDGenerator{}
		id.idGenInit()
	}
	return id
}

type IDGenerator struct {
	mu            sync.Mutex
	Epoch         int64
	MachineId     int
	MachineBit    int
	Sequence      int
	SequenceBit   int
	lastTimestamp int64
}

func init() {

}

func (ID *IDGenerator) idGenInit() {
	ID.Epoch = 1514736000000
	ID.MachineId = 0
	ID.MachineBit = 4
	ID.Sequence = 0
	ID.SequenceBit = 10
	ID.lastTimestamp = 0
}

func (ID *IDGenerator) Get() int64 {

	ID.mu.Lock()

	timestamp := ID.timeGen()

	//fmt.Printf("Timestamp: %d\n", timestamp)
	//fmt.Printf("LastTimestamp: %d\n", ID.lastTimestamp)

	if ID.lastTimestamp > timestamp {
		timestamp = ID.tilNextMillis(ID.lastTimestamp)
	} else if ID.lastTimestamp == timestamp {
		ID.Sequence = (ID.Sequence + 1) & (pow(2, ID.SequenceBit) - 1)
		if ID.Sequence == 0 {
			timestamp = ID.tilNextMillis(ID.lastTimestamp)
		}
	} else {
		ID.Sequence = 0
	}

	ID.lastTimestamp = timestamp

	ID.mu.Unlock()

	timeBits := decBin(timestamp - ID.Epoch)
	machineBits := fmt.Sprintf("%0"+strconv.Itoa(ID.MachineBit)+"s", decBin(int64(ID.MachineId)))
	sequenceBits := fmt.Sprintf("%0"+strconv.Itoa(ID.SequenceBit)+"s", decBin(int64(ID.Sequence)))
	//fmt.Printf("sequence: %d\n", ID.Sequence)
	//fmt.Printf("sequenceBits: %s\n", sequenceBits)
	id := binDec(timeBits + machineBits + sequenceBits)
	//fmt.Printf("sequenceBits: %d\n", id)
	return id
}

func (ID *IDGenerator) timeGen() int64 {
	return time.Now().UnixNano() / 1000000
}

func (ID *IDGenerator) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := ID.timeGen()
	for lastTimestamp > timestamp {
		time.Sleep(time.Duration(1) * time.Millisecond)
		timestamp = ID.timeGen()
	}
	return timestamp
}

func pow(x, n int) int {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}

func decBin(n int64) string {
	if n < 0 {
		log.Println("Decimal to binary error: the argument must be greater than zero.")
		return ""
	}
	if n == 0 {
		return "0"
	}
	s := ""
	for q := n; q > 0; q = q / 2 {
		m := q % 2
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

func binDec(b string) (n int64) {
	s := strings.Split(b, "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println("Binary to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(2, float64(l-i-1))
	}
	return int64(d)
}
