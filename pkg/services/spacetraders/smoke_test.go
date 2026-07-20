// Package spacetraders_test contains a read-only smoke test for the
// SpaceTraders client. It is skipped unless both SPACETRADERS_TOKEN is set
// and ST_SMOKE=1 is set, so it never runs during normal `go test` and never
// performs mutating API calls.
package spacetraders_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/opoccomaxao/spacetraders/pkg/services/spacetraders"
	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	"github.com/stretchr/testify/suite"
)

type SmokeSuite struct {
	suite.Suite

	client *spacetraders.Client
}

func (s *SmokeSuite) SetupSuite() {
	if os.Getenv("ST_SMOKE") != "1" {
		s.T().Skip("skipping smoke test; set ST_SMOKE=1 and SPACETRADERS_TOKEN to run")
	}

	token := os.Getenv("SPACETRADERS_TOKEN")
	if token == "" {
		s.FailNow("SPACETRADERS_TOKEN must be set when ST_SMOKE=1")
	}

	client, err := spacetraders.NewClient(spacetraders.Config{
		Token:         token,
		VerboseErrors: true,
	}, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	s.Require().NoError(err, "NewClient")

	s.client = client
}

func (s *SmokeSuite) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 60*time.Second)
}

func (s *SmokeSuite) firstShip(ctx context.Context) ss.Ship {
	ships, _, err := s.client.ListMyShips(ctx, spacetraders.ListShipsRequest{Limit: 20})
	s.Require().NoError(err, "ListMyShips")
	s.Require().NotEmpty(ships, "no ships returned")

	for _, ship := range ships {
		s.T().Logf("  ship: symbol=%s role=%s nav=%s waypoint=%s fuel=%d/%d",
			ship.Symbol, ship.Registration.Role, ship.Nav.Status,
			ship.Nav.WaypointSymbol, ship.Fuel.Current, ship.Fuel.Capacity)
	}

	return ships[0]
}

func (s *SmokeSuite) TestAgent() {
	ctx, cancel := s.ctx()
	defer cancel()

	agent, err := s.client.GetMyAgent(ctx)
	s.Require().NoError(err, "GetMyAgent")

	s.T().Logf("agent: symbol=%s credits=%d hq=%s faction=%s ships=%d",
		agent.Symbol, agent.Credits, agent.Headquarters, agent.StartingFaction, agent.ShipCount)

	s.Require().NotEmpty(agent.Symbol, "agent symbol is empty")
	s.Require().Positive(agent.Credits, "agent credits not positive")
}

func (s *SmokeSuite) TestShips() {
	ctx, cancel := s.ctx()
	defer cancel()

	first := s.firstShip(ctx)

	single, err := s.client.GetMyShip(ctx, first.Symbol)
	s.Require().NoError(err, "GetMyShip")

	s.T().Logf("GetMyShip: symbol=%s frame=%s cargo=%d/%d",
		single.Symbol, single.Frame.Symbol, single.Cargo.Units, single.Cargo.Capacity)

	cooldown, err := s.client.GetShipCooldown(ctx, first.Symbol)
	s.Require().NoError(err, "GetShipCooldown")

	s.T().Logf("cooldown: remaining=%ds expiration=%s",
		cooldown.RemainingSeconds, cooldown.Expiration)
}

func (s *SmokeSuite) TestContracts() {
	ctx, cancel := s.ctx()
	defer cancel()

	contracts, _, err := s.client.ListMyContracts(ctx, spacetraders.ListContractsRequest{Limit: 20})
	s.Require().NoError(err, "ListMyContracts")

	s.T().Logf("contracts: count=%d", len(contracts))

	for _, c := range contracts {
		s.T().Logf("  contract: id=%s type=%s faction=%s accepted=%v fulfilled=%v",
			c.ID, c.Type, c.FactionSymbol, c.Accepted, c.Fulfilled)
	}
}

func (s *SmokeSuite) TestWaypoints() {
	ctx, cancel := s.ctx()
	defer cancel()

	first := s.firstShip(ctx)

	waypoints, _, err := s.client.ListSystemWaypoints(ctx, spacetraders.ListWaypointsRequest{
		SystemSymbol: first.Nav.SystemSymbol,
		Limit:        50,
	})
	s.Require().NoError(err, "ListSystemWaypoints")

	s.T().Logf("waypoints in %s: count=%d", first.Nav.SystemSymbol, len(waypoints))

	foundField := false

	for _, w := range waypoints {
		s.T().Logf("  waypoint: symbol=%s type=%s", w.Symbol, w.Type)

		if w.Type == ss.WaypointTypeAsteroidField {
			foundField = true
		}
	}

	if !foundField {
		s.T().Logf("warning: no ASTEROID_FIELD found in %s (may be deeper in pagination)",
			first.Nav.SystemSymbol)
	}
}

func (s *SmokeSuite) TestPaginationSeq() {
	ctx, cancel := s.ctx()
	defer cancel()

	ships, _, err := s.client.ListMyShips(ctx, spacetraders.ListShipsRequest{Limit: 20})
	s.Require().NoError(err, "ListMyShips")

	count := 0

	for _, seqErr := range s.client.ListMyShipsSeq(ctx, spacetraders.ListShipsRequest{Limit: 20}) {
		s.Require().NoError(seqErr, "ListMyShipsSeq error")

		count++
		if count >= len(ships) {
			break
		}
	}

	s.T().Logf("ListMyShipsSeq yielded %d ships", count)
}

func TestSmokeSuite(t *testing.T) {
	suite.Run(t, new(SmokeSuite))
}
