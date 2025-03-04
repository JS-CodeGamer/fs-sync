export type Asset = {
  id: string;
  owner: string;
  name: string;
  parent_id: string;
  is_public: boolean;
  type: "folder" | "file";
  children?: Array<Asset>;
  files?: Array<File>;
  created_at: Date;
  updated_at: Date;
};

export type File = {
  id: string;
  size: string;
  version: string;
  mimeType: string;
};
