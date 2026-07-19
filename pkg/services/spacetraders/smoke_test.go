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
)

func smokeClient(t *testing.T) *spacetraders.Client {
	t.Helper()

	if os.Getenv("ST_SMOKE") != "1" {
		t.Skip("skipping smoke test; set ST_SMOKE=1 and SPACETRADERS_TOKEN to run")
	}

	token := os.Getenv("SPACETRADERS_TOKEN")
	if token == "" {
		t.Fatal("SPACETRADERS_TOKEN must be set when ST_SMOKE=1")
	}

	client, err := spacetraders.NewClient(spacetraders.Config{
		Token:         token,
		VerboseErrors: true,
	}, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err != nil {
		t.Fatalf("NewClient: %+v", err)
	}

	return client
}

func smokeCtx(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	return context.WithTimeout(context.Background(), 60*time.Second)
}

func TestSmokeAgent(t *testing.T) {
	client := smokeClient(t)

	ctx, cancel := smokeCtx(t)
	defer cancel()

	agent, err := client.GetMyAgent(ctx)
	if err != nil {
		t.Fatalf("GetMyAgent: %+v", err)
	}

	t.Logf("agent: symbol=%s credits=%d hq=%s faction=%s ships=%d",
		agent.Symbol, agent.Credits, agent.Headquarters, agent.StartingFaction, agent.ShipCount)

	if agent.Symbol == "" {
		t.Fatal("agent symbol is empty")
	}

	if agent.Credits <= 0 {
		t.Fatalf("agent credits not positive: %d", agent.Credits)
	}
}

func firstShip(t *testing.T, client *spacetraders.Client, ctx context.Context) ss.Ship {
	t.Helper()

	ships, _, err := client.ListMyShips(ctx, spacetraders.ListShipsRequest{Limit: 20})
	if err != nil {
		t.Fatalf("ListMyShips: %+v", err)
	}

	if len(ships) == 0 {
		t.Fatal("no ships returned")
	}

	for _, s := range ships {
		t.Logf("  ship: symbol=%s role=%s nav=%s waypoint=%s fuel=%d/%d",
			s.Symbol, s.Registration.Role, s.Nav.Status,
			s.Nav.WaypointSymbol, s.Fuel.Current, s.Fuel.Capacity)
	}

	return ships[0]
}

func TestSmokeShips(t *testing.T) {
	client := smokeClient(t)

	ctx, cancel := smokeCtx(t)
	defer cancel()

	first := firstShip(t, client, ctx)

	single, err := client.GetMyShip(ctx, first.Symbol)
	if err != nil {
		t.Fatalf("GetMyShip: %+v", err)
	}

	t.Logf("GetMyShip: symbol=%s frame=%s cargo=%d/%d",
		single.Symbol, single.Frame.Symbol, single.Cargo.Units, single.Cargo.Capacity)

	cooldown, err := client.GetShipCooldown(ctx, first.Symbol)
	if err != nil {
		t.Fatalf("GetShipCooldown: %+v", err)
	}

	t.Logf("cooldown: remaining=%ds expiration=%s",
		cooldown.RemainingSeconds, cooldown.Expiration)
}

func TestSmokeContracts(t *testing.T) {
	client := smokeClient(t)

	ctx, cancel := smokeCtx(t)
	defer cancel()

	contracts, _, err := client.ListMyContracts(ctx, spacetraders.ListContractsRequest{Limit: 20})
	if err != nil {
		t.Fatalf("ListMyContracts: %+v", err)
	}

	t.Logf("contracts: count=%d", len(contracts))

	for _, c := range contracts {
		t.Logf("  contract: id=%s type=%s faction=%s accepted=%v fulfilled=%v",
			c.ID, c.Type, c.FactionSymbol, c.Accepted, c.Fulfilled)
	}
}

func TestSmokeWaypoints(t *testing.T) {
	client := smokeClient(t)

	ctx, cancel := smokeCtx(t)
	defer cancel()

	first := firstShip(t, client, ctx)

	waypoints, _, err := client.ListSystemWaypoints(ctx, spacetraders.ListWaypointsRequest{
		SystemSymbol: first.Nav.SystemSymbol,
		Limit:        50,
	})
	if err != nil {
		t.Fatalf("ListSystemWaypoints: %+v", err)
	}

	t.Logf("waypoints in %s: count=%d", first.Nav.SystemSymbol, len(waypoints))

	foundField := false

	for _, w := range waypoints {
		t.Logf("  waypoint: symbol=%s type=%s", w.Symbol, w.Type)

		if w.Type == ss.WaypointTypeAsteroidField {
			foundField = true
		}
	}

	if !foundField {
		t.Logf("warning: no ASTEROID_FIELD found in %s (may be deeper in pagination)",
			first.Nav.SystemSymbol)
	}
}

func TestSmokePaginationSeq(t *testing.T) {
	client := smokeClient(t)

	ctx, cancel := smokeCtx(t)
	defer cancel()

	ships, _, err := client.ListMyShips(ctx, spacetraders.ListShipsRequest{Limit: 20})
	if err != nil {
		t.Fatalf("ListMyShips: %+v", err)
	}

	count := 0

	for _, seqErr := range client.ListMyShipsSeq(ctx, spacetraders.ListShipsRequest{Limit: 20}) {
		if seqErr != nil {
			t.Fatalf("ListMyShipsSeq error: %+v", seqErr)
		}

		count++
		if count >= len(ships) {
			break
		}
	}

	t.Logf("ListMyShipsSeq yielded %d ships", count)
}
