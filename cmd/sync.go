package cmd

import (
	"github.com/scribd/keymaster/pkg/keymaster"
	"github.com/spf13/cobra"
	"log"
)

var configPath string

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&configPath, "config", "c", "", "Secret config file or directory containing config yaml's.")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs Secret config yamls with Vault",
	Long: `
Syncs Secret config yamls with Vault.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && configPath == "" {
			log.Fatal("No config files or directories provided.  Cannot continue.  Try again with `dbt keymaster sync -f <file>` or `dbt keymaster sync <file1> <file2> ...`")
		}

		if configPath != "" {
			args = append(args, configPath)
		}

		data, err := keymaster.LoadSecretYamls(args, verbose)
		if err != nil {
			log.Fatalf("failed to load secret definitions: %s", err)
		}

		client, err := auth.Auth()
		if err != nil {
			log.Fatalf("Failed to create Vault client: %s", err)
		}

		km := keymaster.NewKeyMaster(client)

		km.SetTlsAuthCaCert(`-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`)

		clusters := make([]*keymaster.Cluster, 0)
		// bravo
		bravo := keymaster.Cluster{
			Name:         "bravo",
			ApiServerUrl: "https://kube-bravo-master01:6443",
			CACert: `-----BEGIN CERTIFICATE-----
...
Kxq0lynHENJpP/eXjfyC8sLDVJN8YO3n4w==
-----END CERTIFICATE-----`,
			Environment: "production",
			BoundCidrs: []string{
				"1.2.3.4",
				"1.2.3.5",
				"1.2.3.6",
				"1.2.3.7",
				"1.2.3.8",
				"1.2.3.9",
			},
		}

		clusters = append(clusters, &bravo)

		km.SetK8sClusters(clusters)

		for _, config := range data {
			team, err := km.NewTeam(config, verbose)
			if err != nil {
				log.Fatalf("Failed to load secret definitions: %s", err)
			} else {
				err = km.ConfigureTeam(team, verbose)
				if err != nil {
					log.Fatalf("Failed to configure vault for secret definitions: %s", err)
				}
			}
		}
	},
}
