package cmd

import (
	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

type instancePoolItem struct {
	ID    *egoscale.UUID             `json:"id"`
	Name  string                     `json:"name"`
	Zone  string                     `json:"zone"`
	Size  int                        `json:"size"`
	State egoscale.InstancePoolState `json:"state"`
}

type instancePoolListItemOutput []instancePoolItem

func (o *instancePoolListItemOutput) toJSON()  { outputJSON(o) }
func (o *instancePoolListItemOutput) toText()  { outputText(o) }
func (o *instancePoolListItemOutput) toTable() { outputTable(o) }

var instancePoolListCmd = &cobra.Command{
	Use:     "list [zone]",
	Short:   "List instance pool",
	Aliases: gListAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		var zoneName string
		if len(args) == 0 {
			zoneName = gCurrentAccount.DefaultZone
		} else {
			zoneName = args[0]
		}

		zone, err := getZoneByName(zoneName)
		if err != nil {
			return err
		}

		resp, err := cs.RequestWithContext(gContext, egoscale.ListInstancePool{
			ZoneID: zone.ID,
		})
		if err != nil {
			return err
		}
		r := resp.(*egoscale.ListInstancePoolsResponse)
		o := make(instancePoolListItemOutput, 0, r.Count)
		for _, i := range r.ListInstancePoolsResponse {
			o = append(o, instancePoolItem{
				ID:    i.ID,
				Name:  i.Name,
				Zone:  zone.Name,
				Size:  i.Size,
				State: i.State,
			})
		}

		return output(&o, nil)
	},
}

func init() {
	instancePoolCmd.AddCommand(instancePoolListCmd)
}
