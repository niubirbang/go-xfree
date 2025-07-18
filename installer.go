package goxfree

type Installer struct {
	client *client
}

func NewInstaller(option Option) *Installer {
	client := newClientInstaller(option)
	return &Installer{
		client: client,
	}
}

func (i *Installer) Run() error {
	if err := i.client.Run(); err != nil {
		return err
	}
	return nil
}

func (i *Installer) Quit() error {
	if err := i.client.Quit(); err != nil {
		return err
	}
	return nil
}
