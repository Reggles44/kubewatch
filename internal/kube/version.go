package kube

func Version() string {
	ver, err := client.Discovery().ServerVersion()
	if err != nil {
		panic(err)
	}

	return ver.String()
}
