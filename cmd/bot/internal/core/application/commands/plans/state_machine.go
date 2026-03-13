package plans

import "strconv"

type Step int

const (
	StepGrid Step = iota
	StepCatalog
	StepDetails
	StepPeriods
	StepMethods
	StepPayment
)

type requestState struct {
	Step          Step
	PlanID        int
	Bandwidth     int
	Period        int
	PaymentMethod string
}

func (uc *UseCase) parseArgs(args []string) requestState {
	state := requestState{
		Step:      StepDetails,
		PlanID:    0,
		Bandwidth: uc.defaultBandwidth,
	}

	if len(args) == 0 {
		return state
	}

	if args[0] == "grid" {
		state.Step = StepGrid
		return state
	}

	state.PlanID, _ = strconv.Atoi(args[0])

	if len(args) >= 2 {
		s, err := strconv.Atoi(args[1])
		if err == nil {
			state.Bandwidth = s
		}
	}

	switch len(args) {
	case 1, 2:
		state.Step = StepDetails
	case 3:
		state.Period, _ = strconv.Atoi(args[2])
		state.Step = StepPeriods
	case 4:
		state.Period, _ = strconv.Atoi(args[2])
		state.PaymentMethod = args[3]
		state.Step = StepMethods
	case 5:
		state.Period, _ = strconv.Atoi(args[2])
		state.PaymentMethod = args[3]
		state.Step = StepPayment
	}

	return state
}
