package server

import (
	"SOMAS2023/internal/common/objects"
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/common/voting"
	"encoding/json"
	"fmt"
	"math"
	"os"

	baseserver "github.com/MattSScott/basePlatformSOMAS/BaseServer"
	"github.com/google/uuid"
)

const LootBoxCount = BikerAgentCount * 2.5 // 2.5 lootboxes available per Agent
const MegaBikeCount = 11                   // Megabikes should have 8 riders
const BikerAgentCount = 56                 // 56 agents in total

type IBaseBikerServer interface {
	baseserver.IServer[objects.IBaseBiker]
	GetMegaBikes() map[uuid.UUID]objects.IMegaBike
	GetLootBoxes() map[uuid.UUID]objects.ILootBox
	GetAudi() objects.IAudi
	GetJoiningRequests([]uuid.UUID) map[uuid.UUID][]uuid.UUID
	GetRandomBikeId() uuid.UUID
	RulerElection(agents []objects.IBaseBiker, governance utils.Governance) uuid.UUID
	RunRulerAction(bike objects.IMegaBike) uuid.UUID
	RunDemocraticAction(bike objects.IMegaBike, weights map[uuid.UUID]float64) uuid.UUID
	NewGameStateDump(iteration int) GameStateDump
	GetLeavingDecisions(gameState objects.IGameState) []uuid.UUID
	HandleKickoutProcess() []uuid.UUID
	ProcessJoiningRequests(inLimbo []uuid.UUID)
	RunActionProcess()
	AudiCollisionCheck()
	AddAgentToBike(agent objects.IBaseBiker)
	FoundingInstitutions()
	GetWinningDirection(finalVotes map[uuid.UUID]voting.LootboxVoteMap, weights map[uuid.UUID]float64) uuid.UUID
	LootboxCheckAndDistributions()
	ResetGameState()
	GetDeadAgents() map[uuid.UUID]objects.IBaseBiker
	UpdateGameStates()
}

type Server struct {
	baseserver.BaseServer[objects.IBaseBiker]
	lootBoxes map[uuid.UUID]objects.ILootBox
	megaBikes map[uuid.UUID]objects.IMegaBike
	// megaBikeRiders is a mapping from Agent ID -> ID of the bike that they are riding
	// helps with efficiently managing ridership status
	megaBikeRiders  map[uuid.UUID]uuid.UUID
	audi            objects.IAudi
	deadAgents      map[uuid.UUID]objects.IBaseBiker
	foundingChoices map[uuid.UUID]utils.Governance
}

func Initialize(iterations int) IBaseBikerServer {
	server := &Server{
		BaseServer:     *baseserver.CreateServer[objects.IBaseBiker](GetAgentGenerators(), iterations),
		lootBoxes:      make(map[uuid.UUID]objects.ILootBox),
		megaBikes:      make(map[uuid.UUID]objects.IMegaBike),
		megaBikeRiders: make(map[uuid.UUID]uuid.UUID),
		deadAgents:     make(map[uuid.UUID]objects.IBaseBiker),
		audi:           objects.GetIAudi(),
	}
	server.replenishLootBoxes()
	server.replenishMegaBikes()

	return server
}

func (s *Server) RemoveAgent(agent objects.IBaseBiker) {
	id := agent.GetID()
	// add agent to dead agent map
	s.deadAgents[id] = agent
	// remove agent from agent map
	s.BaseServer.RemoveAgent(agent)
	if bikeId, ok := s.megaBikeRiders[id]; ok {
		s.megaBikes[bikeId].RemoveAgent(id)
		delete(s.megaBikeRiders, id)
	}
}

func (s *Server) AddAgentToBike(agent objects.IBaseBiker) {
	// Remove the agent from the old bike, if it was on one
	if oldBikeId, ok := s.megaBikeRiders[agent.GetID()]; ok {
		s.megaBikes[oldBikeId].RemoveAgent(agent.GetID())
	}

	// set agent on desired bike
	bikeId := agent.GetBike()
	s.megaBikes[bikeId].AddAgent(agent)
	s.megaBikeRiders[agent.GetID()] = bikeId
	if !agent.GetBikeStatus() {
		agent.ToggleOnBike()
	}
}

func (s *Server) RemoveAgentFromBike(agent objects.IBaseBiker) {
	bike := s.megaBikes[agent.GetBike()]
	bike.RemoveAgent(agent.GetID())
	agent.ToggleOnBike()

	// get new destination for agent
	targetBike := agent.ChangeBike()
	if _, ok := s.megaBikes[targetBike]; !ok {
		panic("agent requested a bike that doesn't exist")
	}
	agent.SetBike(targetBike)

	if _, ok := s.megaBikeRiders[agent.GetID()]; ok {
		delete(s.megaBikeRiders, agent.GetID())
	}
}

func (s *Server) GetDeadAgents() map[uuid.UUID]objects.IBaseBiker {
	return s.deadAgents
}

func CalculateAverage(values map[uuid.UUID]float64) float64 {
	var total float64
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}

func CalculateVariance(values map[uuid.UUID]float64, mean float64) float64 {
	var sumOfSquares float64
	for _, value := range values {
		difference := value - mean
		sumOfSquares += difference * difference
	}
	return sumOfSquares / float64(len(values))
}

func filterMapByGroup(originalMap map[uuid.UUID]float64, groupMap map[uuid.UUID]int, groupID int) map[uuid.UUID]float64 {
	filteredMap := make(map[uuid.UUID]float64)
	for key, value := range originalMap {
		if groupMap[key] == groupID {
			filteredMap[key] = value
		}
	}
	return filteredMap
}

type GroupStatistics struct {
	GroupID          int
	AverageLifetime  float64
	VarianceLifetime float64
	AverageEnergy    float64
	VarianceEnergy   float64
	AveragePoints    float64
	VariancePoints   float64
}

// func calculateGroupStatistics(statistics GameStatistics, groupID int) GroupStatistics {
//     groupLifetime := filterMapByGroup(statistics.Average.AgentLifetime, statistics.AgentIDToGroupID, groupID)
//     groupEnergyAverage := filterMapByGroup(statistics.Average.AgentEnergyAverage, statistics.AgentIDToGroupID, groupID)
//     groupPointsAverage := filterMapByGroup(statistics.Average.AgentPointsAverage, statistics.AgentIDToGroupID, groupID)

//     avgLifetime := CalculateAverage(groupLifetime)
//     varLifetime := CalculateVariance(groupLifetime, avgLifetime)

//     avgEnergyAverage := CalculateAverage(groupEnergyAverage)
//     varEnergyAverage := CalculateVariance(groupEnergyAverage, avgEnergyAverage)

//     avgPointsAverage := CalculateAverage(groupPointsAverage)
//     varPointsAverage := CalculateVariance(groupPointsAverage, avgPointsAverage)

//     return GroupStatistics{
//         GroupID:          groupID,
//         AverageLifetime:  avgLifetime,
//         VarianceLifetime: varLifetime,
//         AverageEnergy:    avgEnergyAverage,
//         VarianceEnergy:   varEnergyAverage,
//         AveragePoints:    avgPointsAverage,
//         VariancePoints:   varPointsAverage,
//     }
// }

func (s *Server) outputResults(gameStates [][]GameStateDump) {
	statistics := CalculateStatistics(gameStates)

	// statisticsJson, err := json.MarshalIndent(statistics.Average, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Average Statistics:\n" + string(statisticsJson))
	avgLifetime := CalculateAverage(statistics.Average.AgentLifetime)
	fmt.Printf("Average Agent Lifetime: %v\n", avgLifetime)
	fmt.Printf("Variance in Agent Lifetime: %v\n", math.Sqrt(CalculateVariance(statistics.Average.AgentLifetime, avgLifetime)))

	avgEnergy := CalculateAverage(statistics.Average.AgentEnergyAverage)
	fmt.Printf("Average Agent Energy Average: %v\n", avgEnergy)
	fmt.Printf("Variance in Agent Energy: %v\n", math.Sqrt(CalculateVariance(statistics.Average.AgentEnergyAverage, avgEnergy)))

	avgPointsAverage := CalculateAverage(statistics.Average.AgentPointsAverage)
	fmt.Printf("Average Agent Points Average: %v\n", avgPointsAverage)
	fmt.Printf("Variance in Agent Points: %v\n", math.Sqrt(CalculateVariance(statistics.Average.AgentPointsAverage, avgPointsAverage)))

	// Filter maps for Group 2 and Group 8 for all attributes
	// List of group IDs to analyze
	// groupIDs := []int{1, 2, 3, 4, 5, 6, 7, 8} // Add more group IDs as needed

	// // Slice to hold statistics for all groups
	// var allGroupStats []GroupStatistics

	// // Calculate statistics for each group and add to the slice
	// for _, groupID := range groupIDs {
	//     groupStats := calculateGroupStatistics(statistics, groupID)
	//     allGroupStats = append(allGroupStats, groupStats)

	//     // Print statistics for the group
	//     fmt.Printf("Group %d Statistics:\n", groupStats.GroupID)
	//     fmt.Printf("Average Agent Lifetime: %v\n", groupStats.AverageLifetime)
	//     fmt.Printf("Variance in Agent Lifetime: %v\n", groupStats.VarianceLifetime)
	//     fmt.Printf("Average Agent Energy Average: %v\n", groupStats.AverageEnergy)
	//     fmt.Printf("Variance in Agent Energy Average: %v\n", groupStats.VarianceEnergy)
	//     fmt.Printf("Average Agent Points Average: %v\n", groupStats.AveragePoints)
	//     fmt.Printf("Variance in Agent Points Average: %v\n", groupStats.VariancePoints)
	//     fmt.Println()
	// }

	// redColor := "\033[31m"
	// resetColor := "\033[0m"

	// // Rank by Average Agent Lifetime
	// sort.Slice(allGroupStats, func(i, j int) bool {
	//     return allGroupStats[i].AverageLifetime > allGroupStats[j].AverageLifetime
	// })

	// fmt.Println("Ranking by Average Agent Lifetime:")
	// for i, gs := range allGroupStats {
	// 	if gs.GroupID == 2 || gs.GroupID == 8 {
	//         fmt.Printf("Rank %d: %sGroup %d%s with Average Lifetime %v\n", i+1, redColor, gs.GroupID, resetColor, gs.AverageLifetime)
	//     } else {
	//         fmt.Printf("Rank %d: Group %d with Average Lifetime %v\n", i+1, gs.GroupID, gs.AverageLifetime)
	//     }
	// }

	// // Rank by Average Agent Points
	// sort.Slice(allGroupStats, func(i, j int) bool {
	//     return allGroupStats[i].AveragePoints > allGroupStats[j].AveragePoints
	// })

	// fmt.Println("\nRanking by Average Agent Points:")
	// for i, gs := range allGroupStats {
	//     if gs.GroupID == 2 || gs.GroupID == 8 {
	//         fmt.Printf("Rank %d: %sGroup %d%s with Average Point %v\n", i+1, redColor, gs.GroupID, resetColor, gs.AveragePoints)
	//     } else {
	//         fmt.Printf("Rank %d: Group %d with Average Point %v\n", i+1, gs.GroupID, gs.AveragePoints)
	//     }
	// }

	file, err := os.Create("statistics.xlsx")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := statistics.ToSpreadsheet().Write(file); err != nil {
		panic(err)
	}

	file, err = os.Create("game_dump.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(gameStates); err != nil {
		panic(err)
	}
}

func (s *Server) UpdateGameStates() {
	gs := s.NewGameStateDump(0)
	for _, agent := range s.GetAgentMap() {
		agent.UpdateGameState(gs)
	}
}

// had to override to address the fact that agents only have access to the game dump
// version of agents, so if the recipients are set to be those it will panic as they
// can't call the handler functions
func (s *Server) RunMessagingSession() {
	agentArray := s.GenerateAgentArrayFromMap()

	for _, agent := range s.GetAgentMap() {
		allMessages := agent.GetAllMessages(agentArray)
		for _, msg := range allMessages {
			recipients := msg.GetRecipients()
			// make recipient list with actual agents
			usableRecipients := make([]objects.IBaseBiker, len(recipients))
			for i, recipient := range recipients {
				usableRecipients[i] = s.GetAgentMap()[recipient.GetID()]
			}
			for _, recip := range usableRecipients {
				if agent.GetID() == recip.GetID() {
					continue
				}
				msg.InvokeMessageHandler(recip)
			}
		}
	}
}
