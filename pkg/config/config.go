package config

func (e RegistryType) IsUnknown() bool {
	return e == RegistryType_UNKNOWN
}

func (e RegistryType) IsEtcd() bool {
	return e == RegistryType_ETCD
}

func (e RegistryType) IsKubernetes() bool {
	return e == RegistryType_KUBERNETES
}
