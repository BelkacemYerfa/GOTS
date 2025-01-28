export interface gormModel {
  UpdatedAt: Date;
  DeletedAt: Date;
  ID: number;
  CreatedAt: Date;
}

export interface Person {
  Name: string;
  Age: number;
  Hobbies: string[];
  Relations?: Relations;
  Happy: boolean;
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

export interface Collection extends gormModel {
  Components: Component[];
  Website: string;
  Version: string;
  UpdateAt: string;
  ReadmeUrl: string;
  Name: string;
  Users: User[];
  Frameworks: string[];
  Repository: string;
  Tags: string[];
}

export interface User extends gormModel {
  RefreshToken: string;
  Avatar: string;
  Collections: Collection[];
  Username: string;
  Email: string;
  AccessToken: string;
}

export interface File extends gormModel {
  Path: string;
  Type: string;
  Content: string;
  Target: string;
  ComponentID: number;
  Component: Component;
}

export interface CollectionDependency extends gormModel {
  Component: Component;
  Reference: string;
  ComponentID: number;
}

export interface Dependency extends gormModel {
  Component: Component;
  Name: string;
  ComponentID: number;
}

export interface UiDependency extends gormModel {
  Component: Component;
  Name: string;
  ComponentID: number;
}

export interface Version extends gormModel {
  VersionNumber: number;
  Version: string;
  Component: Component;
  ComponentID: number;
}

export interface Component extends gormModel {
  Frameworks: string[];
  Version: Version[];
  Doc: string;
  Type: string;
  Dependencies: Dependency[];
  UiDependencies: UiDependency[];
  CollectionDependencies: CollectionDependency[];
  Files: File[];
  CollectionID: number;
  Collection: Collection;
  Name: string;
  Tags: string[];
}

