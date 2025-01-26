export interface Person {
  Happy: boolean;
  Name: string;
  Age: number;
  Hobbies: string[];
  Relations?: Relations;
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
