package models

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Application struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	PublisherID uint
	Deployment  DeploymentSpec `gorm:"embedded"` // Embedded struct for deployment details
	Publisher   User           `gorm:"foreignKey:PublisherID"`
}

// DeploymentSpec stores deployment-related data
type DeploymentSpec struct {
	Type      string `gorm:"type:varchar(10)"` // "k8s" or "vm"
	RepoURL   string // Only for Kubernetes-based apps
	ChartName string // Only for Kubernetes-based apps
	Image     string // VM image for VM-based apps
	CPU       string // VM CPU configuration (e.g., "2 vCPUs")
	Memory    string // VM Memory configuration (e.g., "4GB RAM")
}

type Deployment struct {
	ID             uint `gorm:"primaryKey"`
	ConsumerID     uint
	ApplicationID  uint
	DeploymentType string `gorm:"type:varchar(10)"` // "k8s" or "vm"

	// Kubernetes-specific
	ClusterName string `gorm:"default:null"` // KIND cluster name (if K8s-based)

	// VM-specific
	VMName string `gorm:"default:null"` // VM instance name (if VM-based)
	VMIP   string `gorm:"default:null"` // IP of the created VM

	Consumer    User        `gorm:"foreignKey:ConsumerID"`
	Application Application `gorm:"foreignKey:ApplicationID"`
}
