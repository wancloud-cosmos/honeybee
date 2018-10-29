package controllers

import (
	"fmt"
	"validator-monitor/app"
	"validator-monitor/app/http"

	"github.com/astaxie/beego"
)

type VoteController struct {
	beego.Controller
}

func (c *VoteController) Get() {
	desc := c.GetString("desc")
	id, err := c.GetInt64("id")
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}
	title := c.GetString("title")

	v, err := http.QueryVote(id, app.GovVoter)
	if nil == err {
		switch v.Ooption {
		case "No":
			c.Data["No"] = true
		case "Yes":
			c.Data["Yes"] = true
		case "Abstain":
			c.Data["Abstain"] = true
		case "NoWithVeto":
			c.Data["NoWithVeto"] = true
		}
	}

	c.Data["ProposalID"] = id
	c.Data["Title"] = title
	c.Data["Description"] = desc
	c.TplName = "vote.tpl"
}

func (c *VoteController) Post() {
	option := c.GetString("optionsRadios")
	id, err := c.GetInt64("id")
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}
	beego.Debug("vote ", option, "for:", id)

	opt, err := http.VoteOptionFromString(option)
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}

	err = http.Vote(id, app.GovVoter, opt)
	if nil != err {
		beego.Error(err)
		c.Data["Error"] = err.Error()
		c.TplName = "error.tpl"

		return
	}

	resp := fmt.Sprintf("成功给ID为%d的提议投了%s票", id, option)
	c.Ctx.WriteString(resp)
}
