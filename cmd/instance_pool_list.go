package cmd

import (
	"fmt"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

type instancePoolItem struct {
	ID    string                     `json:"id"`
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
	Use:   "list",
	Short: "List instance pools",
	Long: fmt.Sprintf(`This command list instance pools.

Supported output template annotations: %s`,
		strings.Join(outputterTemplateAnnotations(&instancePoolItem{}), ", ")),
	Aliases: gListAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneFlag, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}
		zoneFlag = strings.ToLower(zoneFlag)

		var zoneName string
		if zoneFlag == "" {
			zoneName = gCurrentAccount.DefaultZone
		} else {
			zoneName = zoneFlag
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
			z, err := getZoneByName(i.ZoneID.String())
			if err != nil {
				return err
			}

			if z.Name != zoneFlag && zoneFlag != "" {
				continue
			}

			o = append(o, instancePoolItem{
				ID:    i.ID.String(),
				Name:  i.Name,
				Zone:  z.Name,
				Size:  i.Size,
				State: i.State,
			})
		}

		return output(&o, nil)
	},
}

func init() {
	instancePoolListCmd.Flags().StringP("zone", "z", "", "List Instance pool by zone")
	instancePoolCmd.AddCommand(instancePoolListCmd)
}