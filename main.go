package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Hello World!")

	err := initializer.RegisterRpc("healthcheck", RpcHealthCheck)
	if err != nil {
		return err
	}

	err = initializer.RegisterMatch("matchdemo", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &LobbyMatch{}, nil
	})

	if err != nil {
		return err
	}

	if err := initializer.RegisterRpc("create_match_rpc", CreateMatchRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	// Register as matchmaker matched hook, this call should be in InitModule.
	if err := initializer.RegisterMatchmakerMatched(MakeMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}

func CreateMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	//params := make(map[string]interface{})

	//if err := json.Unmarshal([]byte(payload), &params); err != nil {
	//	return "", err
	//}
	logger.Info("payload:%s", payload)
	params := map[string]interface{}{
		"some": "data",
	}

	modulename := "matchdemo" // Name with which match handler was registered in InitModule, see example above.

	if matchId, err := nk.MatchCreate(ctx, modulename, params); err != nil {
		return "", err
	} else {
		params["matchId"] = matchId
		out, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		return string(out), nil
	}
}

func MakeMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	for _, e := range entries {
		logger.Info("Matched user '%s' named '%s'", e.GetPresence().GetUserId(), e.GetPresence().GetUsername())

		for k, v := range e.GetProperties() {
			logger.Info("Matched on '%s' value '%v'", k, v)
		}
	}

	matchId, err := nk.MatchCreate(ctx, "matchdemo", map[string]interface{}{"invited": entries})

	if err != nil {
		return "", err
	}

	return matchId, nil
}
