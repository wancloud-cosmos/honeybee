package controllers

import (
	"validator-monitor/app/gov"
	"validator-monitor/app/http"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	allProposals, err := http.QueryProposals("VotingPeriod", "", "")
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}

	votedProposals, err := http.QueryProposals("VotingPeriod", "", gov.GovVoter)
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}

	var proposals = []*http.Proposal{}

	for _, v := range allProposals {
		if !IsVoted(v.ProposalID, votedProposals) {
			proposals = append(proposals, v)
		}
	}

	if 0 == len(proposals) {
		c.Data["Error"] = "没有正在投票状态的提议"
		c.TplName = "error.tpl"

		return
	}

	c.Data["Proposals"] = proposals
	for _, v := range proposals {
		beego.Debug(v)
	}
	c.TplName = "index.tpl"
}

func IsVoted(id uint64, ps []*http.Proposal) bool {
	for _, v := range ps {
		if id == v.ProposalID {
			return true
		}
	}

	return false
}
