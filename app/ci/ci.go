package ciApp

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/stellar/app/emerge"
	"go.polydawn.net/stellar/gadgets/ingest/git"
	"go.polydawn.net/stellar/gadgets/layout"
)

func Loop(landmarks layout.Landmarks, mod api.Module, stdout, stderr io.Writer) error {
	var hingeIngest api.ImportRef_Ingest
	for _, imp := range mod.Imports {
		switch imp2 := imp.(type) {
		case api.ImportRef_Ingest:
			switch imp2.IngestKind {
			case "git":
				if hingeIngest != (api.ImportRef_Ingest{}) {
					return fmt.Errorf("a module for use in CI mode can only have one ingest!")
				}
				hingeIngest = imp2
			default:
				return fmt.Errorf("a module for use in CI mode can only have one ingest, and it must be 'ingest:git'")
			}
		}
	}
	if hingeIngest == (api.ImportRef_Ingest{}) {
		return fmt.Errorf("a module for use in CI mode must have one ingest, and it must be 'ingest:git'")
	}
	previouslyIngested := api.WareID{}
	for {
		gitResolve := gitingest.Config{landmarks.ModuleRoot}.Resolve
		newlyIngested, _, err := gitResolve(context.Background(), hingeIngest)
		if err != nil {
			return err
		}
		if *newlyIngested == previouslyIngested {
			time.Sleep(1260 * time.Millisecond)
			continue
		}
		fmt.Fprintf(stderr, "found new git hash!  evaluating %s\n", newlyIngested)
		if err := emergeApp.EvalModule(landmarks, nil, mod, stdout, stderr); err != nil {
			return err
		}
		fmt.Fprintf(stderr, "CI execution done, successfully.  Going into standby until more changes.\n")
		previouslyIngested = *newlyIngested
	}
}
