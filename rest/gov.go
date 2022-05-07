package rest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	utils "github.com/jim380/Cendermint/utils"
)

type govInfo struct {
	TotalProposalCount      float64
	VotingProposalCount     float64
	InVotingVotedCount      float64
	InVotingDidNotVoteCount float64
}

type gov struct {
	Proposals  []proposal `json:"proposals"`
	Pagination struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type proposal struct {
	ProposalID string `json:"proposal_id"`
	Status     string `json:"status"`
}

type voteInfo struct {
	Votes struct {
		Option string `json:"option"`
	} `json:"vote"`
}

func (rd *RESTData) getGovInfo() {
	var (
		g                  gov
		gi                 govInfo
		voteInfo           voteInfo
		proposalsInVoting  []string
		inVotingVoted      int
		inVotingDidNotVote int
	)

	res, err := HttpQuery(RESTAddr + "/cosmos/gov/v1beta1/proposals?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &g)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		zap.L().Info("\t", zap.Bool("Success", true), zap.String("Total Proposal Count", g.Pagination.Total))
	}

	for _, value := range g.Proposals {
		if value.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			proposalsInVoting = append(proposalsInVoting, value.ProposalID)
		}
	}

	fmt.Println("Proposals in voting: " + strconv.Itoa(len(proposalsInVoting)))

	for _, value := range proposalsInVoting {
		res, err := HttpQuery(RESTAddr + "/cosmos/gov/v1beta1/proposals/" + value + "/votes/" + utils.GetAccAddrFromOperAddr(OperAddr))
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		json.Unmarshal(res, &voteInfo)
		if voteInfo.Votes.Option != "" {
			inVotingVoted++
			//fmt.Println(value + ":Voter voted")
		} else {
			inVotingDidNotVote++
			//fmt.Println(value + ":Voter didn't vote")
		}
	}

	gi.TotalProposalCount = utils.StringToFloat64(g.Pagination.Total)
	gi.VotingProposalCount = float64(len(proposalsInVoting))
	gi.InVotingVotedCount = float64(inVotingVoted)
	gi.InVotingDidNotVoteCount = float64(inVotingDidNotVote)

	rd.Gov = gi
}
