package primitive

type Directory struct {
	Root string
}

func (dir *Directory) Path() string {
	return dir.Root
}
