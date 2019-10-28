package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 定义 IDGenerator 实例变量
var id *IDGenerator

// 获取 IDGenerator 单例
func GetIDInstance() *IDGenerator {
	if id == nil {
		id = &IDGenerator{}
		id.idGenInit()
	}
	return id
}

// IDGenerator 结构体
type IDGenerator struct {
	mu            sync.Mutex
	Epoch         int64
	MachineId     int
	MachineBit    int
	Sequence      int
	SequenceBit   int
	lastTimestamp int64
}

// 初始化 IDGenerator 结构体变量值，创建 IDGenerator 实例时调用的方法，相当于构造函数
func (ID *IDGenerator) idGenInit() {
	ID.Epoch = 1500000000000
	ID.MachineId = 0
	ID.MachineBit = 4
	ID.Sequence = 0
	ID.SequenceBit = 10
	ID.lastTimestamp = 0
}

// 生成一个长整型 ID
func (ID *IDGenerator) NextID() int64 {
	// 互斥锁，确保同一时间只能有一个线程进入
	ID.mu.Lock()

	timestamp := ID.timeGen()
	if ID.lastTimestamp > timestamp { // 判断上次时间戳大于当前时间戳，防止时钟回拨
		timestamp = ID.tilNextMillis(ID.lastTimestamp)
	} else if ID.lastTimestamp == timestamp { // 同一时间戳内序列号自增
		ID.Sequence = (ID.Sequence + 1) & (pow(2, ID.SequenceBit) - 1)
		if ID.Sequence == 0 { // 自增序列号超过最大值时，等待到下一毫秒
			timestamp = ID.tilNextMillis(ID.lastTimestamp)
		}
	} else { // 上次时间戳小于当前时间戳则序列号回归为0
		ID.Sequence = 0
	}

	var id int64 = 0
	ID.lastTimestamp = timestamp
	if timestamp > 0 { // 等待下一毫秒执行次数过多则返回 ID = 0
		timeBits := decBin(timestamp - ID.Epoch)
		machineBits := fmt.Sprintf("%0"+strconv.Itoa(ID.MachineBit)+"s", decBin(int64(ID.MachineId)))
		sequenceBits := fmt.Sprintf("%0"+strconv.Itoa(ID.SequenceBit)+"s", decBin(int64(ID.Sequence)))
		// 生成 ID
		id = binDec(timeBits + machineBits + sequenceBits)
	}

	// 解除互斥锁
	ID.mu.Unlock()

	return id
}

// 获取当前毫秒级时间戳
func (ID *IDGenerator) timeGen() int64 {
	return time.Now().UnixNano() / 1000000
}

// 获取下一毫秒的毫秒级时间戳（相对于lastTimestamp）
func (ID *IDGenerator) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := ID.timeGen()
	count := 0
	for lastTimestamp > timestamp {
		count++
		// 只执行100次
		if count > 100 {
			return 0
		}
		time.Sleep(time.Duration(1) * time.Millisecond)
		timestamp = ID.timeGen()
	}
	return timestamp
}

// 幂次方值计算
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

// 十进制转二进制
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

// 二进制转十进制
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
