package tcx

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	tcx, err := ReadFile("testdata/sample.tcx")
	if err != nil {
		t.Error(err)
	}
	npts := len(tcx.Acts.Act[0].Laps[0].Trk.Pt)
	nlaps := len(tcx.Acts.Act[0].Laps)
	nacts := len(tcx.Acts.Act)
	if nlaps != 1 || nacts != 1 {
		t.Error("# Laps parsed:", nlaps)
		t.Error("# Activities parsed:", nacts)
	}
	finalPt := tcx.Acts.Act[0].Laps[0].Trk.Pt[npts-1]

	if finalPt.Lat != -33.8010996 || finalPt.Long != 151.2997607 {
		t.Error("Lat/Long parsed incorrectly.")
		t.Error("Got:", finalPt.Lat, finalPt.Long)
	}
	return
}
