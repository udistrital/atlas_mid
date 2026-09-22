package services

import (
	"strconv"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type RiskDecision struct {
	Score             int
	RequireTurnstile  bool
	Block             bool
	Reason            string
	RequestsInWindow  int
	RequestsInBurst   int
	UniquePathsWindow int
}

type clientRiskState struct {
	RequestTimes []time.Time
	PathTimes    map[string]time.Time
	LastSeen     time.Time
}

var riskStore = struct {
	sync.Mutex
	Clients map[string]*clientRiskState
	Counter int
}{
	Clients: make(
		map[string]*clientRiskState,
	),
}

var turnstileAttemptStore = struct {
	sync.Mutex
	Attempts map[string][]time.Time
}{
	Attempts: make(
		map[string][]time.Time,
	),
}

func EvaluateRisk(
	clientKey string,
	path string,
) RiskDecision {

	now := time.Now()

	window :=
		time.Duration(
			configInt(
				"RiskWindowSeconds",
				60,
			),
		) * time.Second

	burstWindow :=
		time.Duration(
			configInt(
				"RiskBurstWindowSeconds",
				5,
			),
		) * time.Second

	maxRequests :=
		configInt(
			"RiskMaxRequestsPerWindow",
			60,
		)

	maxBurst :=
		configInt(
			"RiskMaxBurstRequests",
			12,
		)

	maxUniquePaths :=
		configInt(
			"RiskMaxUniquePathsPerWindow",
			30,
		)

	challengeScore :=
		configInt(
			"RiskChallengeScore",
			35,
		)

	blockScore :=
		configInt(
			"RiskBlockScore",
			80,
		)

	riskStore.Lock()
	defer riskStore.Unlock()

	state, exists :=
		riskStore.Clients[clientKey]

	if !exists {
		state = &clientRiskState{
			PathTimes: make(
				map[string]time.Time,
			),
		}

		riskStore.Clients[clientKey] =
			state
	}

	windowStart :=
		now.Add(-window)

	burstStart :=
		now.Add(-burstWindow)

	filtered :=
		state.RequestTimes[:0]

	for _, requestTime := range state.RequestTimes {

		if requestTime.After(
			windowStart,
		) {
			filtered = append(
				filtered,
				requestTime,
			)
		}
	}

	state.RequestTimes =
		filtered

	state.RequestTimes =
		append(
			state.RequestTimes,
			now,
		)

	for registeredPath, lastSeen := range state.PathTimes {

		if lastSeen.Before(
			windowStart,
		) {
			delete(
				state.PathTimes,
				registeredPath,
			)
		}
	}

	state.PathTimes[path] = now
	state.LastSeen = now

	burstCount := 0

	for _, requestTime := range state.RequestTimes {

		if requestTime.After(
			burstStart,
		) {
			burstCount++
		}
	}

	requestCount :=
		len(
			state.RequestTimes,
		)

	uniquePaths :=
		len(
			state.PathTimes,
		)

	score := 0
	reason := ""

	if requestCount >
		maxRequests {

		score += 40

		reason =
			"demasiadas solicitudes en la ventana de tiempo"
	}

	if burstCount >
		maxBurst {

		score += 35

		if reason == "" {
			reason =
				"ráfaga de solicitudes inusualmente alta"
		}
	}

	if uniquePaths >
		maxUniquePaths {

		score += 25

		if reason == "" {
			reason =
				"consulta de demasiados recursos diferentes"
		}
	}

	/*
		Un volumen extremadamente alto
		se bloquea directamente y no se
		envía a Turnstile.
	*/
	if requestCount >
		maxRequests*2 ||
		burstCount >
			maxBurst*2 {

		score = 100

		reason =
			"volumen de solicitudes excesivo"
	}

	if score > 100 {
		score = 100
	}

	riskStore.Counter++

	if riskStore.Counter%500 == 0 {
		cleanupRiskStore(
			now,
			window*2,
		)
	}

	return RiskDecision{
		Score: score,

		RequireTurnstile: score >= challengeScore &&
			score < blockScore,

		Block: score >= blockScore,

		Reason: reason,

		RequestsInWindow: requestCount,

		RequestsInBurst: burstCount,

		UniquePathsWindow: uniquePaths,
	}
}

func ResetRisk(
	clientKey string,
) {

	riskStore.Lock()
	defer riskStore.Unlock()

	delete(
		riskStore.Clients,
		clientKey,
	)
}

func AllowTurnstileVerification(
	clientKey string,
) bool {

	now := time.Now()

	window :=
		time.Minute

	maxAttempts :=
		configInt(
			"TurnstileVerifyMaxAttemptsPerMinute",
			10,
		)

	turnstileAttemptStore.Lock()
	defer turnstileAttemptStore.Unlock()

	attempts :=
		turnstileAttemptStore.
			Attempts[clientKey]

	filtered :=
		attempts[:0]

	windowStart :=
		now.Add(-window)

	for _, attempt := range attempts {

		if attempt.After(
			windowStart,
		) {
			filtered = append(
				filtered,
				attempt,
			)
		}
	}

	if len(filtered) >=
		maxAttempts {

		turnstileAttemptStore.
			Attempts[clientKey] =
			filtered

		return false
	}

	filtered =
		append(
			filtered,
			now,
		)

	turnstileAttemptStore.
		Attempts[clientKey] =
		filtered

	return true
}

func cleanupRiskStore(
	now time.Time,
	maxAge time.Duration,
) {

	for clientKey, state := range riskStore.Clients {

		if now.Sub(
			state.LastSeen,
		) > maxAge {

			delete(
				riskStore.Clients,
				clientKey,
			)
		}
	}
}

func configInt(
	key string,
	defaultValue int,
) int {

	value :=
		web.AppConfig.
			DefaultString(
				key,
				"",
			)

	if value == "" {
		return defaultValue
	}

	parsed, err :=
		strconv.Atoi(value)

	if err != nil ||
		parsed <= 0 {

		return defaultValue
	}

	return parsed
}

func RiskRechallengeScore() int {
	return configInt(
		"RiskRechallengeScore",
		65,
	)
}
