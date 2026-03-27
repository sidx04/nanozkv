package zk

import (
	"fmt"
	"nanozkv/vm"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func RunTraceProof(code []vm.Instruction) error {
	vm := vm.NewVM(code)
	if err := vm.Run(); err != nil {
		return err
	}

	circuitInput := TraceToCircuit(*vm.Trace)

	var circuit TraceCircuit

	cs, err := frontend.Compile(
		ecc.BN254.ScalarField(),
		r1cs.NewBuilder,
		&circuit,
	)
	if err != nil {
		return err
	}

	w, err := frontend.NewWitness(&circuitInput, ecc.BN254.ScalarField())
	if err != nil {
		return err
	}

	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		return err
	}

	proof, err := groth16.Prove(cs, pk, w)
	if err != nil {
		return err
	}

	publicWitness, _ := w.Public()
	if err = groth16.Verify(proof, vk, publicWitness); err != nil {
		return err
	}

	fmt.Println("✅ Trace proof verified!")
	return nil
}
