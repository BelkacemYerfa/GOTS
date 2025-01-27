package main

type Person struct {
	Name      string     `json:"name"`
	Age       int        `json:"age"`
	Hobbies   []string   `json:"hobbies"`
	Relations *Relations `json:"relations"`
	Happy     bool       `json:"happy"`
}

type Relations struct {
	Parents  []string `json:"parents"`
	Siblings []string `json:"siblings"`
	Children []string `json:"children"`
}

type Address struct {
	Street string `json:"street"`
}

type (
	HAHA struct {
		Name string `json:"name"`
		hehe HEHE   `json:"hehe"`
	}

	HEHE struct {
		Age int `json:"age"`
	}
)

type (
	Collection struct {
		Name       string      `json:"name"`
		Components []Component `gorm:"foreignKey:CollectionID;" json:"components"`
		Users      []User      `gorm:"many2many:user_collection;" json:"users"`
		Frameworks string      `json:"frameworks"`
		Repository string      `json:"repository"`
		Website    string      `json:"website"`
		Tags       string      `json:"tags"`
		Version    string      `json:"version"`
		UpdateAt   string      `json:"update_at"`
		// ! file will be stored on aws s3
		ReadmeUrl string `json:"readme_url"`
	}
)
type User struct {
	Username     string       `gorm:"unique;not null" json:"username"`
	Email        string       `gorm:"unique;not null" json:"email"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	Avatar       string       `json:"avatar"`
	Collections  []Collection `gorm:"many2many:user_collection;" json:"collections"`
}

type File struct {
	Path        string    `gorm:"not null" json:"path"`
	Type        string    `gorm:"not null" json:"type"` // Type of file (e.g., component, utility, hook, etc.)
	Content     string    `json:"content"`
	Target      string    `json:"target"`                       // Target location inside the project
	ComponentID int       `gorm:"not null" json:"component_id"` // Foreign key to Component
	Component   Component `gorm:"constraint:OnDelete:CASCADE;"` // Relationship with cascading delete
}

type (
	CollectionDependency struct {
		ComponentID int       `gorm:"not null"`                     // Foreign key to Component
		Component   Component `gorm:"constraint:OnDelete:CASCADE;"` // Enable cascading delete
		Reference   string    `json:"reference"`
	}

	// Dependency struct represents external dependencies (e.g., npm libs).
	Dependency struct {
		Name        string    `json:"name"`
		ComponentID int       `gorm:"not null"`                     // Foreign key to Component
		Component   Component `gorm:"constraint:OnDelete:CASCADE;"` // Enable cascading delete
	}

	// UiDependency struct represents UI-specific dependencies (e.g., shadcn libs).
	UiDependency struct {
		Name        string    `json:"name"`
		ComponentID int       `gorm:"not null"`                     // Foreign key to Component
		Component   Component `gorm:"constraint:OnDelete:CASCADE;"` // Enable cascading delete
	}

	Version struct {
		VersionNumber uint64    `json:"version_number"`
		Version       string    `json:"version"`
		Component     Component `gorm:"constraint:OnDelete:CASCADE;"` // Enable cascading delete
		ComponentID   int       `gorm:"not null"`
	}

	// Component struct represents a reusable unit within a collection.
	Component struct {
		Name string `json:"name"`
		// Type of the component (e.g., component, utility, hook, etc.)
		Type string `json:"type"`
		// Relationships
		Dependencies           []Dependency           `gorm:"foreignKey:ComponentID;" json:"dependencies"`
		UiDependencies         []UiDependency         `gorm:"foreignKey:ComponentID;" json:"ui_dependencies"`
		CollectionDependencies []CollectionDependency `gorm:"foreignKey:ComponentID;" json:"collection_dependencies"`
		Files                  []File                 `gorm:"foreignKey:ComponentID;" json:"files"`
		// Documentation for the component
		Doc string `json:"doc"`
		// Collection Relationship
		CollectionID int        `gorm:"not null"`                     // Foreign key to Collection
		Collection   Collection `gorm:"constraint:OnDelete:CASCADE;"` // Enable cascading delete
		Tags         string     `json:"tags"`
		Frameworks   string     `json:"frameworks"`
		Version      []Version  `gorm:"foreignKey:ComponentID;" json:"version"`
	}
)
