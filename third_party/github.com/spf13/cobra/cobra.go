package cobra

type Command struct {
	Use     string
	Short   string
	Version string
	RunE    func(cmd *Command, args []string) error
	flags   FlagSet
}

func (c *Command) Execute() error {
	if c.RunE == nil {
		return nil
	}
	return c.RunE(c, nil)
}

func (c *Command) Flags() *FlagSet { return &c.flags }

type FlagSet struct{}

func (f *FlagSet) StringVarP(p *string, name, shorthand, value, usage string) { *p = value }
