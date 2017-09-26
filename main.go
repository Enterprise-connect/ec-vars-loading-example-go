package main

import (
	"flag"
	"fmt"
	"os"

	api "github.build.ge.com/212359746/ecapi"
	config "github.build.ge.com/212359746/ecconf"
	core "github.build.ge.com/212359746/eccore"
	util "github.build.ge.com/212359746/ecutil"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	mapping = map[string]string{
		"_md": "mod",
		"_v":  "ver",
		"_lp": "lpt",
		"_ts": "tid",
		"_gh": "hst",
		"_sh": "sst",
		"_sg": "sgn",
		"_rh": "rht",
		"_rp": "rpt",
		"_id": "aid",
		"_tk": "tkn",
		"_px": "pxy",
		"_ci": "cid",
		"_cs": "csc",
		"_oa": "oa2",
		"_du": "dur",
		"_ct": "crt",
		"_wl": "wtl",
		"_bl": "bkl",
		"_pg": "plg",
		"_if": "inf",
		"_dg": "dbg",
		"_vf": "vfy",
		"_zn": "zon",
		"_pk": "pky",
		"_pc": "pct",
		"_hc": "hca",
		"_gc": "gen",
	}

	cfm = make(map[string]interface{})
)

func init() {
	flag.String("mod", "", `Specify the EC Agent Mode in "client", "server", or "gateway".`)
	flag.Bool("ver", false, "Show EC Agent's version.")
	flag.String("lpt", "", `Specify the EC port# if the "client" mode is set.`)
	flag.String("tid", "", `Specify the Target EC Server Id if the "client" mode is set`)
	flag.String("hst", "", "Specify the EC Gateway URI. E.g. wss://<somedomain>:8989")
	flag.String("sst", "", "Specify the EC Service URI. E.g. https://<service.of.predix.io>")
	flag.Bool("sgn", false, "Start a CA Cert-Signing process.")
	flag.String("rht", "", `Specify the Resource Host if the "server" mode is set. E.g. <someip>, <somedomain>. value will be discard when TLS is specified.`)
	flag.String("rpt", "0", `Specify the Resource Port# if the "server" mode is set. E.g. 8989, 38989`)
	flag.String("aid", "", "Specify the agent Id assigned by the EC Service. You may find it in the Cloud Foundry VCAP_SERVICE")
	flag.String("tkn", "", "Specify the OAuth Token. The token may expire depending on your OAuth provisioner. This flag is ignored if OAuth2 Auto-Refresh were set.")
	flag.String("pxy", "", "Specify a local Proxy service. E.g. http://hello.world.com:8080")
	flag.String("cid", "", "Specify the client Id to auto-refresh the OAuth2 token.")
	flag.String("csc", "", "Specify the client secret to auto-refresh the OAuth2 token.")
	flag.String("oa2", "", "Specify URL of the OAuth2 provisioner. E.g. https://<somedomain>/oauth/token")
	flag.String("oa22", "", "Specify URL of the OAuth2 provisioner. E.g. https://<somedomain>/oauth/token")
	flag.Int("dur", 0, "Specify the duration for the next token refresh in seconds. (default 100 years)")
	flag.String("crt", "", "Specify the relative path of a digital certificate to operate the EC agent. (.pfx, .cer, .p7s, .der, .pem, .crt)")
	flag.String("wtl", "0.0.0.0/0,::/0", "Specify the ip(s) whitelist in the cidr net format. Concatenate ips by comma. E.g. 89.24.9.0/24, 7.6.0.0/16")
	flag.String("bkl", "", "Specify the ip(s) blocklist in the IPv4/IPv6 format. Concatenate ips by comma. E.g. 10.20.30.5, 2002:4559:1FE2::4559:1FE2")
	flag.String("plg", "", `Enable plugin list. Available options "tls", "ip-route", etc. `)
	flag.Bool("inf", false, "The Product Information.")
	flag.Bool("dbg", false, "Turn on debug mode. This will introduce more error information. E.g. connection error.")
	flag.Bool("vfy", false, "Verify the legitimacy of a digital certificate.")
	flag.String("zon", "", `Specify the Zone/Service Inst. Id. required in the "gateway" mode.`)
	flag.String("pky", "", "Specify the relative path to a TLS key when operate as the gateway as desired. E.g. ./path/to/key.pem.")
	flag.String("pct", "", "Specify the relative path to a TLS cert when operate as the gateway as desired. E.g. ./path/to/cert.pem.")
	flag.String("hca", "", `Specify a port# to turn on the Healthcheck API. This flag is always on when in the "gateway mode" with the provisioned local port. Upon provisioned, the api is available at <agent_uri>/health.`)
	flag.Bool("gen", false, "Generate a certificate request for the usage validation purpose.")
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic("Fatal error config file: " + err.Error())
		}
	}

	// override flag value with env var name uppercased prefixed by "EC_"
	viper.SetEnvPrefix("ec")
	viper.AutomaticEnv()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// save the values in the config map
	for k, v := range mapping {
		if viper.IsSet(v) {
			cfm[k] = viper.Get(v)
		}
	}
}

func main() {
	// dump the values to demo a first pass
	for k, v := range cfm {
		if v == nil {
			fmt.Printf("K: %s V: %v\n", k, v)
		} else {
			switch t := v.(type) {
			case string:
				fmt.Printf("K: %s V: %s\n", k, t)
			case *string:
				fmt.Printf("K: %s V: %s\n", k, *t)
			case *bool:
				fmt.Printf("K: %s V: %t\n", k, *t)
			case *int64:
				fmt.Printf("K: %s V: %v\n", k, *t)
			case *int:
				fmt.Printf("K: %s V: %v\n", k, *t)
			default:
				fmt.Printf("K: %s V: %v\n", k, t)
			}
		}
	}
	// first pass -- demonstrate flags, env, and config options
	return

	defer func() {
		if r := recover(); r != nil {
			util.ErrLog(r)
			os.Exit(1)
		}
	}()

	util.Init(*cfm["_md"].(*string), *cfm["_dg"].(*bool))

	conf, err := config.Init(cfm)
	if err != nil {
		panic(err)
	}

	switch conf.ClientType {
	case "server":

		if err := conf.InitServerConfig(cfm); err != nil {
			panic(err)
		}

	case "client":

		if err := conf.InitClientConfig(cfm); err != nil {
			panic(err)
		}

	case "gateway":
		if err := conf.InitGatewayConfig(cfm); err != nil {
			panic(err)
		}

	default:
		panic("Unknown agent mode.")
	}

	var _agt core.AgentIntr

	switch conf.ClientType {
	case "server":
		_agt = core.NewServer(conf, api.HealthCheckSrv)
	case "client":
		_agt = core.NewClient(conf, api.HealthCheck)
	case "gateway":
		_agt = core.NewGateway(conf, api.HealthCheckGwy)
	default:
		panic("no agent mode specified.")
	}

	go func() {
		_agt.Hire()
		defer func() {
			if r := recover(); r != nil {
				util.ErrLog(r)
				os.Exit(1)
			}
		}()
	}()

	core.Operation(_agt)
}
