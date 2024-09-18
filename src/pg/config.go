package pg

import (
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/kelseyhightower/envconfig"
)

type ConfigSpec struct {
	NhsTrudKey       string `envconfig:"NHS_TRUD_KEY" json:"NHS_TRUD_KEY" desc:"See: https://isd.digital.nhs.uk/trud/users/authenticated/filters/0/api"`
	PostgresHost     string `envconfig:"POSTGRES_HOSTNAME" json:"POSTGRES_HOSTNAME" desc:"Postgres host name"`
	PostgresPort     uint   `envconfig:"POSTGRES_PORT" json:"POSTGRES_PORT" desc:"Postgres port"`
	PostgresUsername string `envconfig:"POSTGRES_USER" json:"POSTGRES_USER" desc:"Postgres username"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" json:"POSTGRES_PASSWORD" desc:"Postgres password"`
	PostgresDatabase string `envconfig:"POSTGRES_DB" json:"POSTGRES_DB" desc:"Postgres database name"`
}

var (
	Config  *ConfigSpec
	cfgOnce sync.Once
)

func RegisterEnvironment() (*ConfigSpec, error) {
	var err error
	cfgOnce.Do(func() {
		cfg := &ConfigSpec{}

		err = envconfig.Process("", cfg)
		if err != nil {
			return
		}

		Config = cfg
	})

	return Config, err
}

func EnvironmentUsage() string {
	builder := new(strings.Builder)
	writer := tabwriter.NewWriter(builder, 1, 0, 4, ' ', 0)

	err := envconfig.Usagef("", &Config, writer, envconfig.DefaultTableFormat)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	return builder.String()
}
