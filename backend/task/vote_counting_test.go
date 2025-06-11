package task_test

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"

	"powervoting-server/constant"
	"powervoting-server/model"
	"powervoting-server/task"
	"powervoting-server/utils"
)

var _ = Describe("VoteCounting", func() {
	var voteCount *task.VoteCount
	var err error
	BeforeEach(func() {
		voteCount = task.NewVoteCount(mockISyncService, mockContractRepo, mockDrand)
	})

	Describe("VoteCounting", func() {
		When("Voting algorithm mismatch", func() {
			BeforeEach(func() {
				mockContractRepo.EXPECT().
					GetVotedAlgorithm().
					Return("1", nil)
			})
			It("should return error", func() {
				err = voteCount.VoteCounting(1)
				Expect(err.Error()).To(Equal(constant.ErrAlgorithmMismatch.Error()))
			})
		})

		When("Voting algorithm match", func() {
			BeforeEach(func() {
				mockContractRepo.EXPECT().
					GetVotedAlgorithm().
					Return(constant.VotingAlgorithm, nil)
				mockISyncService.EXPECT().
					UncountedProposalList(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]model.ProposalTbl{
						{
							ProposalId: 1,
							Percentage: model.Percentage{
								SpPercentage:          2500,
								DeveloperPercentage:   2500,
								TokenHolderPercentage: 2500,
								ClientPercentage:      2500,
							},
						},
					}, nil)
			})
			Context("Ensure the overall process", func() {
				It("Need a vote counting algorithm consistently", func() {
					res := strings.ReplaceAll(constant.VotingAlgorithm, " ", "") == strings.ReplaceAll(`((SpPower / totalPower) * SpPercentage + (DeveloperPower / totalPower) * DeveloperPercentage + (ClientPower / totalPower) * ClientPercentage + (TokenPower       / totalPower)*TokenHolderPercentage)/percentage * 100%`, " ", "")
					Expect(res).To(BeTrue())
				})
				It("should return uncounted proposal list", func() {
					res, err := mockISyncService.UncountedProposalList(context.Background(), 1, 1)
					Expect(err).To(BeNil())
					Expect(len(res)).To(Equal(1))
				})

				It("should return all powers", func() {
					res, err := mockSnapshotClient.GetAllAddrPowerByDay(context.Background(), nil)
					Expect(err).To(BeNil())
					Expect(res).ToNot(Equal(BeNil()))
				})

				It("should return uncounted vote list", func() {
					res, err := mockISyncService.GetUncountedVotedList(context.Background(), 1, 1)
					Expect(err).To(BeNil())
					Expect(len(res)).To(Equal(3))
				})

				Describe("CountWeightCredits", func() {
					var powersMap map[string]model.AddrPower
					BeforeEach(func() {
						grpcRes, err := mockSnapshotClient.GetAllAddrPowerByDay(context.Background(), nil)
						Expect(err).To(BeNil())
						var allPower model.SnapshotAllPower
						err = json.Unmarshal([]byte(grpcRes.Info), &allPower)
						Expect(err).To(BeNil())
						powersMap = utils.PowersInfoToMap(allPower.AddrPower)
					})

					It("should return count weight credits", func() {
						votesInfo, err := mockISyncService.GetUncountedVotedList(context.Background(), 1, 1)
						Expect(err).To(BeNil())
						creditsMap, totalCredits, voteList := voteCount.CountWeightCredits(1, powersMap, votesInfo, 1)

						Expect(creditsMap).NotTo(BeNil())
						Expect(totalCredits).To(Equal(model.VoterPowerCount{
							SpPower:        utils.StringToDecimal("400"),
							DeveloperPower: utils.StringToDecimal("400"),
							TokenPower:     utils.StringToDecimal("400"),
							ClientPower:    utils.StringToDecimal("400"),
						}))
						Expect(len(voteList)).To(Equal(3))
					})
					It("should return the approve percentage", func() {
						res := voteCount.CalculateVotesPercentage(model.VoterPowerCount{
							SpPower:        utils.StringToDecimal("100"),
							DeveloperPower: utils.StringToDecimal("100"),
							TokenPower:     utils.StringToDecimal("200"),
							ClientPower:    utils.StringToDecimal("300"),
						}, model.VoterPowerCount{
							SpPower:        utils.StringToDecimal("400"),
							DeveloperPower: utils.StringToDecimal("400"),
							TokenPower:     utils.StringToDecimal("400"),
							ClientPower:    utils.StringToDecimal("400"),
						}, model.Percentage{
							SpPercentage:          2500,
							DeveloperPercentage:   2500,
							TokenHolderPercentage: 2500,
							ClientPercentage:      2500,
						})
						Expect(res.String()).To(Equal("43.75"))
					})
					It("should return the reject percentage", func() {
						res := voteCount.CalculateVotesPercentage(model.VoterPowerCount{
							SpPower:        utils.StringToDecimal("300"),
							DeveloperPower: utils.StringToDecimal("300"),
							TokenPower:     utils.StringToDecimal("200"),
							ClientPower:    utils.StringToDecimal("100"),
						}, model.VoterPowerCount{
							SpPower:        utils.StringToDecimal("400"),
							DeveloperPower: utils.StringToDecimal("400"),
							TokenPower:     utils.StringToDecimal("400"),
							ClientPower:    utils.StringToDecimal("400"),
						}, model.Percentage{
							SpPercentage:          2500,
							DeveloperPercentage:   2500,
							TokenHolderPercentage: 2500,
							ClientPercentage:      2500,
						})
						Expect(res.String()).To(Equal("56.25"))
					})
					It("should return final percentage", func() {
						approve := decimal.NewFromFloat(43.75)
						reject := decimal.NewFromFloat(56.25)
						res := voteCount.CalculateFinalPercentages(approve, reject, 3)
						Expect(res).To(Equal(map[string]float64{
							constant.VoteApprove: 43.75,
							constant.VoteReject:  56.25,
						}))
					})
				})
			})
			It("should return no error", func() {
				err = voteCount.VoteCounting(1)
				Expect(err).To(BeNil())
			})

		})

	})
})
