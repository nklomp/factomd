package simtest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/state"
	. "github.com/FactomProject/factomd/testHelper"
	"github.com/stretchr/testify/assert"
)

var logName string = "simTest"

// KLUDGE likely already exists elsewhere
func encode(s string) []byte {
	b := bytes.Buffer{}
	b.WriteString(s)
	return b.Bytes()
}

func waitForAnyDeposit(s *state.State, ecPub string) int64 {
	return waitForEcBalance(s, ecPub, 1)
}

func waitForZero(s *state.State, ecPub string) int64 {
	fmt.Println("Waiting for Zero Balance")
	return waitForEcBalance(s, ecPub, 0)
}

func waitForEmptyHolding(s *state.State, msg string) {
	t := time.Now()
	s.LogPrintf(logName, "WaitForEmptyHolding %v %v", t, msg)

	for len(s.Holding) > 0 {
		time.Sleep(time.Millisecond * 50)
	}

	t = time.Now()
	s.LogPrintf(logName, "EmptyHolding %v %v", t, msg)

}

func waitForEcBalance(s *state.State, ecPub string, target int64) int64 {

	s.LogPrintf(logName, "WaitForBalance%v:  %v", target, ecPub)

	for {
		bal := engine.GetBalanceEC(s, ecPub)
		time.Sleep(time.Millisecond * 200)
		//fmt.Printf("WaitForBalance: %v => %v\n", ecPub, bal)

		if (target == 0 && bal == 0) || (target > 0 && bal >= target) {
			s.LogPrintf(logName, "FoundBalance%v: %v", target, bal)
			return bal
		}
	}
}

func TestSendingCommitAndReveal(t *testing.T) {
	if RanSimTest {
		return
	}
	RanSimTest = true

	id := "92475004e70f41b94750f4a77bf7b430551113b25d3d57169eadca5692bb043d"
	extids := [][]byte{encode("foo"), encode("bar")}
	a := AccountFromFctSecret("Fs2zQ3egq2j99j37aYzaCddPq9AF3mgh64uG9gRaDAnrkjRx3eHs")
	b := GetBankAccount()

	t.Run("generate accounts", func(t *testing.T) {
		println(b.String())
		println(a.String())
	})

	t.Run("Run sim to create entries", func(t *testing.T) {
		state0 := SetupSim("L", map[string]string{"--debuglog": ""}, 200, 1, 1, t)

		stop := func() {
			ShutDownEverything(t)
			WaitForAllNodes(state0)
		}

		t.Run("Create Chain", func(t *testing.T) {
			e := factom.Entry{
				ChainID: id,
				ExtIDs:  extids,
				Content: encode("Hello World!"),
			}

			c := factom.NewChain(&e)

			commit, _ := ComposeChainCommit(a.Priv, c)
			reveal, _ := ComposeRevealEntryMsg(a.Priv, c.FirstEntry)

			state0.APIQueue().Enqueue(commit)
			state0.APIQueue().Enqueue(reveal)

			t.Run("Fund ChainCommit Address", func(t *testing.T) {
				amt := uint64(11)
				engine.FundECWallet(state0, b.FctPrivHash(), a.EcAddr(), amt*state0.GetFactoshisPerEC())
				waitForAnyDeposit(state0, a.EcPub())
			})
		})

		t.Run("Generate Entries in Batches", func(t *testing.T) {
			waitForZero(state0, a.EcPub())
			GenerateCommitsAndRevealsInBatches(t, state0)
		})

		t.Run("End simulation", func(t *testing.T) {
			// FIXME somehow the EC entry math is wrong - 10 EC end up being stuck here
			//waitForZero(state0, a.EcPub())
			ht := state0.GetDBHeightComplete()
			WaitBlocks(state0, 2)
			newHt := state0.GetDBHeightComplete()
			assert.True(t, ht < newHt, "block height should progress")
			stop()
		})

	})
}

func GenerateCommitsAndRevealsInBatches(t *testing.T, state0 *state.State) {

	// KLUDGE vars duplicated from original test - should refactor
	id := "92475004e70f41b94750f4a77bf7b430551113b25d3d57169eadca5692bb043d"
	a := AccountFromFctSecret("Fs2zQ3egq2j99j37aYzaCddPq9AF3mgh64uG9gRaDAnrkjRx3eHs")
	b := GetBankAccount()

	// NOTE to send more entries/batches change numbers here
	var numEntries int = 100 // set the total number of entries to add
	for BatchID := 0; BatchID < 10; BatchID++ {

		publish := func(i int) {

			extids := [][]byte{encode(fmt.Sprintf("batch%v", BatchID))}

			e := factom.Entry{
				ChainID: id,
				ExtIDs:  extids,
				Content: encode(fmt.Sprintf("batch %v, seq: %v", BatchID, i)), // ensure no duplicate msg hashes
			}
			i++

			commit, _ := ComposeCommitEntryMsg(a.Priv, e)
			reveal, _ := ComposeRevealEntryMsg(a.Priv, &e)

			state0.APIQueue().Enqueue(commit)
			state0.APIQueue().Enqueue(reveal)
		}

		t.Run(fmt.Sprintf("Create Entries Batch %v", BatchID), func(t *testing.T) {

			waitForEmptyHolding(state0, fmt.Sprintf("START%v", BatchID))

			for x := 1; x <= numEntries; x++ {
				publish(x)
			}

			t.Run("Fund EC Address", func(t *testing.T) {
				amt := uint64(numEntries)
				engine.FundECWallet(state0, b.FctPrivHash(), a.EcAddr(), amt*state0.GetFactoshisPerEC())
				waitForAnyDeposit(state0, a.EcPub())
			})

			waitForEmptyHolding(state0, fmt.Sprintf("END%v", BatchID))

			t.Run("Verify Entries", func(t *testing.T) {

				bal := engine.GetBalanceEC(state0, a.EcPub())
				assert.Equal(t, bal, int64(0))

				// TODO: actually check for confirmed entries
				assert.Equal(t, 0, len(state0.Holding), "messages stuck in holding")
				WaitBlocks(state0, 1)
			})
		})

	}

}
