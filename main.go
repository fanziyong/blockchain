package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"math/rand"
	"strings"
	"time"
)

type Block struct {
	Timestamp     int64  // 区块生成时间
	Data          []byte // 区块数据
	PrevBlockHash []byte // 前区块哈希
	Hash          []byte // 当前块哈希
	Nonce         int32  // 随机数
}

type BlockChain struct {
	blocks       []*Block // 区块集合
	currentBlock []byte   // 当前区块
}

func NewBlock(data []byte, prevHash []byte) *Block {
	// ...生成创世区块
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          data,
		PrevBlockHash: prevHash,
		Hash:          []byte{},
		Nonce:         generaNonce(),
	}

	// 生成当前区块的哈希
	block.Hash = block.GenerateHash()

	// 工作量证明
	for !CheckProofOfWork(block) {
		block.Nonce++
	}

	return block
}

func (b *Block) GenerateHash() []byte {
	// ...生成当前区块hash

	// 1. 数据处理
	data := cast.ToString(b.Timestamp) + cast.ToString(b.Nonce) + string(b.Data) + string(b.PrevBlockHash)

	// 2. SHA256哈希
	hash := sha256.Sum256([]byte(data))

	// 3. 存储哈希
	return hash[:]
}

func (bc *BlockChain) AddBlock(data []byte) {
	// ...添加区块
	prevBlock := bc.blocks[len(bc.blocks)-1]   // 上一个区块
	newBlock := NewBlock(data, prevBlock.Hash) // 新区块

	bc.blocks = append(bc.blocks, newBlock) // 区块上链
	bc.currentBlock = newBlock.Hash         // 当前区块哈希
}

func (bc *BlockChain) Valid() bool {
	// ...循环校验每个区块hash
	for _, block := range bc.blocks {
		// 校验区块哈希是否符合要求
		if !IsValidHashDifficulty(block.Hash) {
			return false
		}
	}
	return true
}

// CheckProofOfWork POW验证
// 工作量验证
func CheckProofOfWork(block *Block) bool {
	// 计算hash
	blockHash := block.GenerateHash()

	// 验证hash是否符合难度要求
	if !IsValidHashDifficulty(blockHash) {
		return false
	}
	return true
}

// IsValidHashDifficulty 校验哈希复杂度是否符合要求
// 传入哈希 判断前4位是否符合
func IsValidHashDifficulty(hash []byte) bool {
	// 目标前4位为0即符合难度要求
	difficultyPrefix := strings.Repeat("0", 4)
	return strings.HasPrefix(string(hash), difficultyPrefix)
}

// generaNonce 生成随机数
// 当前时间戳 + 随机数字
func generaNonce() int32 {
	return cast.ToInt32(time.Now().Unix()) + rand.Int31()
}

// 新建一个创世区块
// 上链后，新增2个区块
func main() {
	chain := &BlockChain{}
	chain.blocks = append(chain.blocks, NewBlock([]byte(cast.ToString("创世区块")), []byte{}))
	chain.AddBlock([]byte(cast.ToString(rand.Int())))
	chain.AddBlock([]byte(cast.ToString(rand.Int())))

	fmt.Println(json.Marshal(chain))
}
