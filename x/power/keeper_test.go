package power

import (
	"testing"

	"github.com/quokki/quokki/types"
	"github.com/stretchr/testify/assert"
)

func TestPowerRestore(t *testing.T) {
	ctx, pk, am, _ := createTestInput(t, false)

	_privPower, _ := am.GetAccount(ctx, addrs[0]).Get("QuokkiPower")
	privPower, ok := _privPower.(types.QuokkiPower)
	assert.Equal(t, true, ok)

	ctx = ctx.WithBlockHeight(1)
	pk.PowerUse(ctx, addrs[0], 100, 100)
	for i := 2; i <= 100; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		pk.Tick(ctx)
	}
	_power, _ := am.GetAccount(ctx, addrs[0]).Get("QuokkiPower")
	power, ok := _power.(types.QuokkiPower)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(100), power.Used)
	assert.Equal(t, privPower.Available+privPower.Reserved+privPower.Used, power.Available+power.Reserved+power.Used)

	ctx = ctx.WithBlockHeight(101)
	pk.Tick(ctx)

	_power, _ = am.GetAccount(ctx, addrs[0]).Get("QuokkiPower")
	power, ok = _power.(types.QuokkiPower)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(0), power.Used)
	assert.Equal(t, privPower.Available+privPower.Reserved+privPower.Used, power.Available+power.Reserved+power.Used)
}
