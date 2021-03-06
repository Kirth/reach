package module

import (
	"os"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/reach/gadgets/layout"
)

func Load(landmarks layout.Module) (mod *api.Module, err error) {
	f, err := os.Open(landmarks.ModuleFile())
	if err != nil {
		return
	}
	err = refmt.NewUnmarshallerAtlased(json.DecodeOptions{}, f, api.Atlas_Module).Unmarshal(&mod)
	return
}
