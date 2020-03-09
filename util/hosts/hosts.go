package hosts

type host struct {
	IsControl bool
	IsWorker  bool
	IsEtcd    bool
}
