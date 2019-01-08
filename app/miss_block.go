package app

import (
	"encoding/json"
	"fmt"
	"validator-monitor/app/config"
	"validator-monitor/utils"

	"github.com/astaxie/beego"
)

type Block struct {
	Height int64 `json:"height"`
	Miss   bool  `json:"miss"`
}

type MissBlocks struct {
	Blocks    []*Block `json:"blocks"`
	missCount int      `json:"count"`
	Address   string   `json:"address"`
}

var defaultMissBlocks *MissBlocks

func init() {
	defaultMissBlocks = NewMissBlocks(config.MissBlockWindowSize, config.ValidatorAddresses[0])
}

func NewMissBlocks(size int, address string) *MissBlocks {
	blocks := &MissBlocks{Blocks: make([]*Block, size, size)}
	for i := 0; i < size; i++ {
		blocks.Blocks[i] = &Block{
			Height: 0,
			Miss:   false,
		}
	}
	blocks.missCount = 0
	blocks.Address = address

	return blocks
}

func (mbs *MissBlocks) Size() int64 {
	return int64(len(mbs.Blocks))
}

func (mbs *MissBlocks) Reset() {
	for _, v := range mbs.Blocks {
		v.Height = 0
		v.Miss = false
	}
	mbs.missCount = 0
}

func (mbs *MissBlocks) MissCount() int {
	return mbs.missCount
}

func (mbs *MissBlocks) String() string {
	bin, err := json.Marshal(mbs)
	if nil != err {
		beego.Error(err)
		return err.Error()
	}

	return string(bin)
}

func (mbs *MissBlocks) SetMiss(height int64) {
	index := height % mbs.Size()
	if 0 == index {
		mbs.Reset()
	}

	if !mbs.Blocks[index].Miss {
		mbs.missCount++
	}
	mbs.Blocks[index] = &Block{Height: height, Miss: true}

	if mbs.missCount >= config.MissBlockLimit {
		emailBody := fmt.Sprintf("validator(%s) miss block %d/%d,details:%s", mbs.Address, mbs.missCount, config.MissBlockWindowSize, mbs.String())
		beego.Error(emailBody)
		utils.SendMail("validator miss block", emailBody)
		utils.SendSMS(emailBody)
	}
}
