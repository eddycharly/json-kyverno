package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func Command(parents ...string) *cobra.Command {
	var command options
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "serve",
		Long:         "Serve API",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE:         command.Run,
	}
	// server flags
	cmd.Flags().StringVar(&command.serverFlags.host, "server-host", "0.0.0.0", "server host")
	cmd.Flags().IntVar(&command.serverFlags.port, "server-port", 8080, "server port")
	// gin flags
	cmd.Flags().StringVar(&command.ginFlags.mode, "gin-mode", gin.ReleaseMode, "gin run mode")
	cmd.Flags().BoolVar(&command.ginFlags.log, "gin-log", true, "enable gin logger")
	cmd.Flags().BoolVar(&command.ginFlags.cors, "gin-cors", true, "enable gin cors")
	cmd.Flags().IntVar(&command.ginFlags.maxBodySize, "gin-max-body-size", 2*1024*1024, "gin max body size")
	// cluster flags
	clientcmd.BindOverrideFlags(&command.clusterFlags.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return cmd
}
