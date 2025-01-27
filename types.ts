export interface Person {
  Relations?: Relations;
  Happy: boolean;
  Name: string;
  Age: number;
  Hobbies: string[];
}

export interface Relations {
  Parents: string[];
  Siblings: string[];
  Children: string[];
}

export interface Address {
  Street: string;
}

export interface HAHA {
  Name: string;
  hehe: HEHE;
}

export interface HEHE {
  Age: number;
}

export interface Collection {
  UpdateAt: string;
  ReadmeUrl: string;
  Components: Component[];
  Repository: string;
  Website: string;
  Version: string;
  Name: string;
  Users: User[];
  Frameworks: string;
  Tags: string;
}

export interface User {
  RefreshToken: string;
  Avatar: string;
  Collections: Collection[];
  Username: string;
  Email: string;
  AccessToken: string;
}

export interface File {
  Path: string;
  Type: string;
  Content: string;
  Target: string;
  ComponentID: number;
  Component: Component;
}

export interface CollectionDependency {
  Component: Component;
  Reference: string;
  ComponentID: number;
}

export interface Dependency {
  Component: Component;
  Name: string;
  ComponentID: number;
}

export interface UiDependency {
  Name: string;
  ComponentID: number;
  Component: Component;
}

export interface Version {
  Component: Component;
  ComponentID: number;
  VersionNumber: number;
  Version: string;
}

export interface Component {
  CollectionID: number;
  Collection: Collection;
  Frameworks: string;
  Version: Version[];
  Name: string;
  UiDependencies: UiDependency[];
  CollectionDependencies: CollectionDependency[];
  Doc: string;
  Tags: string;
  Type: string;
  Dependencies: Dependency[];
  Files: File[];
}

