package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
)

type LobbyMatch struct {
	runtime.Match
}

type LobbyMatchState struct {
	presence   map[string]runtime.Presence
	emptyTicks int
}

func (m *LobbyMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := LobbyMatchState{
		emptyTicks: 0,
		presence:   map[string]runtime.Presence{},
	}

	tickRate := 1
	label := ""
	return state, tickRate, label
}

func (m *LobbyMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	logger.Info("LobbyMatch MatchJoinAttempt")
	return state, true, ""
}

func (m *LobbyMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("LobbyMatch MatchJoin")

	lobbyState, ok := state.(LobbyMatchState)
	if !ok {
		logger.Error("state not valid lobby state object")
	}

	for i := 0; i < len(presences); i++ {
		lobbyState.presence[presences[i].GetSessionId()] = presences[i]
	}
	return lobbyState
}

func (m *LobbyMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("LobbyMatch MatchLeave")

	lobbyState, ok := state.(LobbyMatchState)
	if !ok {
		logger.Error("state not valid lobby state object")
	}

	for i := 0; i < len(presences); i++ {
		delete(lobbyState.presence, presences[i].GetSessionId())
	}
	return lobbyState
}

func (m *LobbyMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	lobbyState, ok := state.(LobbyMatchState)
	if !ok {
		logger.Error("state not valid lobby state object")
	}

	if len(lobbyState.presence) == 0 {
		lobbyState.emptyTicks++
	}

	//if lobbyState.emptyTicks > 100 {
	//	return nil
	//}

	return lobbyState
}

func (m *LobbyMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}

func (m *LobbyMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, ""
}
