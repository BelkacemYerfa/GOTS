export interface gormModel {
  DeletedAt: Date;
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
}

export interface Collection extends gormModel {
  Users: User[];
  Frameworks: string[];
  Tags: string[];
  ReadmeUrl: string;
  Name: string;
  Components: Component[];
  Repository: string;
  Website: string;
  Version: string;
  UpdateAt: string;
}

