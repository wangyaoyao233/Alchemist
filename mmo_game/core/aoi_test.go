package core

import (
	"fmt"
	"testing"
)

func TestNewAoiManager(t *testing.T) {
	aoiMgr := NewAOIManager(
		0, 250, 5, 0, 250, 5,
	)

	fmt.Println(aoiMgr)
}

func TestAOIManagerSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(
		0, 250, 5, 0, 250, 5,
	)

	for gid := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid:", gid, "grids len:", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Println("surround grid :", gIDs)
	}
}
