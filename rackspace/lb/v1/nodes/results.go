package nodes

type Node struct {
	Address   string
	ID        int
	Port      int
	Status    Status
	Condition string
	Weight    int
}
