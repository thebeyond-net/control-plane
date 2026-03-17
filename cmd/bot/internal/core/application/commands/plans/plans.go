package plans

import (
	"context"
	"strconv"
	"time"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	sharedPorts "github.com/thebeyond-net/control-plane/internal/core/ports"
)

type UseCase struct {
	bot              ports.Bot
	bandwidths       []int
	periods          []domain.Period
	plans            []domain.Plan
	paymentMethods   domain.Items
	currencies       domain.Items
	yookassa         sharedPorts.Invoice
	telegramStars    sharedPorts.Invoice
	featureFlags     sharedPorts.FeatureFlags
	defaultBandwidth int
}

func NewUseCase(
	bot ports.Bot,
	bandwidths []int,
	periods []domain.Period,
	plans []domain.Plan,
	paymentMethods domain.Items,
	currencies domain.Items,
	yookassa sharedPorts.Invoice,
	telegramStars sharedPorts.Invoice,
	featureFlags sharedPorts.FeatureFlags,
	defaultBandwidth int,
) ports.CommandHandler {
	return &UseCase{
		bot,
		bandwidths,
		periods,
		plans,
		paymentMethods,
		currencies,
		yookassa,
		telegramStars,
		featureFlags,
		defaultBandwidth,
	}
}

func (uc *UseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	state := uc.parseArgs(msg.Args)
	if state.PlanID >= len(uc.plans) {
		state.PlanID = 0
	}

	plan := uc.plans[state.PlanID]

	switch state.Step {
	case StepGrid:
		return uc.renderCatalogGrid(ctx, msg, user)
	case StepDetails:
		return uc.renderPlanDetails(ctx, msg, user, plan, state.Bandwidth)
	case StepPeriods:
		return uc.renderPeriodSelection(ctx, msg, user, plan, state.Bandwidth)
	case StepCalendar:
		year, _ := strconv.Atoi(msg.Args[3])
		month, _ := strconv.Atoi(msg.Args[4])
		return uc.renderCalendar(ctx, msg, user, plan, state.Bandwidth, year, time.Month(month))
	case StepYearSelection:
		month, _ := strconv.Atoi(msg.Args[3])
		if month == 0 {
			month = int(time.Now().Month())
		}
		return uc.renderYearSelection(ctx, msg, user, plan, state.Bandwidth, time.Month(month))
	case StepMethods:
		return uc.renderPaymentMethods(ctx, msg, user, plan, state)
	case StepPayment:
		return uc.initiatePayment(ctx, msg, user, plan, state)
	default:
		return uc.renderPlanDetails(ctx, msg, user, plan, state.Bandwidth)
	}
}
